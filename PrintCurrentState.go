package main

import (
	"fmt"
)

func PrintCurrentState() {

	// Print out existingBuys at steps
	for a := 0; a < nextStepIndex; a++ {
		if Contains(PricesExisting(existingBuys), steps[a]) {
			fmt.Printf("Buy existing at: %f\n", steps[a])
		}
	}

	fmt.Printf("---\nCurrent Price: %f\n---\n", currentPrice)

	// Print out existingSells at steps
	for a := nextStepIndex; a <= len(steps) - 1; a++ {
		if Contains(PricesExisting(existingSells), steps[a]) {
			fmt.Printf("Sell existing at: %f\n", steps[a])
		}
	}

	// Print Balances
	fmt.Printf("\nUSD - Balance: $%f, Hold: $%f, Available: $%f\n", accounts[usdIndex].Balance, accounts[usdIndex].Hold, accounts[usdIndex].Available)
	fmt.Printf("\nBTC - Balance: $%f, Hold: $%f, Available: $%f\n", accounts[btcIndex].Balance, accounts[btcIndex].Hold, accounts[btcIndex].Available)
	fmt.Printf("\nETH - Balance: $%f, Hold: $%f, Available: $%f\n", accounts[ethIndex].Balance, accounts[ethIndex].Hold, accounts[ethIndex].Available)
	fmt.Printf("\nLTC - Balance: $%f, Hold: $%f, Available: $%f\n", accounts[ltcIndex].Balance, accounts[ltcIndex].Hold, accounts[ltcIndex].Available)

	// How to get current price for each currency?
	//fmt.Printf("Total worth in USD: %f\n\n", accounts[usdIndex].Balance + accounts[ethIndex].Balance * currentPriceETH + accounts[ltcIndex].Balance * currentPriceLTC)

	// totalBuys is similar to accounts[usdIndex].Hold, but it takes into consideration partially filled orders
	totalBuys = TotalBuys()
	fmt.Printf("\nBuys: %d\nCost: $%f\n", len(existingBuys), totalBuys)

	// totalSells is similar to accounts[ethIndex].Hold, but it takes into consideration partially filled orders
	totalSells = TotalSells()
	fmt.Printf("\nSells: %d\nCost: $%f\n", len(existingSells), totalSells)
}

func TotalBuys() float64 {
	totalBuys = 0.0
	// Compute total cost of all buys.
	for _, e := range existingBuys {
		totalBuys += e.Price * e.Size
	}
	return totalBuys
}

func TotalSells() float64 {
	totalSells = 0.0
	// Compute total cost of all sells.
	for _, e := range existingSells {
		totalSells += e.Price * e.Size
	}
	return totalSells
}
