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

  fmt.Printf("Total %s Profit: $%f", productId, sellsTotal - buysTotal)
}
