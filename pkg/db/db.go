package database

import (
	"discordbot/pkg/model"

	"github.com/erhanakp/sugaredgorm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseManager struct {
	DB *gorm.DB
}

func InitDatabaseManager() (*DatabaseManager, error) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	gormlogger := sugaredgorm.New(logger, sugaredgorm.Config{})

	var (
		db  *gorm.DB
		err error
	)

	if viper.GetString("db.type") == "sqlite" {
		dbFile := viper.GetString("db.file")
		db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
			Logger: gormlogger,
		})
		if err != nil {
			logger.Errorln("failed to connect database")
			return nil, err
		}
	} else {
		logger.Errorln("unknown db type")
		return nil, err
	}
	return &DatabaseManager{
		DB: db,
	}, nil
}

func (d *DatabaseManager) Migrate() (err error) {
	logger := zap.L().Sugar()
	defer logger.Sync()

	logger.Debug("Migrating database")
	err = d.DB.AutoMigrate(&model.User{})
	if err != nil {
		logger.Panicf("failed to migrate database: %v", err)
	}
	return
}
