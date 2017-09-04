package main

import (
	"fmt"
	"time"

	ws "github.com/gorilla/websocket"
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

func MonitorExchange() {
	// Websocket
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": productId,
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	message := exchange.Message{}

	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		if message.Type == "match" {

			if currentPrice == 0.0 {
				currentPrice = message.Price
				// fmt.Printf("Current Price: $%f\n\n", currentPrice)
				GetOrders()
				InitializeOrders()
			} else {
				SetCurrentPrice(message.Price)
			}



			if message.Side == "buy" {
				// run through existing buys and see if this match aligns with any
				for _, o := range existingBuys {
					if message.MakerOrderId == o.Id {
						println("\n\n** -- ** -- Buy Happened! -- ** -- **\n\n")

						// Check if this match is the complete order?
						// compare message.Size == o.Size ...
						// But is message.Size just the size of that match? (could be partial)

						// create Sell at buy price plus stepGap
						CreateOrder("sell", o.Price + stepGap, o.Size)
						ResetOrders()
						GetOrders()
					}
				}
			} else if message.Side == "sell" {
				// run through my existing sells and see if this match aligns with any
				for _, o := range existingSells {
					if message.MakerOrderId == o.Id {
						println("\n\n** -- ** -- Sell Happened! -- ** -- **\n\n")

						// Check if this match is the complete order?
						// compare message.Size == o.Size ...
						// But is message.Size just the size of that match? (could be partial)

						// create Buy at sell price minus stepGap
						CreateOrder("buy", o.Price - stepGap, o.Size)
						ResetOrders()
						GetOrders()
					}
				}
			}
		}
	}
}

func SetCurrentPrice(price float64) {

	t := time.Now()
	fmt.Printf("%s ||| next step: %f ||| current price: %f\n", t.Format(time.Kitchen), steps[stepsIndex], price)

	// Has the current price surpassed the next step?
	if steps[stepsIndex] < price {

		// 	Is there NOT a sell at current step + 1 AND is there NOT a buy at current step?
		if !Contains(PricesExisting(existingSells), steps[stepsIndex + 1]) && !Contains(PricesExisting(existingBuys), steps[stepsIndex]) {
			println("\n\n** -- ** -- Buy needs to be created! -- ** -- **\n\n")
			CreateOrder("buy", steps[stepsIndex], HowMuchToBuy(steps[stepsIndex]))
		}
		stepsIndex++
	} else if steps[stepsIndex - 1] > price {
		stepsIndex--
	}

	currentPrice = price
}
