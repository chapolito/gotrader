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

func GetStats() {
  stats, err := client.GetStats(productId)
	if err != nil {
		println(err.Error())
	}
  twentyFourHourLow = stats.Low
  twentyFourHourHigh = stats.High

  // This average calculation sucks, the peaks are not representative of the average, could be a single dip/peak throwing it off.
  twentyFourHourAverage = (twentyFourHourHigh + twentyFourHourLow) / 2

	//fmt.Printf("%s 24H || High: %f || Low: %f || Open: %f", productId, stats.High, stats.Low, stats.Open)
}

func GetFills() {

  var sellsTotal, buysTotal float64

	var fills []exchange.Fill
	cursorFills := client.ListFills()

	for cursorFills.HasMore {
		if err := cursorFills.NextPage(&fills); err != nil {
      println(err)
		} else {
      for _, f := range fills {
        if f.ProductId == productId {
          if f.Side == "buy" {
            buysTotal += f.Size * f.Price
          } else if f.Side == "sell" {
            sellsTotal += f.Size * f.Price
          }
        }
      }
    }
  }

  GetStats()

  profit = sellsTotal - buysTotal + accounts[ltcIndex].Balance * twentyFourHourAverage

  fmt.Printf("Total %s Profit: $%f\n\n", productId, profit)
}
