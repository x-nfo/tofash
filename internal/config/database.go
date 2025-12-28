package config

import (
	"fmt"

	orderModel "tofash/internal/modules/order/model"
	productSeeds "tofash/internal/modules/product/database/seeds"
	productModel "tofash/internal/modules/product/model"
	"tofash/internal/modules/user/database/seeds"
	userModel "tofash/internal/modules/user/model"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Legacy Wrapper for backward compatibility
type Postgres struct {
	DB *gorm.DB
}

// ConnectionPostgres is used by legacy consumers to get a fresh connection or singleton
func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	// Reusing the logic from InitDatabase but returning Postgres struct
	// Note: AutoMigrate is done in InitDatabase. Here we just connect.
	// Optimization: Ideally reuse the db instance from InitDatabase if it was global.
	// But legacy code calls this. We will create a new connection or use `InitDatabase` logic.

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres] Failed to connect to database " + cfg.Psql.Host)
		return nil, err
	}

	// Ensure connection pool settings are applied
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
		sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)
	}

	return &Postgres{DB: db}, nil
}

func InitDatabase(cfg *Config) *gorm.DB {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[InitDatabase] Failed to connect to database " + cfg.Psql.Host)
		return nil
	}

	// Auto Migrate ALL modules
	err = db.AutoMigrate(
		// User Module
		&userModel.User{},
		&userModel.Role{},
		&userModel.UserRole{},
		&userModel.VerificationToken{},

		// Product Module
		&productModel.Category{},
		&productModel.Product{},

		// Order Module
		&orderModel.Order{},
		&orderModel.OrderItem{},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("[InitDatabase] Failed to migrate database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[InitDatabase] Failed to get database connection")
		return nil
	}

	// Seeds (User module has seeds)
	seeds.SeedRole(db)
	seeds.SeedAdmin(db)
	productSeeds.SeedProduct(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return db
}
