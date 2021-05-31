package tradier

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
)

const (
	defaultRetries = 3

	// Header indicating the number of requests remaining.
	rateLimitAvailable = "X-Ratelimit-Available"
	// Header indicating the time at which our rate limit will renew.
	rateLimitExpiry = "X-Ratelimit-Expiry"

	// Error returned by Tradier if we make too big of a request.
	ErrBodyBufferOverflow = "protocol.http.TooBigBody"
)

var (
	// ErrNoAccountSelected is returned if account-specific methods
	// are attempted to be used without selecting an account first.
	ErrNoAccountSelected = errors.New("no account selected")
)

type ClientParams struct {
	Endpoint   string
	AuthToken  string
	Client     *http.Client
	Backoff    backoff.BackOff
	RetryLimit int
	Account    string
}

// DefaultParams returns ClientParams initialized with default values.
func DefaultParams(authToken string) ClientParams {
	return ClientParams{
		Endpoint:   APIEndpoint,
		AuthToken:  authToken,
		Client:     &http.Client{},
		Backoff:    backoff.NewExponentialBackOff(),
		RetryLimit: defaultRetries,
	}
}

// Client provides methods for making requests to the Tradier API.
type Client struct {
	client     *http.Client
	endpoint   string
	authHeader string
	backoff    backoff.BackOff
	retryLimit int

	account string
}

func NewClient(params ClientParams) *Client {
	return &Client{
		client:     params.Client,
		endpoint:   params.Endpoint,
		authHeader: fmt.Sprintf("Bearer %s", params.AuthToken),
		backoff:    params.Backoff,
		retryLimit: params.RetryLimit,
		account:    params.Account,
	}
}

func (tc *Client) SelectAccount(account string) {
	tc.account = account
}

func (tc *Client) GetAccountBalances() (*AccountBalances, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/balances"
	var result struct {
		Balances *AccountBalances
	}

	err := tc.getJSON(url, &result)
	return result.Balances, err
}

func (tc *Client) GetAccountPositions() ([]*Position, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/positions"
	var result struct {
		Positions struct {
			Position []*Position
		}
	}
	err := tc.getJSON(url, &result)
	return result.Positions.Position, err
}

func (tc *Client) GetAccountHistory(limit int) ([]*Event, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/history"
	if limit > 0 {
		url += fmt.Sprintf("?limit=%d", limit)
	}
	var result struct {
		History struct {
			Event []*Event
		}
	}
	err := tc.getJSON(url, &result)
	return result.History.Event, err
}

func (tc *Client) GetAccountCostBasis() ([]*ClosedPosition, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/gainloss"
	var result struct {
		GainLoss struct {
			ClosedPosition []*ClosedPosition `json:"closed_position"`
		} `json:"gainloss"`
	}
	err := tc.getJSON(url, &result)
	return result.GainLoss.ClosedPosition, err
}

func (tc *Client) GetOpenOrders() ([]*Order, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders"
	var result openOrdersResponse
	err := tc.getJSON(url, &result)
	return []*Order(result.Orders.Order), err
}

func (tc *Client) GetOrderStatus(orderId int) (*Order, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders/" + strconv.Itoa(orderId)
	var result struct {
		Order *Order
	}
	err := tc.getJSON(url, &result)
	return result.Order, err
}

func (tc *Client) PlaceOrder(order Order) (int, error) {
	if tc.account == "" {
		return 0, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders"
	form, err := orderToParams(order)
	if err != nil {
		return 0, err
	}

	resp, err := tc.do("POST", url, form, 0)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, errors.New(resp.Status + ": " + string(body))
	}

	var result struct {
		Order struct {
			Id     int
			Status string
		}
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return result.Order.Id, err
	} else if result.Order.Status != StatusOK {
		err = fmt.Errorf("received order status: %v", result.Order.Status)
	}
	return result.Order.Id, err
}

