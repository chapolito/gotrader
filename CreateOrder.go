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
		ProductId: productId,
	}

	savedOrder, err := client.CreateOrder(&thisOrder)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("%s order created for %f at $%f\n", savedOrder.Side, savedOrder.Size, savedOrder.Price)

		if side == "sell" {
			existingSells = append(existingSells, Order{"sell", savedOrder.Id, savedOrder.Size, savedOrder.Price})
		} else if side == "buy" {
			existingBuys = append(existingBuys, Order{"buy", savedOrder.Id, savedOrder.Size, savedOrder.Price})
		}
	}
}
