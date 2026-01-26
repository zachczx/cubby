package migration

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlitePath = "./data.db"
)

func Migrate() {
	// 1. Connect to SQLite (Source)
	sqliteDB, err := sqlx.Connect("sqlite3", sqlitePath)
	if err != nil {
		log.Fatalf("❌ Failed to open SQLite: %v", err)
	}
	defer sqliteDB.Close()

	pg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	pgDB, err := sqlx.Connect("pgx", pg)
	if err != nil {
		log.Fatal(err)
	}
	defer pgDB.Close()

	log.Println("⚠️  WARNING: This will DROP ALL TABLES in the public schema of the Postgres DB.")

	// 3. Reset Postgres Schema
	log.Println("🔥 Resetting database...")
	pgDB.MustExec("DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;")

	// 4. Get Tables from SQLite
	var tables []string
	err = sqliteDB.Select(&tables, `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != '_migrations'
	`)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	// 5. Create Tables in Postgres
	for _, tableName := range tables {
		if err := createTable(sqliteDB, pgDB, tableName); err != nil {
			log.Printf("❌ Error creating table %s: %v", tableName, err)
		} else {
			log.Printf("✅ Created table: %s", tableName)
		}
	}

	log.Println("🎉 Database reset and schema creation complete.")
}

func createTable(sqliteDB *sqlx.DB, pgDB *sqlx.DB, tableName string) error {
	// struct to hold PRAGMA info
	type colInfo struct {
		Name      string      `db:"name"`
		Type      string      `db:"type"`
		NotNull   int         `db:"notnull"`
		DfltValue interface{} `db:"dflt_value"` // interface{} because it can be null
		Pk        int         `db:"pk"`
	}

	var cols []colInfo
	// sqlx makes querying structs easy, but PRAGMA is weird, so we use Queryx loop or manual mapping
	// PRAGMA doesn't work well with Select() because it's not a standard SELECT
	rows, err := sqliteDB.Queryx(fmt.Sprintf("PRAGMA table_info(\"%s\")", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var c colInfo
		if err := rows.StructScan(&c); err != nil {
			return err
		}
		cols = append(cols, c)
	}

	var colDefs []string
	var pks []string

	for _, c := range cols {
		pgType := mapTypeToPostgres(c.Name, c.Type)
		def := fmt.Sprintf("\"%s\" %s", c.Name, pgType)

		if c.NotNull == 1 {
			def += " NOT NULL"
		}

		if c.DfltValue != nil {
			val := fmt.Sprintf("%v", c.DfltValue)
			// Handle Defaults logic (same as before)
			if strings.EqualFold(val, "CURRENT_TIMESTAMP") {
				def += " DEFAULT NOW()"
			} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
				def += fmt.Sprintf(" DEFAULT %s", val)
			} else {
				if pgType == "BOOLEAN" {
					if val == "1" || val == "'1'" {
						def += " DEFAULT TRUE"
					} else if val == "0" || val == "'0'" {
						def += " DEFAULT FALSE"
					}
				} else {
					def += fmt.Sprintf(" DEFAULT %s", val)
				}
			}
		}
		colDefs = append(colDefs, def)

		if c.Pk > 0 {
			pks = append(pks, fmt.Sprintf("\"%s\"", c.Name))
		}
	}

	query := fmt.Sprintf("CREATE TABLE \"%s\" (\n\t%s", tableName, strings.Join(colDefs, ",\n\t"))
	if len(pks) > 0 {
		query += fmt.Sprintf(",\n\tPRIMARY KEY (%s)", strings.Join(pks, ", "))
	}
	query += "\n);"

	_, err = pgDB.Exec(query)
	return err
}

func mapTypeToPostgres(colName, sqliteType string) string {
	upperType := strings.ToUpper(sqliteType)
	if colName == "created" || colName == "updated" {
		return "TIMESTAMPTZ"
	}
	if strings.Contains(strings.ToLower(colName), "json") || colName == "metadata" || colName == "settings" {
		return "JSONB"
	}
	switch {
	case strings.Contains(upperType, "INT"):
		return "BIGINT"
	case strings.Contains(upperType, "BOOL"):
		return "BOOLEAN"
	case strings.Contains(upperType, "REAL"), strings.Contains(upperType, "FLOA"), strings.Contains(upperType, "DOUB"):
		return "DOUBLE PRECISION"
	case strings.Contains(upperType, "JSON"):
		return "JSONB"
	case strings.Contains(upperType, "TEXT"), strings.Contains(upperType, "CHAR"):
		return "TEXT"
	case strings.Contains(upperType, "BLOB"):
		return "BYTEA"
	default:
		return "TEXT"
	}
}
