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

### Fetch real-time top-of-book quotes

```Go
package main

import (
  "fmt"
  "net/http"

  "github.com/timpalpant/go-tradier"
)

func main() {
  client := iex.NewClient(&http.Client{})

  quotes, err := client.GetTOPS([]string{"AAPL", "SPY"})
  if err != nil {
      panic(err)
  }

  for _, quote := range quotes {
      fmt.Printf("%v: bid $%.02f (%v shares), ask $%.02f (%v shares) [as of %v]\n",
          quote.Symbol, quote.BidPrice, quote.BidSize,
          quote.AskPrice, quote.AskSize, quote.LastUpdated)
  }
}
```

### Fetch historical top-of-book quote (L1 tick) data.

Historical tick data (TOPS and DEEP) can be parsed using the `PcapScanner`.

```Go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/timpalpant/go-tradier"
	"github.com/timpalpant/go-tradier/iextp/tops"
)

func main() {
	client := iex.NewClient(&http.Client{})

	// Get historical data dumps available for 2016-12-12.
	histData, err := client.GetHIST(time.Date(2016, time.December, 12, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	} else if len(histData) == 0 {
		panic(fmt.Errorf("Found %v available data feeds", len(histData)))
	}

	// Fetch the pcap dump for that date and iterate through its messages.
	resp, err := http.Get(histData[0].Link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	packetDataSource, err := iex.NewPacketDataSource(resp.Body)
	if err != nil {
		panic(err)
	}
	pcapScanner := iex.NewPcapScanner(packetDataSource)

	// Write each quote update message to stdout, in JSON format.
	enc := json.NewEncoder(os.Stdout)

	for {
		msg, err := pcapScanner.NextMessage()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		switch msg := msg.(type) {
		case *tops.QuoteUpdateMessage:
			enc.Encode(msg)
		default:
		}
	}
}
```

## Contributing

Pull requests and issues are welcomed!

## License

go-tradier is released under the [GNU Lesser General Public License, Version 3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html)
