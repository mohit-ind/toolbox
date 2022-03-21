package database

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	logger "github.com/toolboxlogger"
)

// NewMSSQLGormDB creates a GORM DB session towards a Microsoft SQL Server
// it uses the supplied Logger to create the GORM Logger.
func NewMSSQLGormDB(
	log *logger.Logger,
	conf *gorm.Config,
	connectionString string,
	maxOpenConn,
	maxIdleConn int,
	maxConnLifeTime time.Duration) (*gorm.DB, error) {

	if conf == nil {
		conf = &gorm.Config{}
	}
	conf.Logger = log.NewGormLogger("GORM")
	gormDB, err := gorm.Open(sqlserver.Open(connectionString), conf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create the GORM session towards the Microsoft SQL server")
	}
	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	if maxConnLifeTime > 0 {
		db.SetConnMaxLifetime(maxConnLifeTime)
	}
	return gormDB, nil
}
