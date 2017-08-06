package main

import (
	//"time"
	//"fmt"
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
//	1. If a sell happens create a buy at sell.Price-stopGap
// 	2. If a buy happens, create a sell at buy.Price+stopGap
//

var existingBuys, existingSells Orders
var totalBuys, totalSells, currentPrice, firstStop, lastStop, stopGap, totalStops float64
var stops []float64
var btcIndex, usdIndex, ethIndex int
var productId string

var accounts []exchange.Account
var client *exchange.Client

func main() {

	currentPrice = 0.0

	productId = "ETH-USD"
	firstStop = 139.0
	lastStop = 339.0
	stopGap = 5.0
	totalStops = (lastStop - firstStop) / stopGap

	for i := firstStop; i <= lastStop; i += stopGap {
		stops = append(stops, i)
	}

	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client = exchange.NewClient(secret, key, passphrase)
	// client = exchange.NewTestClient()

	GetAccounts()

	MonitorExchange()

}
