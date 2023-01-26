package model

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// 持仓情况
type PositionDetail struct {
	Code                 string `json:"Zqdm"`
	Name                 string `json:"Zqmc"`
	AvailableQuantityStr string `json:"Kysl"`
	TotalQuantityStr     string `json:"Zqsl"`
}

func (p *PositionDetail) GetAvailableQuantity() int {
	res, _ := strconv.Atoi(p.AvailableQuantityStr)
	return res
}

func (p *PositionDetail) GetTotalQuantity() int {
	res, _ := strconv.Atoi(p.TotalQuantityStr)
	return res
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
	Amount      string `json:"Cjsl"` // 成交数量
	Price       string `json:"Cjjg"` // 成交价格
}

// 未成交的订单
type UnClosingOrder struct {
	Time          string `json:"Wtrq"` // 委托日期
	OrderId       string `json:"Wtbh"` // 委托编号
	Code          string `json:"Zqdm"` // 证券代码
	Name          string `json:"Zqmc"` // 证券名称
	Type          string `json:"Mmsm"` // 委托方向
	Amount        string `json:"Wtsl"` // 委托数量
	Price         string `json:"Wtzt"` // 委托价格
	ClosingPrice  string `json:"Cjjg"` // 成交价格
	ClosingAmount string `json:"Cjsl"` // 成交数量
}
