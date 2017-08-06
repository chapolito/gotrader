package main

import (
  "fmt"
)

func PrintCurrentState() {

  // Print Balances
  fmt.Printf("USD - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[usdIndex].Balance, accounts[usdIndex].Hold, accounts[usdIndex].Available)
  fmt.Printf("BTC - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[btcIndex].Balance, accounts[btcIndex].Hold, accounts[btcIndex].Available)

  // Compute total cost of all buys.
  // totalBuys is similar to accounts[usdIndex].Hold, but it takes into consideration partially filled orders
  for _, e := range existingBuys {
    totalBuys += e.Price * e.Size
  }
  fmt.Printf("Buys: %d\nCost: $%f\n\n", len(existingBuys), totalBuys)

  // Compute total cost of all sells
  // totalSells is similar to accounts[btcIndex].Hold, but it takes into consideration partially filled orders
  for _, e := range existingSells {
    totalSells += e.Price * e.Size
  }
  fmt.Printf("Sells: %d\nCost: $%f\n\n", len(existingSells), totalSells)
}
