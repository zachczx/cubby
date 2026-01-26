package migration

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/jmoiron/sqlx"
)

func WipeData(db *sqlx.DB) {
	log.Println("🔥 Truncating all tables...")

	query := `TRUNCATE TABLE users, families, trackers, invites, logs RESTART IDENTITY CASCADE;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌ Truncate failed: %v", err)
	}
	log.Println("✅ All data wiped successfully.")
}
func Create(db *sqlx.DB) {
	log.Println("🚀 Starting schema creation...")

	schema := []string{

		// Users Table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			email TEXT UNIQUE NOT NULL,
			name TEXT,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Families Table
		`CREATE TABLE IF NOT EXISTS families (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			name TEXT NOT NULL,
			owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Trackers Table (Depends on users and families)
		`CREATE TABLE IF NOT EXISTS trackers (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			name TEXT NOT NULL,
			description TEXT,
			interval_value INTEGER NOT NULL,
			interval_unit TEXT NOT NULL,
			category TEXT,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			family_id UUID REFERENCES families(id) ON DELETE CASCADE,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Invites Table (Depends on families and users)
		`CREATE TABLE IF NOT EXISTS invites (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID REFERENCES families(id) ON DELETE CASCADE,
			invited_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
			email TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			family_name_snapshot TEXT,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Logs Table (Depends on trackers and users)
		`CREATE TABLE IF NOT EXISTS logs (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			tracker_id UUID REFERENCES trackers(id) ON DELETE CASCADE,
			performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
			performed_at TIMESTAMPTZ DEFAULT NOW(),
			value TEXT,
			created TIMESTAMPTZ DEFAULT NOW(),
			updated TIMESTAMPTZ DEFAULT NOW()
		);`,
	}

	for i, query := range schema {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("❌ Error executing statement #%d:\nQuery: %s\nError: %v", i+1, query, err)
		}
	}

	log.Println("✅ Schema initialized successfully!")
}
