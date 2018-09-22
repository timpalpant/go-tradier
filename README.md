# go-tradier
A Go library for accessing the Tradier Developer API.

[![GoDoc](https://godoc.org/github.com/timpalpant/go-tradier?status.svg)](http://godoc.org/github.com/timpalpant/go-tradier)
[![Build Status](https://travis-ci.org/timpalpant/go-tradier.svg?branch=master)](https://travis-ci.org/timpalpant/go-tradier)
[![Coverage Status](https://coveralls.io/repos/timpalpant/go-tradier/badge.svg?branch=master&service=github)](https://coveralls.io/github/timpalpant/go-tradier?branch=master)

go-tradier is a library to access the [Tradier Developer API](https://developer.tradier.com/documentation) from [Go](http://www.golang.org).
It provides a thin wrapper for working with the JSON REST endpoints.

[Tradier](https://tradier.com) is the first Brokerage API company, powering the world's leading
trading, investing and digital advisor platforms. Tradier is not affiliated
and does not endorse or recommend this library.

## Usage

### tcli

The `tcli` tool is a small command-line interface for making requests.

```shell
$ go install github.com/timpalpant/go-tradier/tcli
$ tcli -tradier.account XXXXX -tradier.apikey XXXXX -command positions
```

### Fetch real-time top-of-book quotes

```Go
package main

import (
  "fmt"

  "github.com/timpalpant/go-tradier"
)

func main() {
  params := tradier.DefaultParams("your-api-key-here")
  client := tradier.NewClient(params)

  quotes, err := client.GetQuotes([]string{"AAPL", "SPY"})
  if err != nil {
    panic(err)
  }

  for _, quote := range quotes {
    fmt.Printf("%v: bid $%.02f (%v shares), ask $%.02f (%v shares)\n",
      quote.Symbol, quote.Bid, quote.BidSize, quote.Ask, quote.AskSize)
  }
}
```

### Stream real-time top-of-book trades and quotes (L1 TAQ) data.

```Go
package main

import (
	"fmt"

	"github.com/timpalpant/go-tradier"
)

func main() {
	params := tradier.DefaultParams("your-api-key-here")
	client := tradier.NewClient(params)

	eventsReader, err := client.StreamMarketEvents(
		[]string{"AAPL", "SPY"},
		[]tradier.Filter{tradier.FilterQuote, tradier.FilterTrade})
	if err != nil {
		panic(err)
	}

	eventsCh := make(chan *tradier.StreamEvent)
	eventStream := tradier.NewMarketEventStream(eventsReader, eventsCh)
	defer eventStream.Stop()

	demuxer := tradier.StreamDemuxer{
		Quotes: func(quote *tradier.QuoteEvent) {
			fmt.Printf("QUOTE %v: bid $%.02f (%v shares), ask $%.02f (%v shares)\n",
				quote.Symbol, quote.Bid, quote.BidSize, quote.Ask, quote.AskSize)
		},
		Trades: func(trade *tradier.TradeEvent) {
			fmt.Printf("TRADE %v: $%.02f (%v shares) at %v\n",
				trade.Symbol, trade.Price, trade.Size, trade.DateMs)
		},
	}

	demuxer.HandleChan(eventsCh)
}
```

### Place and then cancel an order for SPY.

```Go
package main

import (
	"fmt"
	"time"

	"github.com/timpalpant/go-tradier"
)

func main() {
	params := tradier.DefaultParams("your-api-key-here")
	client := tradier.NewClient(params)
	client.SelectAccount("your-account-id-here")

	// Place a limit order for 1 share of SPY at $1.00.
	orderId, err := client.PlaceOrder(tradier.Order{
		Class:    tradier.Equity,
		Type:     tradier.LimitOrder,
		Symbol:   "SPY",
		Side:     tradier.Buy,
		Quantity: 1,
		Price:    1.00,
		Duration: tradier.Day,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Placed order: %v\n", orderId)

	time.Sleep(2 * time.Second)
	order, err := client.GetOrderStatus(orderId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Order status: %v\n", order.Status)

	// Cancel the order.
	fmt.Printf("Canceling order: %v\n", orderId)
	if err := client.CancelOrder(orderId); err != nil {
		panic(err)
	}
}
```

## Contributing

Pull requests and issues are welcomed!

## License

go-tradier is released under the [GNU Lesser General Public License, Version 3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html)
