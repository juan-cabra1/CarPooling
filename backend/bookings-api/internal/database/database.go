package database

import (
	"fmt"
	"strings"
	"time"

	"bookings-api/internal/config"
	"bookings-api/internal/dao"

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
// These settings prevent:
//   - Connection exhaustion (too many open connections)
//   - Stale connections (connections held too long)
//   - Resource waste (too many idle connections)
//
// GORM Logger Levels:
//   - Development: logger.Info (shows all SQL queries for debugging)
//   - Production: logger.Error (only logs errors, not queries)
//
// Parameters:
//   - cfg: Application configuration containing DatabaseURL
//
// Returns:
//   - *gorm.DB: GORM database instance for queries and transactions
//   - error: Error if connection fails or invalid DSN
//
// Example DatabaseURL (DSN format):
//
//	root:password@tcp(localhost:3306)/bookings_db?charset=utf8mb4&parseTime=True&loc=Local
//
// DSN Components:
//   - username:password - MySQL credentials
//   - tcp(host:port) - MySQL server address
//   - /database - Database name to connect to
//   - charset=utf8mb4 - Full Unicode support (including emojis)
//   - parseTime=True - Parse DATE/DATETIME to Go time.Time automatically
//   - loc=Local - Use local timezone for time parsing
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Validate that DatabaseURL is not empty
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required but not provided in configuration")
	}

	// Log connection attempt (hide password for security)
	sanitizedDSN := sanitizeDSN(cfg.DatabaseURL)
	log.Info().
		Str("dsn", sanitizedDSN).
		Msg("🔌 Connecting to MySQL database...")

	// Configure GORM logger based on environment
	// Development: Show all SQL queries for debugging
	// Production: Only show errors to reduce log noise
	gormLogger := logger.Default
	if cfg.IsProduction() {
		gormLogger = logger.Default.LogMode(logger.Error)
		log.Info().Msg("📝 GORM logger set to ERROR mode (production)")
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
		log.Info().Msg("📝 GORM logger set to INFO mode (development)")
	}

	// Open MySQL connection using GORM
	// GORM Config:
	//   - Logger: Configured based on environment
	//   - NamingStrategy: snake_case for tables/columns (Go convention)
	//   - DisableForeignKeyConstraintWhenMigrating: false (we want FK constraints)
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix: "",                  // No prefix for table names
			SingularTable: false, // Use plural table names (booking -> bookings)
			// NameReplacer: strings.NewReplacer("CID", "Cid"), // Optional: customize naming
		},
		// Keep foreign key constraints enabled for referential integrity
		DisableForeignKeyConstraintWhenMigrating: false,

		// PrepareStmt caches prepared statements for better performance
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	// Get underlying *sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	// ========================================================================
	// CONNECTION POOL CONFIGURATION
	// ========================================================================
	// These settings optimize database performance and resource usage

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	// Idle connections are kept alive for reuse, avoiding connection overhead
	// Too low: Frequent reconnections (slow)
	// Too high: Wasted resources
	// Recommendation: 10-25 for most applications
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	// This includes both idle and in-use connections
	// Too low: Connection pool exhaustion under load
	// Too high: Database server overload
	// MySQL default: 151 connections, leave room for other services
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	// After this time, the connection is closed and a new one is created
	// This prevents issues with:
	//   - Stale connections
	//   - Server-side connection timeouts
	//   - Connection state accumulation
	// Recommendation: 1 hour (MySQL wait_timeout is usually 8 hours)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().
		Int("max_idle_conns", 10).
		Int("max_open_conns", 100).
		Dur("conn_max_lifetime", time.Hour).
		Msg("⚙️  Connection pool configured")

	// ========================================================================
	// CONNECTION TEST
	// ========================================================================
	// Ping the database to verify connection is working
	// This catches connection errors early before attempting migrations
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	log.Info().
		Str("database", extractDatabaseName(cfg.DatabaseURL)).
		Msg("✅ MySQL connection established successfully")

	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models
