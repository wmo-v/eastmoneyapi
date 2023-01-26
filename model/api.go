package model

import "github.com/shopspring/decimal"

// 股票行情数据
type Stockquote struct {
	Code            string  `json:"f57"`
	Name            string  `json:"f58"`
	BuyPrice1       float64 `json:"f19"` // 买1价
	BuyPrice1Count  int     `json:"f20"` // 买1挂单数
	SalePrice1      float64 `json:"f39"` // 卖1价
	SalePrice1Count int     `json:"f40"` // 卖1挂单数

	NewestPrice float64 `json:"f43"` // 当前最新价
	OpenPrice   float64 `json:"f46"` // 开盘价
	HigestPrice float64 `json:"f44"` // 最高价
	LowestPrice float64 `json:"f45"` // 最低价
	AvgPrice    float64 `json:"f71"` // 均价

	PreClosePrice float64 `json:"f60"` // 上一交易日的收盘价
}

func (s *Stockquote) GetActualPrice(magnification float64) {
	s.BuyPrice1 = decimal.NewFromFloat(s.BuyPrice1).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.SalePrice1 = decimal.NewFromFloat(s.SalePrice1).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.NewestPrice = decimal.NewFromFloat(s.NewestPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.OpenPrice = decimal.NewFromFloat(s.OpenPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
	s.HigestPrice = decimal.NewFromFloat(s.HigestPrice).Div(decimal.NewFromFloat(magnification)).InexactFloat64()
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
