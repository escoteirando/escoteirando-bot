package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

func GetCount(tableName string) int64 {
	var count int64
	response := GetDB().Table(tableName).Count(&count)
	if response.Error == nil {
		log.Printf("Count(%s) = %d", tableName, count)
		return count
	}
	log.Printf("Count(%s) -> %v", tableName, response.Error)
	return 0
}

func ParseResponse(response *gorm.DB, description string, id interface{}) error {
	var err error
	if response.Error != nil {
		err = response.Error
	} else if response.RowsAffected < 1 {
		err = fmt.Errorf("returns empty")
	}
	if err != nil {
		log.Printf("Query [%s(%v)] exception: %v", description, id, err)
	}
	return err
}