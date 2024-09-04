package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	// Create users_table
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS users_table (
            id BIGINT PRIMARY KEY,    
            privacy TEXT NOT NULL,
            h3_index TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create users_table: %v", err)
	}

	// Create location_history table
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS location_history (
            id SERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL,            
            h3_index TEXT NOT NULL,
            timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users_table(id)
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create location_history table: %v", err)
	}

	// Create location_sharing table
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS location_sharing (
            id SERIAL PRIMARY KEY,
            sharer_id BIGINT NOT NULL,
            receiver_id BIGINT NOT NULL,
            start_time TIMESTAMP WITH TIME ZONE NOT NULL,
            end_time TIMESTAMP WITH TIME ZONE,
            FOREIGN KEY (sharer_id) REFERENCES users_table(id),
            FOREIGN KEY (receiver_id) REFERENCES users_table(id)
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create location_sharing table: %v", err)
	}

	// Create customers table
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS customers (
            id SERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL,
            customer_id TEXT NOT NULL,
            name TEXT NOT NULL,
            email TEXT NOT NULL,
            status TEXT NOT NULL DEFAULT 'incomplete' ,
            subscription_end_date TIMESTAMP,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users_table(id)
        );
    `)
	if err != nil {
		return fmt.Errorf("could not create customers table: %w", err)
	}

	// Create indexes
	_, err = tx.Exec(`
        CREATE INDEX IF NOT EXISTS idx_users_table_h3_index ON users_table(h3_index);
        CREATE INDEX IF NOT EXISTS idx_location_history_user_id ON location_history(user_id);
        CREATE INDEX IF NOT EXISTS idx_location_sharing_sharer_id ON location_sharing(sharer_id);
        CREATE INDEX IF NOT EXISTS idx_location_sharing_receiver_id ON location_sharing(receiver_id);
        CREATE INDEX IF NOT EXISTS idx_customers_user_id ON customers(user_id);
    `)
	if err != nil {
		return fmt.Errorf("could not create indexes: %v", err)
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	// Drop tables and indexes in reverse order of creation
	_, err := tx.Exec(`
    DROP TABLE IF EXISTS location_sharing;
    DROP TABLE IF EXISTS location_history;
    DROP TABLE IF EXISTS users_table;
    DROP TABLE IF EXISTS customers;
    `)
	if err != nil {
		return fmt.Errorf("could not drop tables: %v", err)
	}
	return nil
}
