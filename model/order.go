package model

import (
	"github.com/shopspring/decimal"
)

// 持仓情况
type PositionDetail struct {
	Code string `json:"Zqdm"`
	Name string `json:"Zqmc"`
	// 可用数量
	AvailableQuantityStr string `json:"Kysl"`
	// 持仓数量
	TotalQuantityStr string `json:"Zqsl"`
	// 成本价（平均）
	CostPriceStr string `json:"Cbjg"`
}

// AccountDetail 账户详情
type AccountDetail struct {
	// 总资产
	TotalAssetStr string `json:"Zzc"`
	// 可用资金
	AvailableFundStr string           `json:"Kyzj"`
	Positions        []PositionDetail `json:"positions"`
}

// TradeType 交易类型
type TradeType string

// TradeTypeBuy 买入
const TradeTypeBuy TradeType = "B"

// TradeTypeSale 卖出
const TradeTypeSale TradeType = "S"

// 提交订单表单
type TradeOrderForm struct {
	Code      string          `json:"Zqdm"`
	Name      string          `json:"Zqmc"`
	Price     decimal.Decimal //价格
	Amount    int             // 数量
	TradeType TradeType       // 交易类型
}

// 成交订单
type DealOrder struct {
	OrderID     string `json:"Wtbh"` // 委托编号
	DealID      string `json:"Cjbh"` // 成交编号
	ClosingTime string `json:"Cjsj"` // 成交时间
	Code        string `json:"Zqdm"` // 证券代码
	Name        string `json:"Zqmc"` // 证券名称
	Type        string `json:"Mmlb"` // 委托方向
	AmountStr   string `json:"Cjsl"` // 成交数量
	PriceStr    string `json:"Cjjg"` // 成交价格
}

// 未成交的订单
type UnClosingOrder struct {
	Date          string `json:"Wtrq"` // 委托日期
	Time          string `json:"Wtsj"` // 委托时间
	OrderId       string `json:"Wtbh"` // 委托编号
	Code          string `json:"Zqdm"` // 证券代码
	Name          string `json:"Zqmc"` // 证券名称
	Type          string `json:"Mmsm"` // 委托方向
	Status        string `json:"Wtzt"` // 委托状态
	Price         string `json:"Wtjg"` // 委托价格
	Amount        string `json:"Wtsl"` // 委托数量
	ClosingPrice  string `json:"Cjjg"` // 成交价格
	ClosingAmount string `json:"Cjsl"` // 成交数量
}
