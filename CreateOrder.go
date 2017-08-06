package main

import (
	"fmt"
  exchange "github.com/preichenberger/go-coinbase-exchange"
)

// Todo: Combine into one CreateOrder func

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
		fmt.Printf("Buy Order Created for %f at $%f\n", savedOrder.Size, savedOrder.Price)

		// Rerun GetOrders?
		// existingBuys = append(existingBuys, Order{"buy", savedOrder.Id, price, size})
		GetOrders()
	}
}

func CreateSellOrder(price float64, size float64) {

	sellPrice := price + stopGap

	thisOrder := exchange.Order {
		Price: sellPrice,
		Size: size,
		Side: "sell",
		PostOnly: true,
		ProductId: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&thisOrder)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Sell Order Created for %f at $%f\n", savedOrder.Size, savedOrder.Price)

		// Rerun GetOrders?
		// existingSells = append(existingSells, Order{"sell", savedOrder.Id, sellPrice, size})
		GetOrders()
	}
}
