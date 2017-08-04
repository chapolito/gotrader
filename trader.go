package main

import (
	"time"
	"fmt"
	"os"
  exchange "github.com/preichenberger/go-coinbase-exchange"
	ws "github.com/gorilla/websocket"
)


func main() {

	secret := os.Getenv("TEST_COINBASE_SECRET")
	key := os.Getenv("TEST_COINBASE_KEY")
	passphrase := os.Getenv("TEST_COINBASE_PASSPHRASE")
	fmt.Printf("secret: " + secret + "\n key: " + key + "\n passphrase: " + passphrase + "\n")

	// or unsafe hardcode way
	// secret = "exposedsecret"
	// key = "exposedkey"
	// passphrase = "exposedpassphrase"

	// Initialization
	// 	1. what is the currentPrice?
	// 	2. what are the buys?
	// 	3. for each level below the currentPrice is there an open buy?
	// 		a. if yes do nothing
	//		b. if no, create a buy
	//
	//
	// Monitoring
	//	1. If a sell happens, set buy at sell.Price-stopGap
	// 	2. If a buy happens, create a sell at buy.Price+stopGap


	client := exchange.NewTestClient()

	var stops []float64
	currentFakePrice := 1.0
	var currentPrice float64
	firstStop := 500.0
	lastStop := 1500.0
	stopGap := 25.0
	totalStops := (lastStop - firstStop) / stopGap

	for i := firstStop; i <= lastStop; i += stopGap {
		stops = append(stops, i)
	}

	//printSlice(stops)

	// var total float64
	// var usd float64
	// var btc float64

	var existingBuys, existingSells [][]interface{}
	var totalBuys, totalSells float64



	// Get Accounts
	accounts, err := client.GetAccounts()
  if err != nil {
    println(err.Error())
  }

	// Print Balances
	var btcIndex, usdIndex int
	for i, a := range accounts {
		println(a.Balance)
    if a.Currency == "USD" {
			usdIndex = i
		}
		if a.Currency == "BTC" {
			btcIndex = i
		}
  }

	fmt.Printf("USD - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[usdIndex].Balance, accounts[usdIndex].Hold, accounts[usdIndex].Available)
	fmt.Printf("BTC - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[btcIndex].Balance, accounts[btcIndex].Hold, accounts[btcIndex].Available)



	// Get all Orders
	var orders []exchange.Order
	cursor := client.ListOrders()

	for cursor.HasMore {
		if err := cursor.NextPage(&orders); err != nil {
			println(err.Error())
			return
		}

		for _, o := range orders {
			if o.Type == "limit" && o.ProductId == "BTC-USD" {
				if o.Side == "sell" {
					sell := []interface{}{o.Price, o.Size, o.Id}
					existingSells = append(existingSells, sell)
				} else if o.Side == "buy" {
					buy := []interface{}{o.Price, o.Size, o.Id}
					existingBuys = append(existingBuys, buy)
				}
			}
		}
	}

	// Compute total cost of all buys.
	// Note: I think this should simply be accounts[usdIndex].Hold
	for i := 0; i < len(existingBuys); i++ {
		price := existingBuys[i][0].(float64)
		size := existingBuys[i][1].(float64)
		totalBuys += price * size
	}
	fmt.Printf("Buys: %d\nCost: $%f\n\n", len(existingBuys), totalBuys)

	// Compute total cost of all sells
	// Note: I think this should simply be accounts[btcIndex].Hold
	for i := 0; i < len(existingSells); i++ {
		price := existingSells[i][0].(float64)
		size := existingSells[i][1].(float64)
		totalSells += price * size
	}
	fmt.Printf("Sells: %d\nCost: $%f\n\n", len(existingSells), totalSells)

	// Figure out how many stops are between the firstStop and currentPrice
	var stopsUnderCurrentPrice int
	for a := range stops {
		if stops[a] <= currentFakePrice {
			stopsUnderCurrentPrice = a + 1
		}
	}

	// Match existing buys orders to stops. If no match create a buy order at that stop.
	for a := 0; a < stopsUnderCurrentPrice; a++ {

		if contains(pricesExisting(existingBuys), stops[a]) {
			fmt.Printf("Buy existing at: %f\n", stops[a])
		} else {
			//createBuyOrder(stops[a], 1.1, client)

			size := float64(int(((accounts[usdIndex].Balance / totalStops) / stops[a]) * 10000)) / 10000

			order := exchange.Order {
				Price: stops[a],
				Size: size,
				Side: "buy",
				PostOnly: true,
				ProductId: "BTC-USD",
			}

			savedOrder, err := client.CreateOrder(&order)
			if err != nil {
				println(err.Error())
			} else {
				fmt.Printf("Buy Order Created for %f at $%f\n", size, stops[a])
				buy := []interface{}{stops[a], size, savedOrder.Id}
				existingBuys = append(existingBuys, buy)
			}
		}
	}

	fmt.Printf("Current Price: %f\n", currentFakePrice)

	// Print out existingSells at stops
	for a := len(stops) - 1; a > stopsUnderCurrentPrice; a-- {
		if contains(pricesExisting(existingSells), stops[a]) {
			fmt.Printf("Sell existing at: %f\n", stops[a])
		}
	}


	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
	    for {
	       select {
	        case <- ticker.C:
						// Get fills
						var fills []exchange.Fill
						cursorFills := client.ListFills()

						for cursorFills.HasMore {
							if err := cursorFills.NextPage(&fills); err != nil {
								println(err.Error())
								return
							}

							for _, o := range fills {
								//fmt.Printf("o: %v\n", o)
								for i, b := range existingBuys {
									//fmt.Printf("buy:  %v\nfill: %v\n\n", b[2], o.FillId)

									// Should check here if the fill is completely filled
									if o.FillId == b[2] {
										fmt.Printf("Buy happened for %f at $%f\n", b[1], b[0])
										// Remove fill from existingBuys

										fmt.Printf("%v\n\n", existingBuys)
										existingBuys = remove(existingBuys, i)
									}
								}
								for _, s := range existingSells {
									//fmt.Printf("s: %v\n", s)
									if o.FillId == s[2] {
										fmt.Printf("Sell happened for %f at $%f\n", s[1], s[0])
									}
								}
							}
						}
	        case <- quit:
	            ticker.Stop()
	            return
	        }
	    }
	 }()



	// Get Orders
	// var orders []exchange.Order
	// cursor := client.ListOrders()
	//
	// for cursor.HasMore {
	// 	if err := cursor.NextPage(&orders); err != nil {
	// 		println(err.Error())
	// 		return
	// 	}
	//
	// 	for _, o := range orders {
	// 		if o.Type == "limit" && o.ProductId == "BTC-USD" {
	// 			//println(o.ProductId + "  " + o.Type + "  " + o.Side)
	// 			var value = o.Price * o.Size
	//
	// 			if o.Side == "sell" {
	// 				btc -= value
	//
	// 			} else if o.Side == "buy" {
	//
	// 				usd -= value
	// 			}
	// 		}
	// 	}
	// }
	// getOrders(btc, usd, client)





	//var totalDollars int64 = int64(total)
	//mt.Printf("Trader NW: $%f", total)
	//fmt.Printf("Liquid USD: $%f\nLiquid BTC: $%f\nTrader NW: $%f\n", usd, btc, total)



	// Websocket
  var wsDialer ws.Dialer
  wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
  if err != nil {
    println(err.Error())
  }

  subscribe := map[string]string{
    "type": "subscribe",
    "product_id": "ETH-USD",
  }
  if err := wsConn.WriteJSON(subscribe); err != nil {
    println(err.Error())
  }

  message:= exchange.Message{}
  for true {
    if err := wsConn.ReadJSON(&message); err != nil {
      println(err.Error())
      break
    }

    if message.Type == "match" {
			currentPrice = message.Price
			fmt.Printf("ETH Price: $%f\n\n", currentPrice)
    }
  }




	// Get Ledger
	// var ledger []exchange.LedgerEntry
	//
  // for _, a := range accounts {
  //   cursor := client.ListAccountLedger(a.Id)
  //   for cursor.HasMore {
  //     if err := cursor.NextPage(&ledger); err != nil {
  //       println(err.Error())
  //     }
	//
  //     for _, e := range ledger {
  //     	fmt.Printf("\n" + e.Type + "\n")
	// 			println(e.Amount)
	// 			println(e.Balance)
  // 		}
  // 	}
	// }
}

// func createBuyOrder(price float64, size float64, client interface{}) {
//
// }

// func printSlice(s []float64) {
// 	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
// }

func contains(s []float64, e float64) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func pricesExisting(e [][]interface{}) []float64 {
	var pricesWithBuys []float64
	for a := range e {
		var price float64
		price = e[a][0].(float64)
		pricesWithBuys = append(pricesWithBuys, price)
	}
	return pricesWithBuys
}

func remove(s [][]interface{}, i int) [][]interface{} {
	// var i int
	// for a = 0, b := range s {
	// 	if b[a] == id {
	// 		i = a
	// 	}
	// }

	fmt.Printf("%v\n\n", s)

	s[len(s)-1], s[i] = s[i], s[len(s)-1]

	return s[:len(s)-1]
}
// func getOrders(btc int, usd int, client) {
//
// }
