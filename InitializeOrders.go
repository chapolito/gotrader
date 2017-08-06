package main

import (
  "fmt"
  // exchange "github.com/preichenberger/go-coinbase-exchange"
)

func InitializeOrders() {

  // Figure out how many stops are between the firstStop and currentPrice
  var stopsUnderCurrentPrice int
  for a := range stops {
    if stops[a] <= currentPrice {
      stopsUnderCurrentPrice = a + 1
    }
  }

  // Match existing buys orders to stops. If no match create a buy order at that stop.
  for a := 0; a < stopsUnderCurrentPrice; a++ {
    if Contains(PricesExisting(existingBuys), stops[a]) {
      fmt.Printf("Buy existing at: %f\n", stops[a])
    } else {
      //CreateOrder("buy", stops[a], float64(int(((accounts[usdIndex].Balance / totalStops) / stops[a]) * 10000)) / 10000)

      // $1 Buys to test live
      CreateOrder("buy", stops[a], float64(int(((40.0 / totalStops) / stops[a]) * 10000)) / 10000)
    }
  }

  fmt.Printf("Current Price: %f\n", currentPrice)

  // Print out existingSells at stops
  for a := len(stops) - 1; a > stopsUnderCurrentPrice; a-- {
    if Contains(PricesExisting(existingSells), stops[a]) {
      fmt.Printf("Sell existing at: %f\n", stops[a])
    }
  }

  PrintCurrentState()
}

func Contains(s []float64, e float64) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func PricesExisting(o Orders) []float64 {
	var pricesWithBuys []float64
	for _, a := range o {
		pricesWithBuys = append(pricesWithBuys, a.Price)
	}
	return pricesWithBuys
}
