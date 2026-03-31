package database

import (
	"fmt"
	"strings"
	"time"

	"payments-api/internal/config"
	"payments-api/internal/dao"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB initializes the MySQL database connection using GORM
//
// This function:
//  1. Parses the MySQL DSN (Data Source Name) from config
//  2. Opens a GORM connection with the MySQL driver
//  3. Configures connection pooling for optimal performance
//  4. Tests the connection with Ping()
//  5. Configures GORM logger based on environment
//
// Connection Pooling:
//   - MaxIdleConns: 10 - Maximum idle connections in the pool
//   - MaxOpenConns: 100 - Maximum total open connections
//   - ConnMaxLifetime: 1 hour - Maximum time a connection can be reused
//
// Parameters:
//   - cfg: Application configuration containing DatabaseURL
//
// Returns:
//   - *gorm.DB: GORM database instance for queries and transactions
//   - error: Error if connection fails or invalid DSN
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Validate that DatabaseURL is not empty
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required but not provided in configuration")
	}

	// Log connection attempt (hide password for security)
	sanitizedDSN := sanitizeDSN(cfg.DatabaseURL)
	log.Info().
		Str("dsn", sanitizedDSN).
		Msg("Connecting to MySQL database...")

	// Configure GORM logger based on environment
	gormLogger := logger.Default
	if cfg.IsProduction() {
		gormLogger = logger.Default.LogMode(logger.Error)
		log.Info().Msg("GORM logger set to ERROR mode (production)")
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
		log.Info().Msg("GORM logger set to INFO mode (development)")
	}

	// Open MySQL connection using GORM
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // Use plural table names
		},
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	// Get underlying *sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	// Connection pool configuration
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().
		Int("max_idle_conns", 10).
		Int("max_open_conns", 100).
		Dur("conn_max_lifetime", time.Hour).
		Msg("Connection pool configured")

	// Ping the database to verify connection is working
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	log.Info().
		Str("database", extractDatabaseName(cfg.DatabaseURL)).
		Msg("MySQL connection established successfully")

	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models
//
// Tables created:
//  1. payments - Payment records
//  2. wallets - Driver wallet balances
//  3. wallet_transactions - Wallet movement history
//  4. withdrawals - Withdrawal requests
//  5. seller_accounts - MP linked accounts
//  6. processed_events - Event idempotency tracking
//
// Parameters:
//   - db: GORM database instance
//
// Returns:
//   - error: Error if migration fails
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Running database auto-migration...")

	// Migrate all models
	err := db.AutoMigrate(
		&dao.SellerAccount{},
		&dao.Payment{},
		&dao.Wallet{},
		&dao.WalletTransaction{},
		&dao.Withdrawal{},
		&dao.ProcessedEvent{},
	)

	if err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Info().
		Strs("tables", []string{
			"seller_accounts",
			"payments",
			"wallets",
			"wallet_transactions",
			"withdrawals",
			"processed_events",
		}).
		Msg("Database tables migrated successfully")

	return nil
}

// CloseDB closes the database connection gracefully
//
// Parameters:
//   - db: GORM database instance to close
//
// Returns:
//   - error: Error if closure fails (usually nil)
func CloseDB(db *gorm.DB) error {
	log.Info().Msg("Closing database connection...")

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Info().Msg("Database connection closed successfully")
	return nil
}

// sanitizeDSN removes the password from DSN for safe logging
func sanitizeDSN(dsn string) string {
	if idx := strings.Index(dsn, ":"); idx != -1 {
		if idx2 := strings.Index(dsn[idx:], "@"); idx2 != -1 {
			return dsn[:idx+1] + "***" + dsn[idx+idx2:]
		}
	}
	return dsn
}

// extractDatabaseName extracts the database name from DSN
func extractDatabaseName(dsn string) string {
	if idx := strings.Index(dsn, "/"); idx != -1 {
		dbPart := dsn[idx+1:]
		if idx2 := strings.Index(dbPart, "?"); idx2 != -1 {
			return dbPart[:idx2]
		}
		return dbPart
	}
	return "unknown"
}
