package tradier

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	SandboxEndpoint = "https://sandbox.tradier.com"
	APIEndpoint     = "https://api.tradier.com"
	StreamEndpoint  = "https://stream.tradier.com"
)

type MarketState string

const (
	MarketPremarket  MarketState = "premarket"
	MarketOpen       MarketState = "open"
	MarketPostmarket MarketState = "postmarket"
	MarketClosed     MarketState = "closed"
)

type Interval string

const (
	IntervalTick    Interval = "tick"
	IntervalMinute  Interval = "1min"
	Interval5Min    Interval = "5min"
	Interval15Min   Interval = "15min"
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
	IntervalMonthly Interval = "monthly"
)

type Filter string

const (
	FilterTrade    Filter = "trade"
	FilterQuote    Filter = "quote"
	FilterTimeSale Filter = "timesale"
	FilterSummary  Filter = "summary"
)

type SecurityType string

const (
	SecurityTypeStock      SecurityType = "stock"
	SecurityTypeIndex      SecurityType = "index"
	SecurityTypeETF        SecurityType = "etf"
	SecurityTypeMutualFund SecurityType = "mutual_fund"
)

var OldestDailyDate = time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)

type TradierError struct {
	Fault struct {
		FaultString string
		Detail      struct {
			ErrorCode string
		}
	}
	HttpStatusCode int
	Message        string
}

func (te TradierError) Error() string {
	return fmt.Sprintf("%d: %s - %s", te.HttpStatusCode, te.Fault.FaultString, te.Message)
}

type Security struct {
	Symbol      string
	Exchange    string
	Type        string
	Description string
}

type OptionSymbol struct {
	Root    string `json:"rootSymbol"`
	Options []string
}

type FloatOrNaN float64

func (f *FloatOrNaN) UnmarshalJSON(data []byte) error {
	var x float64
	var err error
	if err = json.Unmarshal(data, &x); err == nil {
		*f = FloatOrNaN(x)
		return nil
	}

	// Fallback for "NaN" values.
	var s string
	if strErr := json.Unmarshal(data, &s); strErr == nil {
		x, err = strconv.ParseFloat(s, 64)
		*f = FloatOrNaN(x)
	}

	return err
}

func (f FloatOrNaN) Value() (driver.Value, error) {
	if math.IsNaN(float64(f)) {
		return nil, nil
	}

	return float64(f), nil
}

type Greeks struct {
	Delta         float64
	Gamma         float64
	Theta         float64
	Vega          float64
	Rho           float64
	Phi           float64
	BidImpliedVol float64 `json:"bid_iv"`
	MidImpliedVol float64 `json:"mid_iv"`
	AskImpliedVol float64 `json:"ask_iv"`
	SMVImpliedVol float64 `json:"smv_vol"`
	UpdatedAt     float64
}

type Quote struct {
	Symbol           string
	Description      string
	Exchange         string `json:"exch"`
	Type             string
	Change           float64
	ChangePercentage float64 `json:"change_percentage"`
	Volume           int
	AverageVolume    int
	Last             float64
	LastVolume       int
	TradeDate        DateTime `json:"trade_date"`
	Open             float64
	High             float64
	Low              float64
	Close            float64
	PreviousClose    float64 `json:"prevclose"`
	Week52High       float64 `json:"week_52_high"`
	Week52Low        float64 `json:"week_52_low"`
	Bid              float64
	BidSize          int
	BidExchange      string   `json:"bidexch"`
	BidDate          DateTime `json:"bid_date"`
	Ask              float64
	AskSize          int
	AskExchange      string   `json:"askexch"`
	AskDate          DateTime `json:"ask_date"`
	OpenInterest     float64  `json:"open_interest"`
	Underlying       string
	Strike           float64
	ContractSize     int
	ExpirationDate   DateTime `json:"expiration_date"`
	ExpirationType   string   `json:"expiration_type"`
	OptionType       string   `json:"option_type"`
	RootSymbol       string   `json:"root_symbol"`
	Greeks           *Greeks
}

type TimeSale struct {
	Date      DateTime
	Time      DateTime
	Timestamp int64
	Open      FloatOrNaN
	Close     FloatOrNaN
	High      FloatOrNaN
	Low       FloatOrNaN
	Price     FloatOrNaN
	Vwap      FloatOrNaN
	Volume    int64
}

type MarketCalendar struct {
	Date        DateTime
	Status      string
	Description string
	Open        struct {
		Start string
		End   string
	}
	Premarket struct {
		Start string
		End   string
	}
	Postmarket struct {
		Start string
		End   string
	}
}

type MarketStatus struct {
	Time        DateTime `json:"date"`
	State       string
	Description string
	NextChange  DateTime `json:"next_change"`
	NextState   string   `json:"next_state"`
}
