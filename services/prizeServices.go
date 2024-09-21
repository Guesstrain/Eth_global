package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/Guesstrain/ethglobal/database"
	"github.com/Guesstrain/ethglobal/models"
	"github.com/onflow/go-ethereum"
	"github.com/onflow/go-ethereum/accounts/abi"
	"github.com/onflow/go-ethereum/common"
	"github.com/onflow/go-ethereum/ethclient"
)

type PrizeService interface {
	InsertPrize(prize models.PrizeList) error
	UpdatePrize(prizeName string, updatedPrize models.PrizeList) error
	DistributePrize() []models.PrizeReward
}

type PrizeServiceImpl struct {
	dbService database.DatabaseService
}

// Contract ABI
var contractABI = `[{"inputs":[{"internalType":"address[]","name":"_add","type":"address[]"},{"internalType":"uint64","name":"min","type":"uint64"},{"internalType":"uint64","name":"max","type":"uint64"}],"name":"getRandomInRange","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"nonpayable","type":"function"}]`

// Contract Address
const contractAddress = "0x9AB786163fc09E3733e5E9133492eD47a814A029"

// Infura or Alchemy endpoint
const infuraURL = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"

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

func (s *PrizeServiceImpl) DistributePrize() []models.PrizeReward {
	startTime := time.Now().AddDate(-10, 0, 0) // 10 years ago from now
	endTime := time.Now()                      // Current time
	var prizes []models.Prize
	var prizeReward []models.PrizeReward
	const FirstReward float64 = 1000.00
	const SecondReward float64 = 500.00
	const ThirdReward float64 = 100.00
	const MinimumReward float64 = 10.00

	wallets, err := s.dbService.QueryWalletsByTimePeriod(startTime, endTime)
	if err != nil {
		fmt.Println("Failed to query all the wallets")
	}
	fmt.Println("wallets: ", wallets)

	for _, wallet := range wallets {
		randomNumber := CallSmartContract(wallet.Address)
		prizes = append(prizes, models.Prize{wallet.Address, randomNumber})
	}

	sort.Slice(prizes, func(i, j int) bool {
		return prizes[i].RandomNumber < prizes[j].RandomNumber
	})

	for i, wallet := range wallets {
		if i < 10 {
			prizeReward = append(prizeReward, models.PrizeReward{wallet.Address, FirstReward})
		} else if i < 100 {
			prizeReward = append(prizeReward, models.PrizeReward{wallet.Address, SecondReward})
		} else if i < 500 {
			prizeReward = append(prizeReward, models.PrizeReward{wallet.Address, ThirdReward})
		} else {
			prizeReward = append(prizeReward, models.PrizeReward{wallet.Address, MinimumReward})
		}
	}
	return prizeReward
}

func CallSmartContract(address string) uint64 {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	// Load the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Contract address as a common.Address type
	contractAddr := common.HexToAddress(contractAddress)

	// Define the input parameters for getRandomInRange
	addressList := []common.Address{
		common.HexToAddress(address),
	}
	min := big.NewInt(1)   // Minimum range value
	max := big.NewInt(100) // Maximum range value

	// Pack the arguments for getRandomInRange function
	input, err := parsedABI.Pack("getRandomInRange", addressList, min, max)
	if err != nil {
		log.Fatalf("Failed to pack arguments: %v", err)
	}

	// Create a call message
	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: input,
	}

	// Execute the call on the latest block
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	// Unpack the result
	var randomNumber uint64
	err = parsedABI.UnpackIntoInterface(&randomNumber, "getRandomInRange", result)
	if err != nil {
		log.Fatalf("Failed to unpack the result: %v", err)
	}

	// Print the random number generated
	fmt.Printf("Random Number: %d\n", randomNumber)
	return randomNumber
}
