package main

import (
	"fmt"
)

// Get Accounts
func GetAccounts() {
	fmt.Printf("\n\n** GetAccounts **\n\n")

	var err error
	accounts, err = client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	// Figure out which account is which
	for i, a := range accounts {
		if a.Currency == "USD" {
			usdIndex = i
		}
		if a.Currency == "BTC" {
			btcIndex = i
		}
		if a.Currency == "ETH" {
			ethIndex = i
		}
		if a.Currency == "LTC" {
			ltcIndex = i
		}
	}

	// assign account index for current trading pair/coin
	if productId == "LTC-USD" {
		thisCoinAccountIndex = ltcIndex
	} else if productId == "ETH-USD" {
		thisCoinAccountIndex = ethIndex
	}
}
