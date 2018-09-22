package tradier

import (
	"github.com/pkg/errors"
)

// StreamDemuxer demuxes the different types of messages in a market events stream.
type StreamDemuxer struct {
	Quotes    func(quote *QuoteEvent)
	Trades    func(trade *TradeEvent)
	Summaries func(summary *SummaryEvent)
	TimeSales func(timeSale *TimeSaleEvent)
	Errors    func(err error)
}

func (sd *StreamDemuxer) Handle(event *StreamEvent) {
	switch {
	case event.Type == "quote":
		sd.handleQuote(event)
	case event.Type == "trade":
		sd.handleTrade(event)
	case event.Type == "timesale":
		sd.handleTimeSale(event)
	case event.Type == "summary":
		sd.handleSummary(event)
	}
}

func (sd *StreamDemuxer) HandleChan(events <-chan *StreamEvent) {
	for event := range events {
		sd.Handle(event)
	}
}

func (sd *StreamDemuxer) handleQuote(m *StreamEvent) {
	if sd.Quotes != nil {
		if q, err := DecodeQuote(m); err == nil {
			sd.Quotes(q)
		} else {
			sd.Errors(errors.Wrapf(err, "error decoding quote: %v", string(m.Message)))
		}
	}
}

func (sd *StreamDemuxer) handleTrade(m *StreamEvent) {
	if sd.Trades != nil {
		if t, err := DecodeTrade(m); err == nil {
			sd.Trades(t)
		} else {
			sd.Errors(errors.Wrapf(err, "error decoding trade: %v", string(m.Message)))
		}
	}
}

func (sd *StreamDemuxer) handleSummary(m *StreamEvent) {
	if sd.Summaries != nil {
		if s, err := DecodeSummary(m); err == nil {
			sd.Summaries(s)
		} else {
			sd.Errors(errors.Wrapf(err, "error decoding summary: %v", string(m.Message)))
		}
	}
}

func (sd *StreamDemuxer) handleTimeSale(m *StreamEvent) {
	if sd.TimeSales != nil {
		if ts, err := DecodeTimeSale(m); err == nil {
			sd.TimeSales(ts)
		} else {
			sd.Errors(errors.Wrapf(err, "error decoding time sale: %v", string(m.Message)))
		}
	}
}
