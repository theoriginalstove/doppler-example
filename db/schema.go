package db

import (
	"context"
	"fmt"
)

const (
	accountsTableName         = "accounts"
	accountsEmailsTableName   = "accounts_emails"
	accountPassResetTableName = "accounts_password_reset"
	plantsTableName           = "plants"
)

var allTableNames = []string{
	accountsTableName,
	accountPassResetTableName,
	plantsTableName,
}

func (i *Instance) ensureSchemas(ctx context.Context, suffix string) error {
	// account table name
	at := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id UUID PRIMARY KEY NOT NULL,
    name STRING,
    username VARCHAR(64) UNIQUE,
    pw BLOB,
    handle STRING UNIQUE,
    phone STRING UNIQUE,
    is_verified BOOL DEFAULT false,
    emails BLOB
)
    `, accountsTableName)

	_, err := i.db.Exec(ctx, at)
	if err != nil {
		return err
	}

	// accounts_emails table
	aet := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    email STRING NOT NULL UNIQUE,
    is_primary BOOL DEFAULT FALSE,
    is_verified BOOL DEFAULT FALSE,
    verification_token STRING,
    verification_expires TIMESTAMP
)
    `, accountsEmailsTableName)
	_, err = i.db.Exec(ctx, aet)
	if err != nil {
		return err
	}

	// plants table
	pt := fmt.Sprintf(` CREATE TABLE IF NOT EXISTS %s (
	    id STRING NOT NULL PRIMARY KEY,
	    common_name STRING NOT NULL,
	    scientific_name STRING NOT NULL,
	    varients JSONB
	)
	    `, plantsTableName)

	_, err = i.db.Exec(ctx, pt)
	if err != nil {
		return err
	}

	return nil
}

func (i *Instance) cleanupSchemas() error {
	for _, table := range allTableNames {
		q := fmt.Sprintf("DROP TABLE %s;", table)
		fmt.Printf("Dropping table: %s\n", table)
		_, err := i.db.Exec(context.Background(), q)
		if err != nil {
			return err
		}

	}
	return nil
}