func (tc *Client) PreviewOrder(order Order) (*OrderPreview, error) {
	if tc.account == "" {
		return nil, ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders"
	form, err := orderToParams(order)
	if err != nil {
		return nil, err
	}

	form.Add("preview", "true")
	resp, err := tc.do("POST", url, form, tc.retryLimit)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(resp.Status + ": " + string(body))
	}

	var result struct {
		Order *OrderPreview
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return result.Order, err
	} else if result.Order == nil {
		err = fmt.Errorf("didn't receive order preview")
	} else if result.Order.Status != StatusOK {
		err = fmt.Errorf("received order status: %v", result.Order.Status)
	}
	return result.Order, err
}

// Convert the given order to URL parameters for a create order request.
// We also do some sanity checking to prevent placing orders with unset fields.
func orderToParams(order Order) (url.Values, error) {
	form := url.Values{}
	form.Add("class", order.Class)
	form.Add("duration", order.Duration)

	switch order.Class {
	case Equity, Option:
		form.Add("symbol", order.Symbol)
		form.Add("side", order.Side)
		form.Add("quantity", strconv.FormatFloat(order.Quantity, 'f', 0, 64))
		form.Add("type", order.Type)
		if order.Type == LimitOrder || order.Type == StopLimitOrder {
			form.Add("price", strconv.FormatFloat(order.Price, 'f', 2, 64))
		}
		if order.Type == StopOrder || order.Type == StopLimitOrder {
			form.Add("stop", strconv.FormatFloat(order.StopPrice, 'f', 2, 64))
		}
	case Multileg, Combo:
		form.Add("symbol", order.Symbol)
		form.Add("type", order.Type)
		if order.Type == LimitOrder || order.Type == StopLimitOrder {
			form.Add("price", strconv.FormatFloat(order.Price, 'f', 2, 64))
		}
		if order.Type == StopOrder || order.Type == StopLimitOrder {
			form.Add("stop", strconv.FormatFloat(order.StopPrice, 'f', 2, 64))
		}

		for i, leg := range order.Legs {
			form.Add(fmt.Sprintf("option_symbol[%d]", i), leg.OptionSymbol)
			form.Add(fmt.Sprintf("side[%d]", i), leg.Side)
			form.Add(fmt.Sprintf("quantity[%dd]", i), strconv.FormatFloat(leg.Quantity, 'f', 0, 64))
		}
	case OneTriggersOther, OneCancelsOther, OneTriggersOneCancelsOther:
		for i, leg := range order.Legs {
			form.Add(fmt.Sprintf("symbol[%d]", i), leg.Symbol)
			form.Add(fmt.Sprintf("quantity[%d]", i), strconv.FormatFloat(leg.Quantity, 'f', 0, 64))
			form.Add(fmt.Sprintf("type[%d]", i), leg.Type)
			form.Add(fmt.Sprintf("side[%d]", i), leg.Side)
			if leg.OptionSymbol != "" {
				form.Add(fmt.Sprintf("option_symbol[%d]", i), leg.OptionSymbol)
			}
			if leg.Type == LimitOrder || leg.Type == StopLimitOrder {
				form.Add(fmt.Sprintf("price[%d]", i), strconv.FormatFloat(leg.Price, 'f', 2, 64))
			}
			if leg.Type == StopOrder || leg.Type == StopLimitOrder {
				form.Add(fmt.Sprintf("stop[%d]", i), strconv.FormatFloat(leg.StopPrice, 'f', 2, 64))
			}
		}
	default:
		return form, fmt.Errorf("unknown order class: %v", order.Class)
	}
	return form, nil
}

func (tc *Client) ChangeOrder(orderId int, order Order) error {
	if tc.account == "" {
		return ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders/" + strconv.Itoa(orderId)
	form, err := updateOrderParams(order)
	if err != nil {
		return err
	}
	resp, err := tc.do("PUT", url, form, tc.retryLimit)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + ": " + string(body))
	}

	var result struct {
		Order struct {
			Id     int
			Status string
		}
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return err
	} else if result.Order.Status != StatusOK {
		return fmt.Errorf("received order status: %v", result.Order.Status)
	} else if result.Order.Id != orderId {
		return fmt.Errorf("changed order %v but received %v in response", orderId, result.Order.Id)
	}
	return nil
}

