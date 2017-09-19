package main

//
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
// Get Fills
// *~*~*~*~*~*~*~*~*~**~*~*~*~*~*~*~*~*
//
// Demonstrates profit and rate of profit
//
//  1. Display total assets:
//
//    - Net assets in USD:
//      usdIndex +
//      ethIndex * eth4usdExchangeRate + (in theory, this amount would be zero, the program should never hold ETH)
//      existingBuys +
//      existingSells
//
//  2. Calculate Profit for last week
//    - Only count sells that have a price in totalSteps
//      This will allow to run multiple tests on the same account
//    - profit for each sell: (Price * Size) - ((Price - stepGap) * Size)
//
//

// import (
// 	"time"
// 	"fmt"
//   exchange "github.com/preichenberger/go-coinbase-exchange"
// )

// Grab the fills every 10s
// ticker := time.NewTicker(10 * time.Second)
// quit := make(chan struct{})
// go func() {
//     for {
//        select {
//         case <- ticker.C:
// 					// Get fills
// 					var fills []exchange.Fill
// 					cursorFills := client.ListFills()
//
// 					for cursorFills.HasMore {
// 						if err := cursorFills.NextPage(&fills); err != nil {
// 							println(err.Error())
// 							return
// 						}
//
// 						for _, o := range fills {
// 							//fmt.Printf("o: %v\n", o)
// 							for i, b := range existingBuys {
// 								//fmt.Printf("buy:  %v\nfill: %v\n\n", b[2], o.FillId)
//
// 								// If a fill matches an order in existingBuys then create a sell for that buy
// 								// Todo: check here if the fill is completely filled
// 								if o.FillId == b.Id {
// 									fmt.Printf("Buy happened for %f at $%f\n", b.Size, b.Price)
//
// 									fmt.Printf("%v\n\n", existingBuys)
//
// 									CreateSellOrder(b.Price, b.Size)
//
// 									// Note: Can I just run GetOrders here?
// 									// Should update both existingBuys and existingSells
// 									existingBuys = RemoveOrder(existingBuys, i)
//
// 								}
// 							}
// 							for _, s := range existingSells {
// 								//fmt.Printf("s: %v\n", s)
// 								if o.FillId == s.Id {
// 									fmt.Printf("Sell happened for %f at $%f\n", s.Size, s.Price)
// 								}
// 							}
// 						}
// 					}
//         case <- quit:
//             ticker.Step()
//             return
//         }
//     }
//  }()

// func RemoveOrder(o Orders, i int) Orders {
// 	fmt.Printf("%v\n\n", o)
// 	o[len(o)-1], o[i] = o[i], o[len(o)-1]
// 	return o[:len(o)-1]
// }
