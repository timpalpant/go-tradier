package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/timpalpant/go-tradier"
)

func showPositions(client *tradier.Client) {
	fmt.Println("Fetching current positions")
	positions, err := client.GetAccountPositions()
	if err != nil {
		log.Fatal(err)
	}

	if len(positions) == 0 {
		fmt.Println("No current positions")
		return
	}

	for _, p := range positions {
		fmt.Printf("%v (%v): %v shares, $ %.2f\n",
			p.Symbol, p.DateAcquired, p.Quantity, p.CostBasis)
	}
}

func gainLoss(client *tradier.Client) {
	fmt.Println("Fetching gain loss")
	cps, err := client.GetAccountCostBasis()
	if err != nil {
		log.Fatal(err)
	}

	if len(cps) == 0 {
		fmt.Println("No closed positions")
		return
	}

	for _, cp := range cps {
		fmt.Printf(
			"%v: open: %v, close: %v, quantity: %v, cost: $ %.2f, "+
				"proceeds: $ %.2f, gain-loss: $ %.2f, gain-loss percent: %.2f %%\n",
			cp.Symbol, cp.OpenDate.Time, cp.CloseDate.Time,
			cp.Quantity, cp.Cost, cp.Proceeds, cp.GainLoss, cp.GainLossPercent)
	}
}

func openOrders(client *tradier.Client) {
	fmt.Println("Fetching open orders")
	openOrders, err := client.GetOpenOrders()
	if err != nil {
		log.Fatal(err)
	}
	if len(openOrders) == 0 {
		fmt.Println("No open orders")
		return
	}

	var bought, sold, fees float64
	for _, o := range openOrders {
		fmt.Printf(
			"id: %v, type: %v, symbol: %v, side: %v, quantity: %v, status: %v\n",
			o.Id, o.Type, o.Symbol, o.Side, o.Quantity, o.Status)

		costBasis := o.ExecutedQuantity * o.AverageFillPrice
		switch o.Side {
		case tradier.Buy:
			bought += costBasis
		case tradier.Sell:
			sold += costBasis
		}

		if o.ExecutedQuantity > 0 {
			fees += 1.0
		}
	}

	dailyPL := sold - bought - fees
	fmt.Printf(
		"Daily PL = $ %.2f, (%.1f %% of $ %.2f invested)\n",
		dailyPL, 100*dailyPL/bought, bought)
}

func history(client *tradier.Client) {
	fmt.Println("Fetching account history")
	events, err := client.GetAccountHistory(1000)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range events {
		fmt.Printf("%v - %v - %.2f\n", e.Date, e.Type, e.Amount)
		switch e.Type {
		case "trade":
			fmt.Printf("\t%v - %v - %v shares, $ %.2f, commission = $ %.2f, trade type = %v\n",
				e.Trade.Symbol, e.Trade.Description, e.Trade.Quantity,
				e.Trade.Price, e.Trade.Commission, e.Trade.TradeType)
		case "adjustment":
		default:
			fmt.Printf("unknown event type: %v\n", e.Type)
		}
	}
}

func main() {
	subcommand := flag.String("command", "positions", "Command to run (positions, gainloss, openorders)")
	apiKey := flag.String("tradier.apikey", "", "Tradier API key")
	account := flag.String("tradier.account", "", "Tradier account ID")
	flag.Parse()

	params := tradier.DefaultParams(*apiKey)
	client := tradier.NewClient(params)
	client.SelectAccount(*account)

	switch *subcommand {
	case "positions":
		showPositions(client)
	case "gainloss":
		gainLoss(client)
	case "openorders":
		openOrders(client)
	case "history":
		history(client)
	default:
		log.Fatal("unknown command: ", *subcommand)
	}
}