//
// GORM AutoMigrate:
//   - Creates tables if they don't exist
//   - Adds new columns if struct has new fields
//   - Creates indexes defined in struct tags
//   - Creates UNIQUE constraints
//   - DOES NOT delete columns (safe for production)
//   - DOES NOT modify existing column types (requires manual migration)
//
// Tables created:
//  1. bookings - Main booking records
//     - Indexes: booking_uuid (unique), trip_id, passenger_id, status, cancelled_at
//  2. processed_events - Event idempotency tracking
//     - Indexes: event_id (unique), event_type, processed_at
//
// Migration Safety:
//   - AutoMigrate is safe for existing databases
//   - It only adds new tables/columns, never removes
//   - For complex schema changes, use manual migrations instead
//
// Parameters:
//   - db: GORM database instance
//
// Returns:
//   - error: Error if migration fails (e.g., permission denied, syntax error)
//
// Example:
//
//	db, _ := InitDB(cfg)
//	if err := AutoMigrate(db); err != nil {
//	    log.Fatal().Err(err).Msg("Migration failed")
//	}
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("🔄 Running database auto-migration...")

	// Migrate all models
	// Order doesn't matter since we have no foreign key constraints
	// between these tables (they're referenced by external IDs)
	err := db.AutoMigrate(
		&dao.Booking{},        // bookings table
		&dao.ProcessedEvent{}, // processed_events table
	)

	if err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Info().
		Strs("tables", []string{"bookings", "processed_events"}).
		Msg("✅ Database tables migrated successfully")

	// Log created indexes for verification
	log.Info().
		Strs("booking_indexes", []string{"booking_uuid", "trip_id", "passenger_id", "status", "cancelled_at"}).
		Strs("event_indexes", []string{"event_id (UNIQUE)", "event_type", "processed_at"}).
		Msg("📊 Database indexes created")

	return nil
}

// CloseDB closes the database connection gracefully
//
// This function should be called during application shutdown to:
//   - Close all open connections in the pool
//   - Release database resources
//   - Prevent connection leaks
//
// This is important for:
//   - Graceful shutdown (no abrupt connection drops)
//   - Resource cleanup (release MySQL connections)
//   - Avoiding "too many connections" errors
//
// Usage:
//
//	defer CloseDB(db) // Called automatically on program exit
//
// Or in main():
//
//	// On shutdown signal
//	if err := CloseDB(db); err != nil {
//	    log.Error().Err(err).Msg("Error closing database")
//	}
//
// Parameters:
//   - db: GORM database instance to close
//
// Returns:
//   - error: Error if closure fails (usually nil)
func CloseDB(db *gorm.DB) error {
	log.Info().Msg("🔌 Closing database connection...")

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	// Close all connections in the pool
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Info().Msg("✅ Database connection closed successfully")
	return nil
}

// ============================================================================
// HELPER FUNCTIONS (Private)
// ============================================================================

// sanitizeDSN removes the password from DSN for safe logging
//
// Example:
// Input:  "user:secret@tcp(localhost:3306)/db"
// Output: "user:***@tcp(localhost:3306)/db"
//
// This prevents password leaks in logs while still showing connection details
func sanitizeDSN(dsn string) string {
	// Find password section (between : and @)
	// Example DSN: username:password@tcp(host:port)/database?params
	if idx := strings.Index(dsn, ":"); idx != -1 {
		if idx2 := strings.Index(dsn[idx:], "@"); idx2 != -1 {
			// Replace password with ***
			return dsn[:idx+1] + "***" + dsn[idx+idx2:]
		}
	}
	// If format doesn't match, return as-is (no password to hide)
	return dsn
}

// extractDatabaseName extracts the database name from DSN
//
// Example:
// Input:  "user:pass@tcp(localhost:3306)/my_database?charset=utf8"
// Output: "my_database"
//
// Used for logging to confirm which database we connected to
func extractDatabaseName(dsn string) string {
	// Find database name (between / and ?)
	// Example DSN: username:password@tcp(host:port)/database?params
	if idx := strings.Index(dsn, "/"); idx != -1 {
		dbPart := dsn[idx+1:]
		// Remove query parameters if present
		if idx2 := strings.Index(dbPart, "?"); idx2 != -1 {
			return dbPart[:idx2]
		}
		return dbPart
	}
	return "unknown"
}