func updateOrderParams(order Order) (url.Values, error) {
	form := url.Values{}
	if order.Type != MarketOrder && order.Type != LimitOrder && order.Type != StopOrder && order.Type != StopLimitOrder {
		return form, fmt.Errorf("unknown order type: %v", order.Type)
	}
	form.Add("type", order.Type)
	if order.Duration != GTC && order.Duration != Day {
		return form, fmt.Errorf("unknown order duration: %v", order.Duration)
	}
	form.Add("duration", order.Duration)
	if order.Type == LimitOrder || order.Type == StopLimitOrder {
		if order.Price <= 0 {
			return form, fmt.Errorf("cannot place limit order without limit price")
		}
		form.Add("price", strconv.FormatFloat(order.Price, 'f', 2, 64))
	}
	if order.Type == StopOrder || order.Type == StopLimitOrder {
		if order.StopPrice <= 0 {
			return form, fmt.Errorf("cannot place stop order without stop price")
		}
		form.Add("stop", strconv.FormatFloat(order.StopPrice, 'f', 2, 64))
	}
	return form, nil
}

func (tc *Client) CancelOrder(orderId int) error {
	if tc.account == "" {
		return ErrNoAccountSelected
	}

	url := tc.endpoint + "/v1/accounts/" + tc.account + "/orders/" + strconv.Itoa(orderId)
	resp, err := tc.do("DELETE", url, nil, tc.retryLimit)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + ": " + string(body))
	}

	var result struct {
		Order struct {
			Id     int
			Status string
		}
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return err
	} else if result.Order.Status != StatusOK {
		return fmt.Errorf("received order status: %v", result.Order.Status)
	} else if result.Order.Id != orderId {
		return fmt.Errorf(
			"asked to cancel order %v but received %v in response",
			orderId, result.Order.Id)
	}
	return nil

}

// Get a list of symbols matching the given parameters.
func (tc *Client) LookupSecurities(
	types []SecurityType, exchanges []string, query string) (
	[]Security, error) {
	url := tc.endpoint + "/v1/markets/lookup"
	if len(types) > 0 {
		strTypes := make([]string, len(types))
		for i, t := range types {
			strTypes[i] = string(t)
		}
		url = url + "?types=" + strings.Join(strTypes, ",")
	}
	if exchanges != nil && len(exchanges) > 0 {
		url = url + "&exchanges=" + strings.Join(exchanges, ",")
	}
	if query != "" {
		url = url + "&q=" + query
	}

	var result struct {
		Securities struct {
			// TODO: If there is only one data point, then Tradier returns
			// a single object (not a list) and this fails to parse it.
			Security []Security
		}
	}
	err := tc.getJSON(url, &result)
	return result.Securities.Security, err
}

// Get a list of symbols matching the given parameters.
func (tc *Client) LookupOptionSymbols(underlying string) ([]OptionSymbol, error) {
	url := fmt.Sprintf("%s/v1/markets/options/lookup?underlying=%s", tc.endpoint, underlying)

	var result struct {
		Symbols []OptionSymbol
	}
	err := tc.getJSON(url, &result)
	return result.Symbols, err
}

// Get the securities on the Easy-to-Borrow list.
func (tc *Client) GetEasyToBorrow() ([]Security, error) {
	url := tc.endpoint + "/v1/markets/etb"
	var result struct {
		Securities struct {
			Security []Security
		}
	}
	err := tc.getJSON(url, &result)
	return result.Securities.Security, err
}

// Get an option's expiration dates.
func (tc *Client) GetOptionExpirationDates(symbol string) ([]time.Time, error) {
	params := "?symbol=" + symbol
	url := tc.endpoint + "/v1/markets/options/expirations" + params
	var result struct {
		Expirations struct {
			Date []DateTime
		}
	}
	err := tc.getJSON(url, &result)

	times := make([]time.Time, len(result.Expirations.Date))
	for i, dt := range result.Expirations.Date {
		times[i] = dt.Time
	}

	return times, err
}

