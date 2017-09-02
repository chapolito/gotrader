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
// Websocket Monitoring
//	1. If a sell happens create a buy at sell.Price-stepGap
// 	2. If a buy happens, create a sell at buy.Price+stepGap
//	3. If the price increases to a new high, keep buying
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

	productId = "ETH-USD"
	// firstStep = 239.0
	// lastStep = 439.0
	// stepGap = 5.00
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

	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client = exchange.NewClient(secret, key, passphrase)

	GetAccounts()

	MonitorExchange()

}
