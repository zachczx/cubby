package migration

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zachczx/cubby/api/internal/entry"
)

const (
	sqlitePath = "./internal/migration/data.db"
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

	if err := MigrateEntries(sqliteDB, pgDB); err != nil {
		fmt.Println(err)
	}
}

// var trackerUUIDs = Trackers{
// Gummy:       uuid.MustParse("019c5b83-e669-707a-9e66-c3f8032d1168"),
// Spray:       uuid.MustParse("019c5b83-e66e-731b-b3f7-962df23bc89b"),
// Towel:       uuid.MustParse("019c5b83-e672-7f03-9122-bca3a4c1508e"),
// PetChewable: uuid.MustParse("019c5b83-e664-710e-a12e-89755f2e2b33"),
// PetBath:  uuid.MustParse("019c5b83-e65e-7f34-a2ff-baade24726c4"),
// BedSheet: uuid.MustParse("019c5b83-e659-705a-b03e-f45d7b229314"),
// }

var mapperIDs = map[string]IDPair{
	"8t9hsvqah63rs7h": {
		Name: "Gummy",
		UUID: uuid.MustParse("019c5b83-e669-707a-9e66-c3f8032d1168"),
	},
	"381t91o03thrvyd": {
		Name: "Spray",
		UUID: uuid.MustParse("019c5b83-e66e-731b-b3f7-962df23bc89b"),
	},
	"vvd9jnl0uw8qnie": {
		Name: "Towel",
		UUID: uuid.MustParse("019c5b83-e672-7f03-9122-bca3a4c1508e"),
	},
	"h3e3xkbmoxma6dv": {
		Name: "PetChewable",
		UUID: uuid.MustParse("019c5b83-e664-710e-a12e-89755f2e2b33"),
	},
	"hhz09lsfj1o5mbp": {
		Name: "PetBath",
		UUID: uuid.MustParse("019c5b83-e65e-7f34-a2ff-baade24726c4"),
	},
	"vk58159wczyxmus": {
		Name: "BedSheet",
		UUID: uuid.MustParse("019c5b83-e659-705a-b03e-f45d7b229314"),
	},
	"3u9v9goodptrvez": {
		Name: "M365Personal",
		UUID: uuid.MustParse("019c5c92-853b-71dc-9a82-18e2a5f6b4b1"),
	},
	"n63w3ft8oxh41y4": {
		Name: "AdobeCreativeCloud",
		UUID: uuid.MustParse("019c5c94-270e-76e7-83df-98ff472792ab"),
	},
	"wnfhp96bg54b41v": {
		Name: "PayOCBCCreditCards",
		UUID: uuid.MustParse("019c5c95-72c1-7c5b-a6c7-66ee2948f9b0"),
	},
	"q458g0bnwmygf7b": {
		Name: "KitVaccination",
		UUID: uuid.MustParse("019c5c96-92fb-7f0a-838e-e1c25b527bca"),
	},
	"rbq3jszzfw49u21": {
		Name: "Electricity",
		UUID: uuid.MustParse("019c5c97-cd2b-7ec8-8bb7-ed1ae99d8860"),
	},
	"fdc5aof8fdz3jnr": {
		Name: "PQSleep",
		UUID: uuid.MustParse("019c5cb5-5d41-7182-8981-5699d7873e57"),
	},
	"t72sc8e8652r9vv": {
		Name: "DentalCleaning",
		UUID: uuid.MustParse("019c5cb4-e5de-7a07-8a03-2ecd648c57d0"),
	},
	"932ed1fenn3uv6n": {
		Name: "WasherCleaning",
		UUID: uuid.MustParse("019c5cb4-1451-7a2f-83eb-985fe32c4d38"),
	},
}

var currentUserUUID = uuid.MustParse("019c5b83-e63e-75ba-b532-007c77f129a5")

// var sqliteIDs = SqliteTrackers{
// 	Gummy:              "8t9hsvqah63rs7h",
// 	Spray:              "381t91o03thrvyd",
// 	Towel:              "vvd9jnl0uw8qnie",
// 	PetChewable:        "h3e3xkbmoxma6dv",
// 	PetBath:            "hhz09lsfj1o5mbp",
// 	BedSheet:           "vk58159wczyxmus",
// 	Electricity:        "rbq3jszzfw49u21",
// 	KitVaccination:     "q458g0bnwmygf7b",
// 	PayOCBCCreditCards: "q458g0bnwmygf7b",
// 	AdobeCreativeCloud: "n63w3ft8oxh41y4",
// 	M365Personal:       "3u9v9goodptrvez",
// }

type IDPair struct {
	Name string
	UUID uuid.UUID
}

// type Trackers struct {
// 	Spray       uuid.UUID
// 	Gummy       uuid.UUID
// 	Towel       uuid.UUID
// 	BedSheet    uuid.UUID
// 	PetBath     uuid.UUID
// 	PetChewable uuid.UUID
// }

// type SqliteTrackers struct {
// 	Spray              string
// 	Gummy              string
// 	Towel              string
// 	BedSheet           string
// 	PetBath            string
// 	PetChewable        string
// 	Electricity        string
// 	KitVaccination     string
// 	PayOCBCCreditCards string
// 	AdobeCreativeCloud string
// 	M365Personal       string
// }

type Entry struct {
	ID           string  `db:"id"`
	Created      string  `db:"created"`
	Updated      string  `db:"updated"`
	Interval     int     `db:"interval"`
	IntervalUnit string  `db:"intervalUnit"`
	Time         string  `db:"time"`
	Tracker      string  `db:"tracker"`
	Remark       *string `db:"remark"`
}

const pbTimeFormat = "2006-01-02 15:04:05.999Z07:00"

func MigrateEntries(sqliteDB *sqlx.DB, pgDB *sqlx.DB) error {
	var sqliteEntries []Entry

	if err := sqliteDB.Select(&sqliteEntries, `SELECT * FROM logs`); err != nil {
		return fmt.Errorf("select entries err: %w", err)
	}

	var pgEntries []entry.Entry

	for _, s := range sqliteEntries {
		pair, ok := mapperIDs[s.Tracker]
		if !ok {
			return fmt.Errorf("no mapper found for tracker ID %s", s.Tracker)
		}

		t, err := time.Parse(pbTimeFormat, s.Time)
		if err != nil {
			return fmt.Errorf("format time err: %w", err)
		}

		createdAt, err := time.Parse(pbTimeFormat, s.Created)
		if err != nil {
			return fmt.Errorf("format time err: %w", err)
		}

		updatedAt, err := time.Parse(pbTimeFormat, s.Updated)
		if err != nil {
			return fmt.Errorf("format time err: %w", err)
		}

		p := entry.Entry{
			TrackerID:    pair.UUID,
			Interval:     s.Interval,
			IntervalUnit: s.IntervalUnit,
			PerformedBy:  currentUserUUID,
			PerformedAt:  t,
			Remark:       *s.Remark,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		}

		pgEntries = append(pgEntries, p)
	}

	q := `INSERT INTO entries (
			tracker_id, 
			interval, 
			interval_unit, 
			performed_by, 
			performed_at, 
			remark, 
			created_at, 
			updated_at)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8
			)`

	for _, p := range pgEntries {

		_, err := pgDB.Exec(q, p.TrackerID, p.Interval, p.IntervalUnit, p.PerformedBy, p.PerformedAt, p.Remark, p.CreatedAt, p.UpdatedAt)
		if err != nil {
			return fmt.Errorf("insert %s: %w", p.TrackerID, err)
		}
	}

	return nil
}
