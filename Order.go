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
	fmt.Printf("\n\n** CreateOrder **\n")

	thisOrder := exchange.Order{
		Price:     price,
		Size:      size,
		Side:      side,
		PostOnly:  true,
		ProductId: productId,
	}

	savedOrder, err := client.CreateOrder(&thisOrder)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	} else {
		fmt.Printf("%s order created for %f at $%f", savedOrder.Side, savedOrder.Size, savedOrder.Price)

		// Update Lists of Buys / Sells
		if side == "sell" {
			existingSells = append(existingSells, Order{"sell", savedOrder.Id, savedOrder.Price, savedOrder.Size})
		} else if side == "buy" {
			existingBuys = append(existingBuys, Order{"buy", savedOrder.Id, savedOrder.Price, savedOrder.Size})
		}

		return nil
	}
}

func CancelOrder(id string) error {
	fmt.Printf("\n** CancelOrder **\n\n")

	err := client.CancelOrder(id)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	} else {
		return nil
	}
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

	// Loop through existing orders, cancel them, then recreate them
	// This ensures any profits get pulled in so they start compounding.
	for _, o := range existingBuys {
		CancelOrder(o.Id)
		CreateOrder("buy", o.Price, HowMuchToBuy(o.Price))
	}

	ResetOrders()
	GetOrders()

	// Match existing buys orders to steps. If no match create a buy order at that step.
	//
	// NOTE stepsIndex always points to the next higher step
	for a := 0; a < stepsIndex; a++ {
		if Contains(PricesExisting(existingBuys), steps[a]) {
			fmt.Printf("Buy existing at: %f\n", steps[a])
		} else if !Contains(PricesExisting(existingSells), steps[a + 1]) {
			CreateOrder("buy", steps[a], HowMuchToBuy(steps[a]))
		}
	}

	fmt.Printf("Current Price: %f\n", currentPrice)

	// Print out existingSells at steps
	for a := len(steps) - 1; a >= stepsIndex; a-- {
		if Contains(PricesExisting(existingSells), steps[a]) {
			fmt.Printf("Sell existing at: %f\n", steps[a])
		}
	}

	PrintCurrentState()
}

func HowMuchToBuy(price float64) float64 {
	// Minimium Buys are 0.01 ETH/LTC/BTC
	// Dividing by 2 to split between ETH and LTC
	return float64(int(((accounts[usdIndex].Balance / (totalSteps - float64(len(existingSells)))) / price) * 10000)) / 10000 / 2
}

func PricesExisting(o Orders) []float64 {
	var pricesWithBuys []float64
	for _, a := range o {
		pricesWithBuys = append(pricesWithBuys, a.Price)
	}
	return pricesWithBuys
}

func ResetOrders() {
	println("\n\n** ResetOrders **\n\n")
	// Clear out previously recorded orders
	existingSells = existingSells[:0]
	existingBuys = existingBuys[:0]
}
