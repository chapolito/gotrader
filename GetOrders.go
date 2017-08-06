package main

import (
  "fmt"
  exchange "github.com/preichenberger/go-coinbase-exchange"
)

// Get all Orders
func GetOrders() {
  var rawOrders []exchange.Order
  cursor := client.ListOrders()

  for cursor.HasMore {
    if err := cursor.NextPage(&rawOrders); err != nil {
      println(err.Error())
      return
    }

    for _, o := range rawOrders {
      if o.Type == "limit" && o.ProductId == "BTC-USD" {
        if o.Side == "sell" {
          existingSells = append(existingSells, Order{"sell", o.Id, o.Price, o.Size})
        } else if o.Side == "buy" {
          existingBuys = append(existingBuys, Order{"buy", o.Id, o.Price, o.Size})
        }
      }
    }
  }

  // Figure out how many stops are between the firstStop and currentPrice
	var stopsUnderCurrentPrice int
	for a := range stops {
		if stops[a] <= currentFakePrice {
			stopsUnderCurrentPrice = a + 1
		}
	}

	// Match existing buys orders to stops. If no match create a buy order at that stop.
	for a := 0; a < stopsUnderCurrentPrice; a++ {
		if contains(pricesExisting(existingBuys), stops[a]) {
			fmt.Printf("Buy existing at: %f\n", stops[a])
		} else {
			CreateBuyOrder(stops[a], float64(int(((accounts[usdIndex].Balance / totalStops) / stops[a]) * 10000)) / 10000)
		}
	}

	fmt.Printf("Current Price: %f\n", currentFakePrice)

	// Print out existingSells at stops
	for a := len(stops) - 1; a > stopsUnderCurrentPrice; a-- {
		if contains(pricesExisting(existingSells), stops[a]) {
			fmt.Printf("Sell existing at: %f\n", stops[a])
		}
	}
}
