package main

import (
	"fmt"
  exchange "github.com/preichenberger/go-coinbase-exchange"
)

func CreateOrder(side string, price float64, size float64) {

	if side == "sell" {
		price += stopGap
	}

	thisOrder := exchange.Order {
		Price: price,
		Size: size,
		Side: side,
		PostOnly: true,
		ProductId: ProductId,
	}

	savedOrder, err := client.CreateOrder(&thisOrder)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("%s order created for %f at $%f\n", savedOrder.Side, savedOrder.Size, savedOrder.Price)
		GetOrders()
	}
}