// Get an option's expiration dates.
func (tc *Client) GetOptionStrikes(symbol string, expiration time.Time) ([]float64, error) {
	params := "?symbol=" + symbol + "&expiration=" + expiration.Format("2006-01-02")
	url := tc.endpoint + "/v1/markets/options/strikes" + params
	var result struct {
		Strikes struct {
			Strike []float64
		}
	}
	err := tc.getJSON(url, &result)
	return result.Strikes.Strike, err
}

// Get an option chain.
func (tc *Client) GetOptionChain(symbol string, expiration time.Time) ([]*Quote, error) {
	params := "?symbol=" + symbol + "&expiration=" + expiration.Format("2006-01-02")
	url := tc.endpoint + "/v1/markets/options/chains" + params
	var result struct {
		Options struct {
			Option []*Quote
		}
	}
	err := tc.getJSON(url, &result)
	return result.Options.Option, err
}

func (tc *Client) GetQuotes(symbols []string, greeks bool) ([]*Quote, error) {
	endpoint := fmt.Sprintf("%s/v1/markets/quotes", tc.endpoint)

	body := url.Values{}
	body.Set("greeks", strconv.FormatBool(greeks))
	body.Set("symbols", strings.Join(symbols, ","))

	var result struct {
		Quotes struct {
			Quote []*Quote
		}
	}
	err := tc.postJSON(endpoint, body, &result)
	return result.Quotes.Quote, err
}

func (tc *Client) getTimeSalesUrl(symbol string, interval Interval, start, end time.Time) string {
	url := tc.endpoint
	timeFormat := "2006-01-02T15:04:05"
	tz := time.UTC
	var err error
	if interval == IntervalDaily || interval == IntervalWeekly || interval == IntervalMonthly {
		url = url + "/v1/markets/history"
		timeFormat = "2006-01-02"
	} else {
		url = url + "/v1/markets/timesales"
		tz, err = time.LoadLocation("America/New_York")
		if err != nil {
			panic(err)
		}
	}
	url = url + "?symbol=" + symbol
	if interval != "" {
		url = url + "&interval=" + string(interval)
	}
	if !start.IsZero() {
		url = url + "&start=" + start.In(tz).Format(timeFormat)
	}
	if !end.IsZero() {
		url = url + "&end=" + end.In(tz).Format(timeFormat)
	}
	return url
}

// NOTE: If there is only one data point, then Tradier returns
// a single object (not a list). So first try list, and if parsing
// fails then fall back to try parsing a single object.
type timeSaleList []TimeSale

func (tsl *timeSaleList) UnmarshalJSON(data []byte) error {
	tss := make([]TimeSale, 0)
	if err := json.Unmarshal(data, &tss); err == nil {
		*tsl = tss
		return nil
	}

	ts := TimeSale{}
	err := json.Unmarshal(data, &ts)
	if err == nil {
		*tsl = []TimeSale{ts}
	}
	return err
}

func decodeTimeSales(reader io.Reader, interval Interval) ([]TimeSale, error) {
	dec := json.NewDecoder(reader)
	var timeSales []TimeSale
	if interval == IntervalDaily || interval == IntervalWeekly || interval == IntervalMonthly {
		var result struct {
			History struct {
				Day timeSaleList
			}
		}
		err := dec.Decode(&result)
		if err != nil {
			return nil, err
		}
		timeSales = result.History.Day
	} else {
		var result struct {
			Series struct {
				Data timeSaleList
			}
		}
		err := dec.Decode(&result)
		if err != nil {
			return nil, err
		}
		timeSales = result.Series.Data
	}

	return timeSales, nil
}

func bisect(start, end time.Time) time.Time {
	if end.IsZero() {
		end = time.Now()
	}

	delta := end.Sub(start)
	middle := start.Add(delta / 2)
	return middle
}

