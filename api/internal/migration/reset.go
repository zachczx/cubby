package migration

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/jmoiron/sqlx"
)

func WipeData(db *sqlx.DB) {
	log.Println("🔥 Truncating all tables...")

	query := `DROP TABLE IF EXISTS entries, invites, trackers, families, users CASCADE;`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌ Drop tables failed: %v", err)
	}
	log.Println("✅ All tables dropped successfully.")
}

func Create(db *sqlx.DB) {
	log.Println("🚀 Starting schema creation...")

	schema := []string{
		// Users
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			email TEXT UNIQUE NOT NULL,
			name TEXT,
			task_lookahead_days INTEGER DEFAULT 14,
			sound_on BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Families
		`CREATE TABLE IF NOT EXISTS families (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			name TEXT NOT NULL,
			owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		`CREATE TABLE IF NOT EXISTS families_users (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(family_id, user_id)
		);`,

		// Trackers
		`CREATE TABLE IF NOT EXISTS trackers (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
			family_id UUID REFERENCES families(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			display TEXT,           
			interval INTEGER NOT NULL,
			interval_unit TEXT NOT NULL,    
			category TEXT,            
			kind TEXT,              
			action_label TEXT,      
			icon TEXT,               
			pinned BOOLEAN DEFAULT FALSE,
			show BOOLEAN DEFAULT TRUE,  
			start_date TIMESTAMPTZ,  
			cost DOUBLE PRECISION,    
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// Invites
		`CREATE TABLE IF NOT EXISTS invites (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID REFERENCES families(id) ON DELETE CASCADE,
			invited_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
			email TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			family_name_snapshot TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// entries
		`CREATE TABLE IF NOT EXISTS entries (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			tracker_id UUID REFERENCES trackers(id) ON DELETE CASCADE,
			interval INTEGER NOT NULL,
			interval_unit VARCHAR(10) NOT NULL,
			performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
			performed_at TIMESTAMPTZ DEFAULT NOW(),
			remark TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// vacation
		`CREATE TABLE IF NOT EXISTS vacations (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
			created_by UUID REFERENCES users(id) ON DELETE SET NULL,
			start_date_time TIMESTAMPTZ NOT NULL,
			end_date_time TIMESTAMPTZ NOT NULL,
			label TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT valid_vacation_period CHECK (end_date_time > start_date_time)
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
