package main

// Get Accounts
func GetAccounts() {
	var err error
	accounts, err = client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	// Print Balances
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
	}
}