// Return daily, minute, or tick price bars for the given symbol.
// Tick data is available for the past 5 days, minute data for the past 20 days,
// and daily data since 1980-01-01.
// NOTE: The results are split, but not dividend-adjusted.
// https://developer.tradier.com/documentation/markets/get-history
// https://developer.tradier.com/documentation/markets/get-timesales
func (tc *Client) GetTimeSales(
	symbol string, interval Interval,
	start, end time.Time) ([]TimeSale, error) {

	url := tc.getTimeSalesUrl(symbol, interval, start, end)

	resp, err := tc.do("GET", url, nil, tc.retryLimit)
	if err != nil {
		if err, ok := err.(TradierError); ok {
			if err.Fault.Detail.ErrorCode == ErrBodyBufferOverflow {
				// Too much data for a single request!
				// Split the requested time interval in half and recurse.
				middle := bisect(start, end)
				if end.Sub(middle) < time.Duration(1*time.Minute) {
					// Give up if the interval is < 1min to prevent infinite recursion.
					return nil, err
				}

				firstHalf, err := tc.GetTimeSales(symbol, interval, start, middle)
				if err != nil {
					return nil, err
				}
				secondHalf, err := tc.GetTimeSales(symbol, interval, middle, end)
				if err != nil {
					return nil, err
				}
				allResults := make([]TimeSale, 0, len(firstHalf)+len(secondHalf))
				allResults = append(allResults, firstHalf...)
				allResults = append(allResults, secondHalf...)
				return allResults, nil
			}
		}

		// Some other error that we don't know how to handle.
		return nil, err
	}

	defer resp.Body.Close()
	return decodeTimeSales(resp.Body, interval)
}

// Subscribe to a stream of market events for the given symbols.
// Filter restricts the type of events streamed and can include:
// summary, trade, quote, timesale. If nil then all events are streamed.
// https://developer.tradier.com/documentation/streaming/get-markets-events
func (tc *Client) StreamMarketEvents(
	symbols []string, filter []Filter) (io.ReadCloser, error) {
	if len(symbols) == 0 {
		return nil, errors.New("list of symbols is required")
	}

	// First create a streaming session.
	createSessionUrl := tc.endpoint + "/v1/markets/events/session"

	createSessionResp, err := tc.do("POST", createSessionUrl, nil, tc.retryLimit)
	if err != nil {
		return nil, err
	}
	defer createSessionResp.Body.Close()
	if createSessionResp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(createSessionResp.Body)
		return nil, errors.New(createSessionResp.Status + ": " + string(body))
	}

	dec := json.NewDecoder(createSessionResp.Body)
	var sessionResp struct {
		Stream struct {
			SessionId string
			Url       string
		}
	}
	err = dec.Decode(&sessionResp)
	if err != nil {
		return nil, err
	}

	// Now open the stream.
	form := url.Values{}
	form.Add("linebreak", "true")
	form.Add("sessionid", sessionResp.Stream.SessionId)
	form.Add("symbols", strings.Join(symbols, ","))
	if len(filter) > 0 {
		strFilters := make([]string, len(filter))
		for i, f := range filter {
			strFilters[i] = string(f)
		}
		form.Add("filter", strings.Join(strFilters, ","))
	}
	// TODO: Make validOnly/flags configurable.
	form.Add("advancedDetails", "true")
	// If we fail here then just make a new session rather than retrying.
	// This prevents repeated failures to a session that doesn't exist for
	// some reason.
	resp, err := tc.do("POST", sessionResp.Stream.Url, form, 0)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return nil, errors.New("nil response with no error")
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(resp.Status + ": " + string(body))
	}

	return resp.Body, nil
}

// Get the market calendar for a given month.
func (tc *Client) GetMarketCalendar(year int, month time.Month) ([]MarketCalendar, error) {
	params := fmt.Sprintf("?year=%d&month=%d", year, month)
	url := tc.endpoint + "/v1/markets/calendar" + params
	var result struct {
		Calendar struct {
			Days struct {
				Day []MarketCalendar
			}
		}
	}

	err := tc.getJSON(url, &result)
	return result.Calendar.Days.Day, err
}

