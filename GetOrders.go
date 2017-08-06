package main

import (
  exchange "github.com/preichenberger/go-coinbase-exchange"
)

// Get all Orders
func GetOrders() {
  println("\n\n** GetOrders ** \n\n")
  var rawOrders []exchange.Order
  cursor := client.ListOrders()

  for cursor.HasMore {
    if err := cursor.NextPage(&rawOrders); err != nil {
      println(err.Error())
      return
    }

    for _, o := range rawOrders {
      if o.Type == "limit" && o.ProductId == productId {
        if o.Side == "sell" {
          existingSells = append(existingSells, Order{"sell", o.Id, o.Price, o.Size})
        } else if o.Side == "buy" {
          existingBuys = append(existingBuys, Order{"buy", o.Id, o.Price, o.Size})
        }
      }
    }
  }
}

func ResetOrders() {
  // Clear out previously recorded orders
  existingSells = existingSells[:0]
  existingBuys = existingBuys[:0]
}
