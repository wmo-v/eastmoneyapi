package service

import (
	"context"
	"eastmoneyapi/api"
	"eastmoneyapi/client"
	"eastmoneyapi/config"
	"eastmoneyapi/model"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var zInst *z_513050Svc

const (
	buyTransaction  = "买入事务"
	saleTransaction = "卖出事务"
)
const (
	orderNotSubmit = "未报"
	orderWithdrawn = "已撤"
	orderClosed    = "已成"
	orderSubmit    = "已报"
)
const (
	// 5 min
	orderMaxLifeTicket = 300
)

// 中概互联网ETF
type z_513050Svc struct {
	ctx       context.Context
	closeFn   context.CancelFunc
	closeOnce sync.Once
	c         *client.EastMoneyClient

	code              string
	transactionStatus string
	// 委托买单详情
	buyOrder *model.Order
	// 委托卖单详情
	saleOrder *model.Order
	// 局部行情
	localQuotes    *model.QuoteQueue
	currentQuotech chan *model.Stockquote

	// 当日最大止损次数为3次
	maxStopLossCount int
	// 停止交易的间隔;发生止损时,该ticket=600(10min);当ticket>0,禁止买入
	pauseTradingTicket int
	// 订单的存活时间,避免长时间的挂单
	orderLifeTicket int
}

func NewZ513050Svc() *z_513050Svc {
	var ctx, closeFn = context.WithCancel(context.Background())
	var z = &z_513050Svc{
		c:       client.NewEastMoneyClient(),
		ctx:     ctx,
		closeFn: closeFn,

		code:              "513050",
		transactionStatus: buyTransaction,
		// 每3秒获取一次值，这里就是15min的局部行情
		localQuotes:    model.NewQueue(300),
		currentQuotech: make(chan *model.Stockquote, 1),

		maxStopLossCount:   2,
		pauseTradingTicket: 0,
	}
	zInst = z
	return zInst
}

// 关闭定时任务
func (z *z_513050Svc) Close() {
	z.closeOnce.Do(func() {
		z.closeFn()
		close(z.currentQuotech)
	})
}

func (z *z_513050Svc) listenQuote() {
	for {
		select {
		case <-z.ctx.Done():
			logrus.Info("listenQuote 协程关闭")
			return
		default:
			quote, err := api.GetQuote(z.code)
			if err != nil || quote.Code == "" {
				logrus.Warning("获取行情数据失败：", err)
				continue
			}
			z.localQuotes.Enqueue(quote)
			z.currentQuotech <- quote
		}
		time.Sleep(3 * time.Second)
	}
}
func (z *z_513050Svc) ticketService() {
	for {
		select {
		case <-z.ctx.Done():
			return
		default:
			// 只有一个协程进行值的赋值,不需要考虑并发问题
			z.pauseTradingTicket--
			z.orderLifeTicket--
		}
		time.Sleep(time.Second)
	}
}

func (z *z_513050Svc) Start() {
	go z.listenQuote()
	go z.doAction()
	go z.ticketService()
}

func (z *z_513050Svc) doAction() error {
	for {
		switch z.transactionStatus {
		case buyTransaction:
			z.doBuy()
		case saleTransaction:
			z.doSale()
		}
	}
}

func (z *z_513050Svc) chekIsBuyPoint(curQuote *model.Stockquote) bool {
	if z.pauseTradingTicket > 0 {
		return false
	}
	// 买1价挂单量必须是卖1价挂单量的2倍
	if curQuote.BuyPrice1Count/curQuote.SalePrice1Count < 2 {
		return false
	}

	// 当前值是队列中的局部最小值，并且该值与队列中的最高值相差大于0.5%
	curPrice := decimal.NewFromFloat(curQuote.NewestPrice)
	// 局部最大值不能是当日的最大值
	if curPrice.Equal(decimal.NewFromFloat(curQuote.HighestPrice)) {
		return false
	}
	localMax := z.localQuotes.GetHighestPrice()
	return curPrice.Equal(z.localQuotes.GetLowestPrice()) &&
		curPrice.Div(localMax).LessThan(decimal.NewFromFloat(1-0.006))
}

func (z *z_513050Svc) doBuy() {
	for {
		// 获取最新的行情价格
		curQuote, ok := <-z.currentQuotech
		if !ok {
			return
		}
		if z.buyOrder == nil {
			if !z.chekIsBuyPoint(curQuote) {
				continue
			}
			orderId, err := z.c.SubmitTrade(model.TradeOrderForm{
				Code:      z.code,
				TradeType: model.TradeTypeBuy,
				Amount:    config.GetConfg().User.MaxAmount,
				Price:     decimal.NewFromFloat(curQuote.BuyPrice1),
			})
			if err != nil {
				logrus.Error("订单委托失败: ", err)
				continue
			}
			z.orderLifeTicket = orderMaxLifeTicket
			z.buyOrder = &model.Order{OrderId: orderId}
		} else {
			orderDetail, err := z.getOrder(z.buyOrder.OrderId)
			if err != nil {
				logrus.Warning("获取订单详情失败: ", err)
				continue
			}
			// 订单已撤
			if orderDetail.Status == orderWithdrawn {
				z.buyOrder = nil
			}
			// 订单完成
			if orderDetail.Status == orderClosed {
				z.buyOrder = orderDetail
				z.transactionStatus = saleTransaction
				return
			}
			// 订单已报，且长时间未成交，取消订单
			if orderDetail.Status == orderSubmit && z.orderLifeTicket < 0 {
				msg, _ := z.c.RevokeOrders([]*model.Order{z.buyOrder})
				logrus.Info("买单长时间未成交,取消订单：", msg)
			}
		}
	}
}

func (z *z_513050Svc) doSale() {
	for {
		// 获取最新的行情价格
		curQuote, ok := <-z.currentQuotech
		if !ok {
			return
		}
		if z.saleOrder == nil {
			amount, _ := strconv.Atoi(z.buyOrder.AmountStr)
			// 默认每单获利0.2%
			price := decimal.NewFromFloat(curQuote.SalePrice2)
			// 触发了止损条件
			if z.pauseTradingTicket > 0 {
				price = decimal.NewFromFloat(curQuote.SalePrice1)
			}
			orderId, err := z.c.SubmitTrade(model.TradeOrderForm{
				Code:      z.code,
				TradeType: model.TradeTypeSale,
				Amount:    amount,
				Price:     price,
			})
			if err != nil {
				logrus.Error("委托卖单失败: ", err)
				continue
			}
			z.saleOrder = &model.Order{OrderId: orderId}
		} else {
			orderDetail, err := z.getOrder(z.saleOrder.OrderId)
			if err != nil {
				logrus.Warning("获取订单详情失败: ", err)
				continue
			}

			// 订单已撤(有可能是手动撤单的结果)
			if orderDetail.Status == orderWithdrawn {
				z.saleOrder = nil
				continue
			}
			// 订单完成
			if orderDetail.Status == orderClosed {
				z.buyOrder = nil
				z.saleOrder = nil
				z.transactionStatus = buyTransaction

				// 订单在止损情况下完成的，最大止损次数-1
				if z.pauseTradingTicket > 0 {
					z.maxStopLossCount--
				}
				if z.maxStopLossCount == 0 {
					logrus.Warn("当日止损次数已达上限,程序终止")
					z.Close()
				}
				return
			}
			// 止损条件下,卖单应该是马上成交的，如果没有成交，撤单即可
			if orderDetail.Status == orderSubmit && z.iSNeedStopLoss(curQuote) {
				msg, _ := z.c.RevokeOrders([]*model.Order{z.saleOrder})
				logrus.Info("触发止损：", msg)
				z.pauseTradingTicket = 600
			}
		}
	}
}

// 止损
func (z *z_513050Svc) iSNeedStopLoss(q *model.Stockquote) bool {
	// costPrice, _ := decimal.NewFromString(z.buyOrder.ClosingPriceStr)
	/*
		止损0.5%
		x:现价, y:成本价
		(x-y)/y<-000.5
		x/y<0.995
	*/
	costPrice, _ := decimal.NewFromString(z.buyOrder.Price)
	curPrice := decimal.NewFromFloat(q.NewestPrice)
	return curPrice.Div(costPrice).LessThan(decimal.NewFromFloat(1 - 0.01))
}

func (z *z_513050Svc) getOrder(orderId string) (*model.Order, error) {
	list, err := z.c.GetOrdersList()
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].OrderId == orderId {
			return list[i], nil
		}
	}
	return nil, errors.New("没有找到对应的委托编号：" + orderId)
}
