package main

import (
	"fmt"

  exchange "github.com/preichenberger/go-coinbase-exchange"
)

//
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
// Get Fills
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
//
// Demonstrates profit and rate of profit
//

func CalculateProfit() {
  fmt.Printf("\n** CalculateProfit **\n")

  GetFills()
  GetStats()

  profit = filledSellsTotal - FilledBuysTotal + accounts[thisCoinAccountIndex].Balance * twentyFourHourAverage

  fmt.Printf("\nTotal %s Profit: $%f\n", productId, profit)
}

func GetStats() {
  fmt.Printf("\n** GetStats **\n")

  stats, err := client.GetStats(productId)
	if err != nil {
		println(err.Error())
	}
  twentyFourHourLow = stats.Low
  twentyFourHourHigh = stats.High

  // TODO This average calculation sucks, the peaks are not representative of the average, could be a single dip/peak throwing it off.
  twentyFourHourAverage = (twentyFourHourHigh + twentyFourHourLow) / 2

	//fmt.Printf("%s 24H || High: %f || Low: %f || Open: %f", productId, stats.High, stats.Low, stats.Open)
}

func GetFills() {

  fmt.Printf("\n** GetFills **\n")

  var fills []exchange.Fill
	cursorFills := client.ListFills()

	for cursorFills.HasMore {
		if err := cursorFills.NextPage(&fills); err != nil {
      println(err)
		} else {
      for _, f := range fills {
        if f.ProductId == productId {
          if f.Side == "buy" {
            FilledBuysTotal += f.Size * f.Price
          } else if f.Side == "sell" {
            filledSellsTotal += f.Size * f.Price
          }
        }
      }
    }
  }
}
