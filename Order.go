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
	fmt.Printf("\n** CreateOrder **\n")

	thisOrder := exchange.Order{
		Price:     price,
		Size:      size,
		Side:      side,
		PostOnly:  true,
		ProductId: productId,
	}

	savedOrder, err := client.CreateOrder(&thisOrder)

	if err != nil {
		fmt.Printf("\nCreateOrder error: %v\n", err)
		return err
	} else {
		fmt.Printf("\n%s order created for %f at $%f\n", savedOrder.Side, savedOrder.Size, savedOrder.Price)

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
	fmt.Printf("\n** CancelOrder **\n")

	err := client.CancelOrder(id)

	if err != nil {
		fmt.Printf("\nCancelOrder error: %v\n", err)
		return err
	} else {
		return nil
	}
}


func GetOrders() error {

	fmt.Printf("\n** GetOrders **\n")

	// Zero out previously recorded orders
	existingSells = existingSells[:0]
	existingBuys = existingBuys[:0]

	// Get all Orders
	var rawOrders []exchange.Order
	cursor := client.ListOrders()
	for cursor.HasMore {
		if err := cursor.NextPage(&rawOrders); err != nil {
			return err
		}

		// Record and categorize orders
		for _, o := range rawOrders {
			if o.Type == "limit" && o.ProductId == productId {
				if o.Side == "sell" {
					existingSells = append(existingSells, Order{"sell", o.Id, o.Price, o.Size})
				} else if o.Side == "buy" {
					existingBuys = append(existingBuys, Order{"buy", o.Id, o.Price, o.Size})
					// Should I cleanse Orders here?
				}
			}
		}
	}

	return nil
}

func InitializeOrders() {

	fmt.Printf("\n** InitializeOrders **\n")

	// Set existingBuys and existingSells
	GetOrders()

	// Prune orders out of existingBuys
	PruneBuys()

	// Create any missing buys that there are steps for
	CreateMissingBuys()

	// Cancel and recreate everything in existingBuys
	CompoundOrders()

	// Run GetOrders again to refresh existingBuys and existingSells
	GetOrders()

	PrintCurrentState()
}

func PruneBuys() {
	fmt.Printf("\n** PruneBuys **\n")

	// Should run when:
	//   1. The currentPrice passes above or below a step

	// Loop through existing buys
	// check if price is still within steps[], if it is not, cancel and remove from existingBuys
	for _, o := range existingBuys {
		if !Contains(steps, o.Price) {
			CancelOrder(o.Id)
		}
	}

	// Run GetOrders to refresh existingBuys (and existingSells)
	GetOrders()
}

func CompoundOrders() {
	fmt.Printf("\n** CompoundOrders **\n")

	// Loop through existing orders, cancel them, then recreate
	// HowMuchToBuy() ensures any profits get pulled in to start compounding.
	for _, o := range existingBuys {
		CancelOrder(o.Id)
		CreateOrder("buy", o.Price, HowMuchToBuy(o.Price))
	}
}

func CreateMissingBuys() {
	fmt.Printf("\n** CreateMissingBuys **\n")

	for a := 0; a < nextStepIndex; a++ {

		// Is there NOT a sell at current step + 1 AND is there NOT a buy at current step?
		if !Contains(PricesExisting(existingSells), steps[a + 1]) && !Contains(PricesExisting(existingBuys), steps[a]) {
			fmt.Printf("\n** Buy needs to be created! **\n")
			CreateOrder("buy", steps[a], HowMuchToBuy(steps[a]))
		}
	}
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
