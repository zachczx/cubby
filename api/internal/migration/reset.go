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
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
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

		// entries
		`CREATE TABLE IF NOT EXISTS entries (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			tracker_id UUID NOT NULL REFERENCES trackers(id) ON DELETE CASCADE,
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

		// Invites
		`CREATE TABLE IF NOT EXISTS invites (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
			invitee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(family_id, invitee_id)
		);`,

		`CREATE TABLE IF NOT EXISTS push_tokens (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token TEXT NOT NULL,
			platform TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(user_id, token)
		);`,

		`CREATE TABLE IF NOT EXISTS notification_logs (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			tracker_id UUID NOT NULL REFERENCES trackers(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tracker_id, user_id)
		);`,

		// FK, lookup indexes
		`CREATE INDEX IF NOT EXISTS idx_families_owner_id ON families(owner_id);`,
		`CREATE INDEX IF NOT EXISTS idx_families_users_user_id ON families_users(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_trackers_owner_id ON trackers(owner_id);`,
		`CREATE INDEX IF NOT EXISTS idx_trackers_family_id ON trackers(family_id);`,
		`CREATE INDEX IF NOT EXISTS idx_entries_tracker_id ON entries(tracker_id);`,
		`CREATE INDEX IF NOT EXISTS idx_entries_performed_by ON entries(performed_by);`,
		`CREATE INDEX IF NOT EXISTS idx_vacations_family_id ON vacations(family_id);`,
		`CREATE INDEX IF NOT EXISTS idx_invites_invitee_id ON invites(invitee_id);`,
		`CREATE INDEX IF NOT EXISTS idx_push_tokens_user_id ON push_tokens(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notification_logs_lookup ON notification_logs(tracker_id, user_id);`,

		// Date time filter indexes
		`CREATE INDEX IF NOT EXISTS idx_entries_tracker_time ON entries(tracker_id, performed_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_vacations_dates ON vacations(start_date_time, end_date_time);`,

		// Partial Indexes for highly filtered data
		`CREATE INDEX IF NOT EXISTS idx_trackers_active_owner ON trackers(owner_id) WHERE show = true;`,
		`CREATE INDEX IF NOT EXISTS idx_trackers_active_family ON trackers(family_id) WHERE show = true;`,
		`CREATE INDEX IF NOT EXISTS idx_invites_pending ON invites(invitee_id) WHERE status = 'pending';`,

		// Global sort indexes
		`CREATE INDEX IF NOT EXISTS idx_entries_performed_at ON entries(performed_at DESC);`,

		// Case insensitive lookups
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_lower ON users(LOWER(email));`,
	}

	for i, query := range schema {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("❌ Error executing statement #%d:\nQuery: %s\nError: %v", i+1, query, err)
		}
	}

	log.Println("✅ Schema initialized successfully!")
}
