package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/Guesstrain/ethglobal/database"
	"github.com/Guesstrain/ethglobal/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func CallSmartContract() {
	// Connect to an Ethereum node (Infura or local node)
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Load the private key
	privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Derive the sender's public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce for the sender
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Set the gas price and limit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Define the contract address and ABI
	contractAddress := common.HexToAddress("YOUR_CONTRACT_ADDRESS")
	contractABI, err := abi.JSON(strings.NewReader(`YOUR_CONTRACT_ABI`)) // Use your contract's ABI as JSON string
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Prepare the function call
	data, err := contractABI.Pack("someFunction", param1, param2) // Replace with your function name and parameters
	if err != nil {
		log.Fatalf("Failed to pack contract function: %v", err)
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), 300000, gasPrice, data)

	// Sign the transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
}
