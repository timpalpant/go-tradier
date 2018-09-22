package tradier

//go:generate ffjson $GOFILE
import (
	"encoding/json"
)

type Margin struct {
	FedCall           int     `json:"fed_call"`
	MaintenanceCall   int     `json:"maintenance_call"`
	OptionBuyingPower float64 `json:"option_buying_power"`
	StockBuyingPower  float64 `json:"stock_buying_power"`
	StockShortValue   float64 `json:"stock_short_value"`
	Sweep             int
}

type Cash struct {
	CashAvailable  float64 `json:"cash_available"`
	Sweep          int
	UnsettledFunds float64 `json:"unsettled_funds"`
}

type PDT struct {
	DayTradeBuyingPower float64 `json:"day_trade_buying_power"`
	FedCall             int     `json:"fed_call"`
	MaintenanceCall     int     `json:"maintenance_call"`
	OptionBuyingPower   float64 `json:"option_buying_power"`
	StockBuyingPower    float64 `json:"stock_buying_power"`
	StockShortValue     float64 `json:"stock_short_value"`
}

type AccountBalances struct {
	AccountNumber      string  `json:"account_number"`
	AccountType        string  `json:"account_type"`
	ClosePL            float64 `json:"close_pl"`
	CurrentRequirement float64 `json:"current_requirement"`
	Equity             float64
	LongMarketValue    float64 `json:"long_market_value"`
	MarketValue        float64 `json:"market_value"`
	OpenPL             float64 `json:"open_pl"`
	OptionLongValue    float64 `json:"option_long_value"`
	OptionRequirement  float64 `json:"option_requirement"`
	OptionShortValue   float64 `json:"option_short_value"`
	PendingOrdersCount int     `json:"pending_orders_count"`
	ShortMarketValue   float64 `json:"short_market_value"`
	StockLongValue     float64 `json:"stock_long_value"`
	TotalCash          float64 `json:"total_cash"`
	TotalEquity        float64 `json:"total_equity"`
	UnclearedFunds     float64 `json:"uncleared_funds"`
	Margin             Margin
	Cash               Cash
	PDT                PDT
}

type Position struct {
	CostBasis    float64  `json:"cost_basis"`
	DateAcquired DateTime `json:"date_acquired"`
	Id           int
	Quantity     float64
	Symbol       string
}

type Trade struct {
	Commission  float64
	Description string
	Price       float64
	Quantity    float64
	Symbol      string
	TradeType   string `json:"trade_type"`
}

type Adjustment struct {
	Description string
	Quantity    float64
}

type Event struct {
	Amount     float64
	Date       DateTime
	Type       string
	Trade      Trade
	Adjustment Adjustment
}

type ClosedPosition struct {
	CloseDate       DateTime `json:"close_date"`
	Cost            float64
	GainLoss        float64  `json:"gain_loss"`
	GainLossPercent float64  `json:"gain_loss_percent"`
	OpenDate        DateTime `json:"open_date"`
	Proceeds        float64
	Quantity        float64
	Symbol          string
	Term            int
}

const (
	// Order classes.
	Equity   = "equity"
	Option   = "option"
	Multileg = "multileg"
	Combo    = "combo"

	// Order sides.
	Buy         = "buy"
	BuyToCover  = "buy_to_cover"
	BuyToOpen   = "buy_to_open"
	BuyToClose  = "buy_to_close"
	Sell        = "sell"
	SellShort   = "sell_short"
	SellToOpen  = "sell_to_open"
	SellToClose = "sell_to_close"

	// Order types.
	MarketOrder    = "market"
	LimitOrder     = "limit"
	StopOrder      = "stop"
	StopLimitOrder = "stop_limit"
	Credit         = "credit"
	Debit          = "debit"
	Even           = "even"

	// Order durations.
	Day        = "day"
	GTC        = "gtc"
	PreMarket  = "pre"
	PostMarket = "post"

	// Option types.
	Put  = "put"
	Call = "call"

	// Order statuses.
	StatusOK        = "ok"
	Filled          = "filled"
	Canceled        = "canceled"
	Open            = "open"
	Expired         = "expired"
	Rejected        = "rejected"
	Pending         = "pending"
	PartiallyFilled = "partially_filled"
	Submitted       = "submitted"
)

type Order struct {
	Id                int
	Type              string
	Symbol            string
	OptionSymbol      string `json:"option_symbol"`
	Side              string
	Quantity          float64
	Status            string
	Duration          string
	Price             float64
	StopPrice         float64  `json:"stop_price"`
	OptionType        string   `json:"option_type"`
	ExpirationDate    DateTime `json:"expiration_date"`
	Exchange          string   `json:"exch"`
	AverageFillPrice  float64  `json:"avg_fill_price"`
	ExecutedQuantity  float64  `json:"exec_quantity"`
	ExecutionExchange string   `json:"exec_exch"`
	LastFillPrice     float64  `json:"last_fill_price"`
	LastFillQuantity  float64  `json:"last_fill_quantity"`
	RemainingQuantity float64  `json:"remaining_quantity"`
	CreateDate        DateTime `json:"create_date"`
	TransactionDate   DateTime `json:"transaction_date"`
	Class             string
	NumLegs           int `json:"num_legs"`
	Legs              []Order
	Strategy          string
}

// If there is only a single event, then tradier sends back
// an object, but if there are multiple events, then it sends
// a list of objects...
type OpenOrders []*Order

func (oo *OpenOrders) UnmarshalJSON(data []byte) error {
	orders := make([]*Order, 0)
	if err := json.Unmarshal(data, &orders); err == nil {
		*oo = orders
		return nil
	}

	order := Order{}
	err := json.Unmarshal(data, &order)
	if err == nil {
		*oo = []*Order{&order}
	}
	return err
}

// Helper struct for decoding Tradier's response to open orders request,
// which returns null if there are no open orders, a list if there are
// multiple open orders, or a single object if there is just one.
type openOrdersResponse struct {
	Orders struct {
		Order OpenOrders
	}
}

func (oor *openOrdersResponse) UnmarshalJSON(data []byte) error {
	// If there are no open orders, then Tradier returns "null".
	var noOrders struct {
		Orders string
	}

	err := json.Unmarshal(data, &noOrders)
	if err == nil {
		return nil
	}

	// Otherwise, unmarshal the results using the default unmarshaler.
	var result struct {
		Orders struct {
			Order OpenOrders
		}
	}
	err = json.Unmarshal(data, &result)
	oor.Orders = result.Orders
	return err
}

type OrderPreview struct {
	Commission    float64
	Cost          float64
	ExtendedHours bool `json:"extended_hours"`
	Fees          float64
	MarginChange  float64 `json:"margin_change"`
	Quantity      float64
	Status        string
}
