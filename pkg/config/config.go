package config

type Category int64

const (
	Status Category = iota
	ExchangeInfo
	DailyTicker
	CandleStick
)

func (c Category) String() string {
	switch c {
	case Status:
		return "/api/v3/ping"
	case ExchangeInfo:
		return "/api/v3/exchangeInfo?permissions=SPOT"
	case DailyTicker:
		return "/api/v3/ticker/24hr"
	case CandleStick:
		return "/api/v3/klines"
	}
	return "unknown"
}
