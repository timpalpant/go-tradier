package tradier

import (
	"bufio"
	"encoding/json"
	"io"
)

// StreamEvent is used to unmarshal stream events before they are demuxed.
// Message contains the remainder of the type-specific message.
//
// StreamEvents can be demuxed into type-specific events using
// the StreamDemuxer.
type StreamEvent struct {
	Type    string
	Symbol  string
	Message json.RawMessage
	Error   error
}

func UnmarshalStreamEvent(buf []byte, se *StreamEvent) error {
	se.Message = make([]byte, len(buf))
	copy(se.Message, buf)
	se.Error = json.Unmarshal(buf, se)
	return se.Error
}

type QuoteEvent struct {
	Symbol      string
	Bid         float64
	BidSize     int64  `json:"bidsz"`
	BidExchange string `json:"bidexch"`
	BidDateMs   int64  `json:"biddate,string"`
	Ask         float64
	AskSize     int64  `json:"asksz"`
	AskExchange string `json:"askexch"`
	AskDateMs   int64  `json:"askdate,string"`
}

type TimeSaleEvent struct {
	Symbol   string
	Exchange string  `json:"exch"`
	Bid      float64 `json:",string"`
	Ask      float64 `json:",string"`
	Last     float64 `json:",string"`
	Size     int64   `json:",string"`
	DateMs   int64   `json:"date,string"`
}

type TradeEvent struct {
	Symbol           string
	Exchange         string  `json:"exch"`
	Price            float64 `json:",string"`
	Last             float64 `json:",string"`
	Size             int64   `json:",string"`
	CumulativeVolume int64   `json:"cvol,string"`
	DateMs           int64   `json:"date,string"`
}

type SummaryEvent struct {
	Symbol        string
	Open          float64 `json:",string"`
	High          float64 `json:",string"`
	Low           float64 `json:",string"`
	PreviousClose float64 `json:"prevClose,string'`
}

// MarketEventStream scans the newline-delimited market stream
// sent by Tradier and decodes each event into a StreamEvent.
type MarketEventStream struct {
	// A message on this channel indicates to the http consumer to shutdown the stream.
	// All channels will be closed by the goroutine that owns this stream.
	closeChan chan struct{}
}

func NewMarketEventStream(input io.ReadCloser, output chan *StreamEvent) *MarketEventStream {
	mes := &MarketEventStream{
		closeChan: make(chan struct{}),
	}
	go mes.consumeEvents(input, output)
	return mes
}

func (mes *MarketEventStream) Stop() {
	close(mes.closeChan)
}

func (mes *MarketEventStream) consumeEvents(
	input io.ReadCloser,
	output chan *StreamEvent) {
	scanner := bufio.NewScanner(input)
	defer input.Close()
	defer close(output)

	for scanner.Scan() {
		select {
		case <-mes.closeChan:
			return
		default:
			event := &StreamEvent{}
			if err := UnmarshalStreamEvent(scanner.Bytes(), event); err != nil {
				Logger.Println(err)
			}

			select {
			case output <- event:
			case <-mes.closeChan:
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		Logger.Println(err)
	}
}

func DecodeQuote(e *StreamEvent) (*QuoteEvent, error) {
	q := &QuoteEvent{Symbol: e.Symbol}
	err := json.Unmarshal(e.Message, q)
	return q, err
}

func DecodeTrade(e *StreamEvent) (*TradeEvent, error) {
	t := &TradeEvent{Symbol: e.Symbol}
	err := json.Unmarshal(e.Message, t)
	return t, err
}

func DecodeTimeSale(e *StreamEvent) (*TimeSaleEvent, error) {
	ts := &TimeSaleEvent{Symbol: e.Symbol}
	err := json.Unmarshal(e.Message, ts)
	return ts, err
}

func DecodeSummary(e *StreamEvent) (*SummaryEvent, error) {
	s := &SummaryEvent{Symbol: e.Symbol}
	err := json.Unmarshal(e.Message, s)
	return s, err
}
