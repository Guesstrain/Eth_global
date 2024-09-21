package services

import (
	"fmt"

	"github.com/Guesstrain/ethglobal/database"
	"github.com/Guesstrain/ethglobal/models"
)

type WalletService interface {
	InsertWallet(wallet models.Wallet) error
	SelectWalletByAddress(address string) (*models.Wallet, error)
	DeleteWalletByAddress(address string) error
}

type WalletServiceImpl struct {
	dbService database.DatabaseService
}

func NewWalletService(dbService database.DatabaseService) WalletService {
	return &WalletServiceImpl{dbService: dbService}
}

// var dsn = "user:password@tcp(127.0.0.1:3306)/dbname"
// var dbService, _ = database.NewGormDatabaseService(dsn)

func (s *WalletServiceImpl) InsertWallet(wallet models.Wallet) error {

	err := s.dbService.Insert(&wallet)
	if err != nil {
		fmt.Println("Failed to insert wallet:", err)
	}
	return err
}

func (s *WalletServiceImpl) SelectWalletByAddress(address string) (*models.Wallet, error) {
	var wallet models.Wallet

	err := s.dbService.SelectByField(&wallet, "Address", address)
	if err != nil {
		fmt.Println("Failed to find wallet with address:", address, err)
		return nil, err
	}
	return &wallet, nil
}

func (s *WalletServiceImpl) DeleteWalletByAddress(address string) error {
	var wallet models.Wallet

	err := s.dbService.DeleteByField(&wallet, "Address", address)
	if err != nil {
		fmt.Println("Failed to delete wallet with address:", address, err)
		return err
	}
	return nil
}
