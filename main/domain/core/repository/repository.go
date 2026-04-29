package repository

import (
	"moss/domain/core/entity"
	"moss/infrastructure/persistent/db"
	"moss/infrastructure/support/log"
)

func init() {
	MigrateTable()
}

func MigrateTable() {

	if err := Article.MigrateTable(); err != nil {
		log.Error("migrate article table error", log.Err(err))
	}

	if err := Category.MigrateTable(); err != nil {
		log.Error("migrate category table error", log.Err(err))
	}

	if err := Tag.MigrateTable(); err != nil {
		log.Error("migrate tag table error", log.Err(err))
	}

	if err := Mapping.MigrateTable(); err != nil {
		log.Error("migrate mapping table error", log.Err(err))
	}

	if err := Link.MigrateTable(); err != nil {
		log.Error("migrate link table error", log.Err(err))
	}

	if err := Store.MigrateTable(); err != nil {
		log.Error("migrate store table error", log.Err(err))
	}

	// New tables for frontend features
	if err := db.DB.AutoMigrate(&entity.Favorite{}); err != nil {
		log.Error("migrate favorite table error", log.Err(err))
	}

	if err := db.DB.AutoMigrate(&entity.Like{}); err != nil {
		log.Error("migrate like table error", log.Err(err))
	}

	if err := db.DB.AutoMigrate(&entity.ViewHistory{}); err != nil {
		log.Error("migrate view_history table error", log.Err(err))
	}

	if err := db.DB.AutoMigrate(&entity.User{}); err != nil {
		log.Error("migrate user table error", log.Err(err))
	}
}
