package database

import (
	"log"
	"time"

	"github.com/Guesstrain/ethglobal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseService defines the methods to interact with the database
type DatabaseService interface {
	Insert(data interface{}) error
	SelectByField(dest interface{}, fieldName string, value interface{}) error
	DeleteByField(model interface{}, fieldName string, value interface{}) error
	QueryWalletsByTimePeriod(startTime, endTime time.Time) ([]models.Wallet, error)
	UpdateByStruct(model interface{}, fieldName string, value interface{}, updatedData interface{}) error
}

// GormDatabaseService implements the DatabaseService interface using GORM
type GormDatabaseService struct {
	db *gorm.DB
}

// NewGormDatabaseService initializes a new GormDatabaseService
func NewGormDatabaseService(dsn string) DatabaseService {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return nil
	}

	return &GormDatabaseService{db: db}
}

// Insert inserts data into the database (works for any struct)
func (s *GormDatabaseService) Insert(data interface{}) error {
	result := s.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// SelectByField selects data from the database based on a specific field
func (s *GormDatabaseService) SelectByField(dest interface{}, fieldName string, value interface{}) error {
	result := s.db.Where(fieldName+" = ?", value).First(dest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateByField updates data in the database based on a specific field
func (s *GormDatabaseService) UpdateByStruct(model interface{}, fieldName string, value interface{}, updatedData interface{}) error {
	result := s.db.Model(model).Where(fieldName+" = ?", value).Updates(updatedData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteByField deletes data from the database based on a specific field
func (s *GormDatabaseService) DeleteByField(model interface{}, fieldName string, value interface{}) error {
	result := s.db.Where(fieldName+" = ?", value).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *GormDatabaseService) QueryWalletsByTimePeriod(startTime, endTime time.Time) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := s.db.Where("timestamp BETWEEN ? AND ?", startTime, endTime).Find(&wallets).Error
	return wallets, err
}
