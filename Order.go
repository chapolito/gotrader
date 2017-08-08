package main

import (
	"fmt"

	exchange "github.com/preichenberger/go-coinbase-exchange"
)

type Order struct {
	Type  string
	Id    string
	Price float64
	Size  float64
}

type Orders []Order

func CreateOrder(side string, price float64, size float64) error {

	thisOrder := exchange.Order{
		Price:     price,
		Size:      size,
		Side:      side,
		PostOnly:  true,
		ProductId: productId,
	}

	savedOrder, err := client.CreateOrder(&thisOrder)
	if err != nil {
		return err
	}

	fmt.Printf("%s order created for %f at $%f\n", savedOrder.Side, savedOrder.Size, savedOrder.Price)

	if side == "sell" {
		existingSells = append(existingSells, Order{"sell", savedOrder.Id, savedOrder.Size, savedOrder.Price})
	} else if side == "buy" {
		existingBuys = append(existingBuys, Order{"buy", savedOrder.Id, savedOrder.Size, savedOrder.Price})
	}

	return nil
}

// Get all Orders
func GetOrders() error {
	println("\n\n** GetOrders ** \n\n")
	var rawOrders []exchange.Order

	cursor := client.ListOrders()
	for cursor.HasMore {
		if err := cursor.NextPage(&rawOrders); err != nil {
			return err
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

	return nil
}

func InitializeOrders() {

	// Figure out how many stops are between the firstStop and currentPrice
	for a := range stops {
		if stops[a] <= currentPrice {
			stopsIndex = a + 1
		}
	}

	// Match existing buys orders to stops. If no match create a buy order at that stop.
	for a := 0; a < stopsIndex; a++ {
		if Contains(PricesExisting(existingBuys), stops[a]) {
			fmt.Printf("Buy existing at: %f\n", stops[a])
		} else {
			//CreateOrder("buy", stops[a], float64(int(((accounts[usdIndex].Balance / totalStops) / stops[a]) * 10000)) / 10000)

			// Minimium Buys (0.01 ETH) to test live
			CreateOrder("buy", stops[a], float64(int(((120.0/totalStops)/stops[a])*10000))/10000)
		}
	}

	fmt.Printf("Current Price: %f\n", currentPrice)

	// Print out existingSells at stops
	for a := len(stops) - 1; a > stopsIndex; a-- {
		if Contains(PricesExisting(existingSells), stops[a]) {
			fmt.Printf("Sell existing at: %f\n", stops[a])
		}
	}

	PrintCurrentState()
}

func PricesExisting(o Orders) []float64 {
	var pricesWithBuys []float64
	for _, a := range o {
		pricesWithBuys = append(pricesWithBuys, a.Price)
	}
	return pricesWithBuys
}

func ResetOrders() {
	// Clear out previously recorded orders
	existingSells = existingSells[:0]
	existingBuys = existingBuys[:0]
}
