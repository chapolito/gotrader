package main

import (
	"fmt"
	"time"
	"math"

	ws "github.com/gorilla/websocket"
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

// Websocket Monitoring
//	1. If a sell happens create a buy at sell.Price-stepGap
// 	2. If a buy happens, create a sell at buy.Price+stepGap
//	3. If the price increases to a new high, keep buying

func MonitorExchange() {
	fmt.Printf("\n** MonitorExchange **\n")

	// Websocket Connection
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

		if message.Type == "done" && message.Reason == "filled" {

			fmt.Printf("\n\nRemaining Size: %f || Old Size: %f || New Size: %f || Side: %s || Price: %f \n\n", message.RemainingSize, message.OldSize, message.NewSize, message.Side, message.Price)
			// fmt.Printf("\nRemaining Size: %f\n", message.RemainingSize)

			SetCurrentPrice(message.Price)

			if message.Side == "buy" {
				// Run through existing buys and see if this match aligns with any
				for _, o := range existingBuys {
					if message.MakerOrderId == o.Id {
						fmt.Printf("\n** Buy Happened! **\n")

						// Create Sell at buy price plus stepGap
						CreateOrder("sell", o.Price + stepGap, o.Size)
					}
				}
			} else if message.Side == "sell" {
				// Run through existing sells and see if this match aligns with any
				for _, o := range existingSells {
					if message.MakerOrderId == o.Id {
						fmt.Printf("\n** Sell Happened! **\n")

						// Create buy order at sell price minus stepGap
						CreateOrder("buy", o.Price - stepGap, o.Size)
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
		fmt.Printf("\n** First Time! **\n")
		CreateSteps()
		InitializeOrders()

	} else {
		currentPrice = price

		t := time.Now()
		fmt.Printf("%s ||| next step: %f ||| current price: %f\n", t.Format(time.Kitchen), steps[nextStepIndex], price)

		// Has the current price surpassed the next step?
		if steps[nextStepIndex] < price {
			fmt.Printf("\n** Surpassed! **\n")
			CreateSteps()
			CreateMissingBuys()
			PruneBuys()

		// Has the current price dipped below the next lowest step?
		} else if steps[nextStepIndex - 1] > price {
			fmt.Printf("\n** Dipped! **\n")
			CreateSteps()
			CreateMissingBuys()
		}
	}
}

func CreateSteps()  {

	fmt.Printf("\n** CreateSteps **\n")

	currentNextStepPrice := math.Floor(currentPrice/stepGap) * stepGap + stepGap
	firstStep = currentNextStepPrice - (holdSteps * stepGap)
	lastStep = currentNextStepPrice + (holdSteps * stepGap)

	// Should also be (holdSteps * 2)
	totalSteps = (lastStep - firstStep) / stepGap

	// Reset steps
	steps = steps[:0]

	// Create steps
	for i := firstStep; i <= lastStep; i += stepGap {
		steps = append(steps, i)
	}

	// Figure out how many steps are between the firstStep and currentPrice
	nextStepIndex = int(holdSteps)

	fmt.Printf("\nFirst Step: %f || currentPrice: %f || Next step: %f || Last Step: %f || Total Steps: %f\n", firstStep, currentPrice, steps[nextStepIndex], lastStep, totalSteps)
}
