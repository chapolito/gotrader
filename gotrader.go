package main

import (
	"os"
	"strconv"
	"fmt"

  exchange "github.com/preichenberger/go-coinbase-exchange"
)

//
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
// Go Trader
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
//
// Initialization
// 	1. what is the currentPrice?
// 	2. what are the buys?
// 	3. for each level below the currentPrice is there an open buy?
// 		a. if yes do nothing
//		b. if no, create a buy
//

var existingBuys, existingSells Orders
var totalBuys, totalSells, currentPrice, firstStep, lastStep, stepGap, totalSteps float64

var steps []float64
var btcIndex, usdIndex, ethIndex, stepsIndex int
var productId string

var accounts []exchange.Account
var client *exchange.Client

func main() {

	currentPrice = 0.0


	// firstStep = 179.0
	// lastStep = 399.0
	// stepGap = 5.0

	// min buy is .01 ETH
	// in order for stepGapSmall to be 1.0, we need $900 dedicated to it
	//
	// stepGapSmall = 2.0
	// stepGapMedium = 8.0
	// stepGapLarge = 16.0

	var firstStepErr, lastStepErr, stepGapErr error

	firstStep, firstStepErr = strconv.ParseFloat(os.Getenv("FIRST_STEP"), 64)
	lastStep, lastStepErr = strconv.ParseFloat(os.Getenv("LAST_STEP"), 64)
	stepGap, stepGapErr = strconv.ParseFloat(os.Getenv("STEP_GAP"), 64)

	if firstStepErr != nil || lastStepErr != nil || stepGapErr != nil {
	  println("ERROR parsing env variables as floats.\n")
	}

	totalSteps = (lastStep - firstStep) / stepGap

	for i := firstStep; i <= lastStep; i += stepGap {
		steps = append(steps, i)
	}

	// productId = "ETH-USD"

	productId = os.Getenv("PRODUCT_ID")
	fmt.Printf("ProductId: %v", productId)

	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client = exchange.NewClient(secret, key, passphrase)

	GetAccounts()

	MonitorExchange()

}
