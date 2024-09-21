package services

import (
	"fmt"
	"time"

	"github.com/Guesstrain/ethglobal/database"
	"github.com/Guesstrain/ethglobal/models"
)

type PrizeService interface {
	InsertPrize(prize models.PrizeList) error
	UpdatePrize(prizeName string, updatedPrize models.PrizeList) error
	DistributePrize() []models.Prize
}

type PrizeServiceImpl struct {
	dbService database.DatabaseService
}

func NewPrizeService(dbService database.DatabaseService) PrizeService {
	return &PrizeServiceImpl{dbService: dbService}
}

func (s *PrizeServiceImpl) InsertPrize(prize models.PrizeList) error {

	err := s.dbService.Insert(&prize)
	if err != nil {
		fmt.Println("Failed to insert prize:", err)
	}
	return err
}

func (s *PrizeServiceImpl) UpdatePrize(prizeName string, updatedPrize models.PrizeList) error {
	// Find the prize by ID
	var existingPrize models.PrizeList
	if err := s.dbService.SelectByField(&existingPrize, "prize_name", prizeName); err != nil {
		fmt.Println("Failed to find prize:", err)
		return err
	}

	// Update the fields
	existingPrize.PrizeName = prizeName
	existingPrize.Amount = updatedPrize.Amount
	existingPrize.Probability = updatedPrize.Probability

	// Save the updated prize
	if err := s.dbService.UpdateByStruct(&existingPrize, "prize_name", prizeName, &existingPrize); err != nil {
		fmt.Println("Failed to update prize:", err)
		return err
	}

	return nil
}

func (s *PrizeServiceImpl) DistributePrize() []models.Prize {
	startTime := time.Now().AddDate(-10, 0, 0) // 10 years ago from now
	endTime := time.Now()                      // Current time
	var prizes []models.Prize

	wallets, err := s.dbService.QueryWalletsByTimePeriod(startTime, endTime)
	if err != nil {
		fmt.Println("Failed to query all the wallets")
	}
	fmt.Println("wallets: ", wallets)

	//calculate the odds
	//set the threshold --select all the staking money, calculate the prize pool
	for _, wallet := range wallets {
		prizes = append(prizes, models.Prize{wallet.Address, 11})
	}

	return prizes
}
