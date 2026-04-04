package migration

import (
	"log"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/jmoiron/sqlx"
)

func WipeData(db *sqlx.DB) {
	query := `DROP TABLE IF EXISTS timer_profiles, gym_routine_exercises, gym_routines, gym_sets, gym_workouts, tracker_user_settings, notification_logs, push_tokens, invites, vacations, entries, trackers, families_users, families, users, market_prices CASCADE;`
	_, err := db.Exec(query)
	if err != nil {
		slog.Error("failed to drop tables", "error", err)
		os.Exit(1)
	}

	slog.Info("drop tables", "status", "success")
}

func Create(db *sqlx.DB) {
	slog.Info("create schema", "status", "success")

	schema := []string{
		// Users
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			email TEXT UNIQUE NOT NULL,
			name TEXT,
			task_lookahead_days INTEGER DEFAULT 14,
			sound_on BOOLEAN DEFAULT TRUE,
			preferred_character VARCHAR(50) DEFAULT 'default',
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
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(name, family_id)
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

		`CREATE TABLE IF NOT EXISTS tracker_user_settings (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			tracker_id UUID NOT NULL REFERENCES trackers(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			is_muted BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(tracker_id, user_id)
		);`,

		// Gym
		`CREATE TABLE IF NOT EXISTS gym_workouts (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		`CREATE TABLE IF NOT EXISTS gym_sets (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			workout_id UUID NOT NULL REFERENCES gym_workouts(id) ON DELETE CASCADE,
			exercise_id TEXT NOT NULL,
			weight_kg NUMERIC(6,2),
			reps SMALLINT,
			set_type VARCHAR(50) DEFAULT 'working',
			is_completed BOOLEAN DEFAULT FALSE,
			position SMALLINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		`CREATE TABLE IF NOT EXISTS gym_favourite_exercises (
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			exercise_id TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (user_id, exercise_id)
		);`,

		// Gym Routines
		`CREATE TABLE IF NOT EXISTS gym_routines (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			position SMALLINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		`CREATE TABLE IF NOT EXISTS gym_routine_exercises (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			routine_id UUID NOT NULL REFERENCES gym_routines(id) ON DELETE CASCADE,
			exercise_id TEXT NOT NULL,
			sets SMALLINT NOT NULL DEFAULT 3,
			position SMALLINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(routine_id, exercise_id)
		);`,

		// Market Prices
		`CREATE TABLE IF NOT EXISTS market_prices (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
			logged_by UUID REFERENCES users(id) ON DELETE SET NULL,
			item_name TEXT NOT NULL,
			category TEXT,
			country TEXT,
			store TEXT,
			unit TEXT,
			quantity NUMERIC(8,2),
			price NUMERIC(8,2) NOT NULL,
			is_promo BOOLEAN DEFAULT FALSE,
			remarks TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
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
		`CREATE INDEX IF NOT EXISTS idx_tracker_user_settings_user_id ON tracker_user_settings(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tracker_user_settings_tracker_id ON tracker_user_settings(tracker_id);`,

		`CREATE INDEX IF NOT EXISTS idx_gym_workouts_user_id ON gym_workouts(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_gym_sets_workout_id ON gym_sets(workout_id);`,
		`CREATE INDEX IF NOT EXISTS idx_gym_sets_chronological ON gym_sets(workout_id, created_at ASC);`,

		`CREATE INDEX IF NOT EXISTS idx_gym_routines_user_id ON gym_routines(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_gym_routine_exercises_routine_id ON gym_routine_exercises(routine_id);`,

		`CREATE INDEX IF NOT EXISTS idx_market_prices_family_id ON market_prices(family_id);`,
		`CREATE INDEX IF NOT EXISTS idx_market_prices_item_name ON market_prices(item_name);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_market_prices_no_duplicate ON market_prices(family_id, LOWER(item_name), COALESCE(LOWER(store), ''), price, DATE(created_at AT TIME ZONE 'UTC'));`,

		// Date time filter indexes
		`CREATE INDEX IF NOT EXISTS idx_entries_tracker_time ON entries(tracker_id, performed_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_vacations_dates ON vacations(start_date_time, end_date_time);`,

		`CREATE INDEX IF NOT EXISTS idx_gym_workouts_start_time ON gym_workouts(start_time DESC);`,

		// Timer profiles
		`CREATE TABLE IF NOT EXISTS timer_profiles (
			id UUID PRIMARY KEY DEFAULT uuidv7(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			segments JSONB NOT NULL DEFAULT '[]'::jsonb,
			is_default BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		`CREATE INDEX IF NOT EXISTS idx_timer_profiles_user_id ON timer_profiles(user_id);`,
		// Partial unique index: PG doesn't support WHERE on constraints, but unique indexes enforce the same way.
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_timer_profiles_one_default ON timer_profiles(user_id) WHERE is_default = TRUE;`,

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
			log.Fatalf("error executing statement #%d:\nQuery: %s\nError: %v", i+1, query, err)
		}
	}

	slog.Info("schema init", "status", "success")
}
