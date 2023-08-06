package model

import "github.com/shopspring/decimal"

// 股票行情数据
type Stockquote struct {
	Code string `json:"f57"`
	Name string `json:"f58"`

	SalePrice5      float64 `json:"f31"` // 卖5
	SalePrice5Count int     `json:"f32"`
	SalePrice4      float64 `json:"f33"` // 卖4
	SalePrice4Count int     `json:"f34"`
	SalePrice3      float64 `json:"f35"` // 卖3
	SalePrice3Count int     `json:"f36"`
	SalePrice2      float64 `json:"f37"` // 卖2
	SalePrice2Count int     `json:"f38"`
	SalePrice1      float64 `json:"f39"` // 卖1价
	SalePrice1Count int     `json:"f40"`

	BuyPrice1      float64 `json:"f19"` // 买1价
	BuyPrice1Count int     `json:"f20"`
	BuyPrice2      float64 `json:"f17"` // 买2价
	BuyPrice2Count int     `json:"f18"`
	BuyPrice3      float64 `json:"f15"` // 买3价
	BuyPrice3Count int     `json:"f16"`
	BuyPrice4      float64 `json:"f13"` // 买4价
	BuyPrice4Count int     `json:"f14"`
	BuyPrice5      float64 `json:"f11"` // 买5价
	BuyPrice5Count int     `json:"f12"`

	NewestPrice  float64 `json:"f43"` // 当前最新价
	OpenPrice    float64 `json:"f46"` // 开盘价
	HighestPrice float64 `json:"f44"` // 最高价
	LowestPrice  float64 `json:"f45"` // 最低价
	AvgPrice     float64 `json:"f71"` // 均价

	PreClosePrice float64 `json:"f60"` // 上一交易日的收盘价
}

func (s *Stockquote) GetActualPrice(magnification float64) {
	s.SalePrice5 = decimal.NewFromFloat(s.SalePrice5).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.SalePrice4 = decimal.NewFromFloat(s.SalePrice4).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.SalePrice3 = decimal.NewFromFloat(s.SalePrice3).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.SalePrice2 = decimal.NewFromFloat(s.SalePrice2).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.SalePrice1 = decimal.NewFromFloat(s.SalePrice1).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.BuyPrice1 = decimal.NewFromFloat(s.BuyPrice1).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.BuyPrice2 = decimal.NewFromFloat(s.BuyPrice2).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.BuyPrice3 = decimal.NewFromFloat(s.BuyPrice3).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.BuyPrice4 = decimal.NewFromFloat(s.BuyPrice4).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.BuyPrice5 = decimal.NewFromFloat(s.BuyPrice5).Div(decimal.NewFromFloat(magnification)).InexactFloat64()

	s.NewestPrice = decimal.NewFromFloat(s.NewestPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.OpenPrice = decimal.NewFromFloat(s.OpenPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.HighestPrice = decimal.NewFromFloat(s.HighestPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.LowestPrice = decimal.NewFromFloat(s.LowestPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.PreClosePrice = decimal.NewFromFloat(s.PreClosePrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.AvgPrice = decimal.NewFromFloat(s.AvgPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
}

type KlineType string

const (
	DailyKlineType   KlineType = "101"
	WeeklyKlineType  KlineType = "102"
	MonthlyKlineType KlineType = "103"
)

// k线数据
type QueryKlineParam struct {
	Code  string    // 股票代码
	Begin string    // 开始日期
	End   string    // 结束日期
	Type  KlineType // K线类型
}

// K线
type Kline struct {
	Code        string          // 股票代码
	Date        string          // 日期
	Amplitude   decimal.Decimal // 涨幅
	Volume      decimal.Decimal // 成交量
	OpenPrice   decimal.Decimal // 开盘价
	ClosePrice  decimal.Decimal // 收盘价格
	HigestPrice decimal.Decimal // 最高价
	LowestPrice decimal.Decimal // 最低价
}
