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
	println("\n\n** CreateOrder ** \n\n")

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

	// Figure out how many steps are between the firstStep and currentPrice
	for a := range steps {
		if steps[a] <= currentPrice {
			stepsIndex = a + 1
		}
	}

	// Match existing buys orders to steps. If no match create a buy order at that step.
	for a := 0; a < stepsIndex; a++ {
		if Contains(PricesExisting(existingBuys), steps[a]) {
			fmt.Printf("Buy existing at: %f\n", steps[a])
		} else {
			CreateOrder("buy", steps[a], HowMuchToBuy(steps[a]))
		}
	}

	fmt.Printf("Current Price: %f\n", currentPrice)

	// Print out existingSells at steps
	for a := len(steps) - 1; a > stepsIndex; a-- {
		if Contains(PricesExisting(existingSells), steps[a]) {
			fmt.Printf("Sell existing at: %f\n", steps[a])
		}
	}

	PrintCurrentState()
}

func HowMuchToBuy(price float64) float64 {
	return float64(int(((accounts[usdIndex].Balance / totalSteps) / price) * 10000)) / 10000

	// Minimium Buys (0.01 ETH) to test live
	//return float64(int(((120.0/totalSteps)/price)*10000))/10000
}

func PricesExisting(o Orders) []float64 {
	var pricesWithBuys []float64
	for _, a := range o {
		pricesWithBuys = append(pricesWithBuys, a.Price)
	}
	return pricesWithBuys
}

func ResetOrders() {
	println("\n\n** ResetOrders ** \n\n")
	// Clear out previously recorded orders
	existingSells = existingSells[:0]
	existingBuys = existingBuys[:0]
}
