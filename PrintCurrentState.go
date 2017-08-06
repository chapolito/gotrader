package main

import (
	"fmt"
)

func PrintCurrentState() {

	// Print Balances
	fmt.Printf("USD - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[usdIndex].Balance, accounts[usdIndex].Hold, accounts[usdIndex].Available)
	//fmt.Printf("BTC - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[btcIndex].Balance, accounts[btcIndex].Hold, accounts[btcIndex].Available)
	fmt.Printf("ETH - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[ethIndex].Balance, accounts[ethIndex].Hold, accounts[ethIndex].Available)

	// totalBuys is similar to accounts[usdIndex].Hold, but it takes into consideration partially filled orders
	totalBuys = TotalBuys()
	fmt.Printf("Buys: %d\nCost: $%f\n\n", len(existingBuys), totalBuys)

	// totalSells is similar to accounts[btcIndex].Hold, but it takes into consideration partially filled orders
	totalSells = TotalSells()
	fmt.Printf("Sells: %d\nCost: $%f\n\n", len(existingSells), totalSells)
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