// Get the current state of the market (open/closed/etc.)
func (tc *Client) GetMarketState() (MarketStatus, error) {
	url := tc.endpoint + "/v1/markets/clock"
	var result struct {
		Clock MarketStatus
	}
	err := tc.getJSON(url, &result)
	return result.Clock, err
}

// Get corporate calendars.
func (tc *Client) GetCorporateCalendars(symbols []string) (
	GetCorporateCalendarsResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/calendars" + params
	var result GetCorporateCalendarsResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get company fundamentals.
func (tc *Client) GetCompanyInfo(symbols []string) (GetCompanyInfoResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/company" + params
	var result GetCompanyInfoResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get corporate actions.
func (tc *Client) GetCorporateActions(symbols []string) (GetCorporateActionsResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/corporate_actions" + params
	var result GetCorporateActionsResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get dividends.
func (tc *Client) GetDividends(symbols []string) (GetDividendsResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/dividends" + params
	var result GetDividendsResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get corporate ratios.
func (tc *Client) GetRatios(symbols []string) (GetRatiosResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/ratios" + params
	var result GetRatiosResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get financial reports.
func (tc *Client) GetFinancials(symbols []string) (GetFinancialsResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/financials" + params
	var result GetFinancialsResponse
	err := tc.getJSON(url, &result)
	return result, err
}

// Get price statistics.
func (tc *Client) GetPriceStatistics(symbols []string) (GetPriceStatisticsResponse, error) {
	params := "?symbols=" + strings.Join(symbols, ",")
	url := tc.endpoint + "/beta/markets/fundamentals/statistics" + params
	var result GetPriceStatisticsResponse
	err := tc.getJSON(url, &result)
	return result, err
}

func (tc *Client) getJSON(url string, result interface{}) error {
	resp, err := tc.do("GET", url, nil, tc.retryLimit)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + ": " + string(body))
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(result)
}

func (tc *Client) postJSON(url string, body url.Values, result interface{}) error {
	resp, err := tc.do("POST", url, body, tc.retryLimit)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(resp.Status + ": " + string(body))
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(result)
}

func (tc *Client) do(method, url string, body url.Values, maxRetries int) (*http.Response, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var sleep time.Duration
	for i := 0; i <= maxRetries; i++ {
		// Request must be made within retry loop, because body will be re-read each time.
		req, err = tc.makeSignedRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		resp, err = tc.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break // Successful request
		}

		if err != nil {
			Logger.Println(err)
			sleep = tc.backoff.NextBackOff()
		} else if resp.StatusCode != http.StatusOK {
			var respBody []byte
			respBody, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			tradierErr := TradierError{
				HttpStatusCode: resp.StatusCode,
			}
			if jsonErr := json.Unmarshal(respBody, &tradierErr); jsonErr == nil {
				// We extracted an error message, don't retry.
				return resp, tradierErr
			} else {
				tradierErr.Fault.FaultString = string(respBody)
			}
			// Assign an error since we have read the body. If this is the last retry,
			// we need to return a non-nil error.
			err = tradierErr
			rateLimitExpiry := parseQuotaViolationExpiration(tradierErr.Fault.FaultString)
			if rateLimitExpiry.After(time.Now().Add(sleep)) {
				sleep = rateLimitExpiry.Sub(time.Now()) + (1 * time.Second)
			} else {
				sleep = tc.backoff.NextBackOff()
			}
		}

		if i+1 <= maxRetries && sleep != backoff.Stop {
			Logger.Printf("Retrying after %v\n", sleep)
			time.Sleep(sleep)
		}
	}
	return resp, err
}

func (tc *Client) makeSignedRequest(method, url string, body url.Values) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = strings.NewReader(body.Encode())
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", tc.authHeader)
	if method != http.MethodDelete {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}
