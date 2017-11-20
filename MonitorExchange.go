package main

import (
	"fmt"
	"time"

	ws "github.com/gorilla/websocket"
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

// Websocket Monitoring
//	1. If a sell happens create a buy at sell.Price-stepGap
// 	2. If a buy happens, create a sell at buy.Price+stepGap
//	3. If the price increases to a new high, keep buying

func MonitorExchange() {
	fmt.Printf("\n\n** MonitorExchange **\n\n")
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

			SetCurrentPrice(message.Price)

			if message.Side == "buy" {
				// Run through existing buys and see if this match aligns with any
				for _, o := range existingBuys {
					if message.MakerOrderId == o.Id {
						println("\n\n** -- ** -- Buy Happened! -- ** -- **\n\n")

						// Is this match a complete order?
						if message.RemainingSize != 0.0 {
							fmt.Printf("INCOMPLETE ORDER: only %f of %f filled\n\n", message.RemainingSize, o.Size)
						} else {
							// Create Sell at buy price plus stepGap
							CreateOrder("sell", o.Price + stepGap, o.Size)
							GetOrders()
						}
					}
				}
			} else if message.Side == "sell" {
				// Run through existing sells and see if this match aligns with any
				for _, o := range existingSells {
					if message.MakerOrderId == o.Id {
						println("\n\n** -- ** -- Sell Happened! -- ** -- **\n\n")

						// Is this match a complete order?
						if message.RemainingSize != 0.0 {
							fmt.Printf("INCOMPLETE ORDER: only %f of %f filled\n\n", message.RemainingSize, o.Size)
						} else {
							// Create buy order at sell price minus stepGap
							CreateOrder("buy", o.Price - stepGap, o.Size)
							GetOrders()
						}
					}
				}
			}
		}
	}
}

func SetCurrentPrice(price float64) {

	// Is this the very first match that we're seeing?
	if currentPrice == 0.0 {
		currentPrice = price
		CreateSteps()
		InitializeOrders()

	} else {
		t := time.Now()
		fmt.Printf("%s ||| next step: %f ||| current price: %f\n", t.Format(time.Kitchen), steps[stepsIndex], price)

		// Has the current price surpassed the next step?
		if steps[stepsIndex] < price {

			// 	Is there NOT a sell at current step + 1 AND is there NOT a buy at current step?
			if !Contains(PricesExisting(existingSells), steps[stepsIndex + 1]) && !Contains(PricesExisting(existingBuys), steps[stepsIndex]) {
				println("\n\n** -- ** -- Buy needs to be created! -- ** -- **\n\n")
				CreateOrder("buy", steps[stepsIndex], HowMuchToBuy(steps[stepsIndex]))
			}
			//stepsIndex++
			CreateSteps()

		// Has the current price dipped below the next lowest step?
		} else if steps[stepsIndex - 1] > price {
			//stepsIndex--
			CreateSteps()
		}

		currentPrice = price
	}
}

func CreateSteps()  {

	println("\n\n** CreateSteps **\n")

	currentNextStepPrice := float64(int(currentPrice / stepGap)) * stepGap + 1.0
	firstStep = currentNextStepPrice - (holdSteps * stepGap)
	lastStep = currentNextStepPrice + (holdSteps * stepGap)

	// Should also be (holdSteps * 2) + 1
	totalSteps = (lastStep - firstStep) / stepGap

	// Reset steps
	steps = steps[:0]
	
	// Create steps
	for i := firstStep; i <= lastStep; i += stepGap {
		steps = append(steps, i)
	}

	// Figure out how many steps are between the firstStep and currentPrice
	for a := range steps {
		if steps[a] <= currentPrice {
			stepsIndex = a + 1
		}
	}

	fmt.Printf("First Step: %f || Next step: %f || Last Step: %f || Total Steps: %f", firstStep, steps[stepsIndex], lastStep, totalSteps)
}
