package main

import (
	"os"
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
	firstStep = 179.0
	lastStep = 379.0
	stepGap = 2.50
	totalSteps = (lastStep - firstStep) / stepGap

	for i := firstStep; i <= lastStep; i += stepGap {
		steps = append(steps, i)
	}

	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client = exchange.NewClient(secret, key, passphrase)
	// client = exchange.NewTestClient()

	GetAccounts()

	//CreateOrder("sell", 313.87, 0.01)

	MonitorExchange()

}
