package main

import (
	"fmt"
  exchange "github.com/preichenberger/go-coinbase-exchange"
)

func CreateBuyOrder(price float64, size float64) {

	thisOrder := exchange.Order {
		Price: price,
		Size: size,
		Side: "buy",
		PostOnly: true,
		ProductId: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&thisOrder)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Buy Order Created for %f at $%f\n", size, price)
		existingBuys = append(existingBuys, Order{"buy", savedOrder.Id, price, size})
	}
}
