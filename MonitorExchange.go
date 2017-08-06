package main

import (
	"fmt"
  exchange "github.com/preichenberger/go-coinbase-exchange"
	ws "github.com/gorilla/websocket"
)

func MonitorExchange() {
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

  message := exchange.Message{}

  for true {
    if err := wsConn.ReadJSON(&message); err != nil {
      println(err.Error())
      break
    }

    if message.Type == "match"  {
			
			if currentPrice == 0.0 {
				GetOrders()
			}

			currentPrice = message.Price
      fmt.Printf("Current Price: $%f\n\n", currentPrice)

			if message.Side == "buy" {
        // run through existing buys and see if this match aligns with any
        for _, o := range existingBuys {
          if message.MakerOrderId == o.Id {
            // Buy Happened!

            // Check if this match is the complete order?
            // compare message.Size == o.Size ...
            // But is message.Size just the size of that match? (could be partial)

            // create Sell at buy price plus stopGap
            CreateSellOrder(o.Price + stopGap, o.Size)
          }
        }
      } else if message.Side == "sell" {
        // run through my existing sells and see if this match aligns with any
        for _, o := range existingSells {
          if message.MakerOrderId == o.Id {
            // Sell Happened!

            // Check if this match is the complete order?
            // compare message.Size == o.Size ...
            // But is message.Size just the size of that match? (could be partial)

            // create Buy at sell price minus stopGap
            CreateBuyOrder(o.Price - stopGap, o.Size)
          }
        }
      }
    }
  }
}
