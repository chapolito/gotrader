package main

import (
	//"time"
	//"fmt"
	//"os"
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
//	1. If a sell happens create a buy at sell.Price-stopGap
// 	2. If a buy happens, create a sell at buy.Price+stopGap
//

var existingBuys, existingSells Orders
var totalBuys, totalSells, currentPrice, currentFakePrice, firstStop, lastStop, stopGap, totalStops float64
var stops []float64
var btcIndex, usdIndex int

var accounts []exchange.Account
var client *exchange.Client

func main() {

	currentFakePrice = 1.0
	firstStop = 500.0
	lastStop = 1500.0
	stopGap = 25.0
	totalStops = (lastStop - firstStop) / stopGap

	for i := firstStop; i <= lastStop; i += stopGap {
		stops = append(stops, i)
	}

	// secret := os.Getenv("TEST_COINBASE_SECRET")
	// key := os.Getenv("TEST_COINBASE_KEY")
	// passphrase := os.Getenv("TEST_COINBASE_PASSPHRASE")
	// fmt.Printf("secret: " + secret + "\n key: " + key + "\n passphrase: " + passphrase + "\n")

	client = exchange.NewTestClient()

	GetAccounts()

	GetOrders()

	MonitorExchange()

}
