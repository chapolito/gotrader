package main

// import (
//   exchange "github.com/preichenberger/go-coinbase-exchange"
// )


// Get Accounts
func GetAccounts() {
  var err error
  accounts, err = client.GetAccounts()
  if err != nil {
    println(err.Error())
  }

  // Print Balances
  for i, a := range accounts {
    println(a.Balance)
    if a.Currency == "USD" {
      usdIndex = i
    }
    if a.Currency == "BTC" {
      btcIndex = i
    }
  }
}
