package main

import (
	"os"
	"strconv"

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
var totalBuys, totalSells, currentPrice, firstStep, lastStep, stepGap, totalSteps, holdSteps float64

var steps []float64
var btcIndex, usdIndex, ethIndex, ltcIndex, stepsIndex int
var productId string

var accounts []exchange.Account
var client *exchange.Client

func main() {

	currentPrice = 0.0
	holdSteps = 14.0

	productId = os.Getenv("PRODUCT_ID")

	var stepGapErr error
	// var firstStepErr, lastStepErr, stepGapErr error
	// firstStep, firstStepErr = strconv.ParseFloat(os.Getenv("FIRST_STEP"), 64)
	// lastStep, lastStepErr = strconv.ParseFloat(os.Getenv("LAST_STEP"), 64)
	stepGap, stepGapErr = strconv.ParseFloat(os.Getenv("STEP_GAP"), 64)

	if stepGapErr != nil {
	  println("ERROR parsing env vars as floats.\n")
	}

	// if firstStepErr != nil || lastStepErr != nil || stepGapErr != nil {
	//   println("ERROR parsing env vars as floats.\n")
	// }

	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")
	client = exchange.NewClient(secret, key, passphrase)

	GetFills()

	GetAccounts()

	MonitorExchange()

}
