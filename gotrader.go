package main

import (
	"time"
	"fmt"
	//"os"
  exchange "github.com/preichenberger/go-coinbase-exchange"
	ws "github.com/gorilla/websocket"
)

//
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
// Go Trader
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
//
// Initialization
// 	1. what is the currentPrice?
// 	2. what are the buys?
// 	3. for each level below the currentPrice is there an open buy?
// 		a. if yes do nothing
//		b. if no, create a buy
//
//
// Monitoring
//	1. If a sell happens:
//		a. set buy at sell.Price-stopGap
//		b. remove sell from existingSells
// 	2. If a buy happens, create a sell at buy.Price+stopGap
//

var existingBuys, existingSells Orders
var totalBuys, totalSells, currentPrice, currentFakePrice, firstStop, lastStop, stopGap, totalStops float64
var stops []float64
var btcIndex, usdIndex int

var accounts []exchange.Account
var client *exchange.Client

func main() {

	currentFakePrice = 1.0
	firstStop = 500.0
	lastStop = 1500.0
	stopGap = 25.0
	totalStops = (lastStop - firstStop) / stopGap

	// secret := os.Getenv("TEST_COINBASE_SECRET")
	// key := os.Getenv("TEST_COINBASE_KEY")
	// passphrase := os.Getenv("TEST_COINBASE_PASSPHRASE")
	// fmt.Printf("secret: " + secret + "\n key: " + key + "\n passphrase: " + passphrase + "\n")

	client = exchange.NewTestClient()

	for i := firstStop; i <= lastStop; i += stopGap {
		stops = append(stops, i)
	}

	GetAccounts()

	fmt.Printf("USD - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[usdIndex].Balance, accounts[usdIndex].Hold, accounts[usdIndex].Available)
	fmt.Printf("BTC - Balance: $%f, Hold: $%f, Available: $%f\n\n", accounts[btcIndex].Balance, accounts[btcIndex].Hold, accounts[btcIndex].Available)

	GetOrders()

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


	// Grab the fills every 10s
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

									// If a fill matches an order in existingBuys then create a sell for that buy
									// Todo: check here if the fill is completely filled
									if o.FillId == b.Id {
										fmt.Printf("Buy happened for %f at $%f\n", b.Size, b.Price)

										fmt.Printf("%v\n\n", existingBuys)

										CreateSellOrder(b.Price, b.Size)

										// Note: Can I just run GetOrders here?
										// Should update both existingBuys and existingSells
										existingBuys = RemoveOrder(existingBuys, i)

									}
								}
								for _, s := range existingSells {
									//fmt.Printf("s: %v\n", s)
									if o.FillId == s.Id {
										fmt.Printf("Sell happened for %f at $%f\n", s.Size, s.Price)
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
}


func contains(s []float64, e float64) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func pricesExisting(o Orders) []float64 {
	var pricesWithBuys []float64
	for _, a := range o {
		pricesWithBuys = append(pricesWithBuys, a.Price)
	}
	return pricesWithBuys
}

func RemoveOrder(o Orders, i int) Orders {
	fmt.Printf("%v\n\n", o)
	o[len(o)-1], o[i] = o[i], o[len(o)-1]
	return o[:len(o)-1]
}
