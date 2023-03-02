# 东方财富网页版登录API
# 前言
仅供个人学习和个人自动交易使用，禁止用于其他用途。
其主要的重点在于东财的登录过程，可以使用任何的语言实现。剩余的其他API，仅满足本人的交易需求。
# 要求
1.python环境，安装[ddddocr](https://pypi.org/project/ddddocr/)库，用于验证码识别。  
2.在根目录下创建config.yaml文件  
```
user:  
  account: "资金账号"
  password: "交易密码"
```
## 新建东方财富客户端
```go
	client.NewEastMoneyClient()
```
## 提交委托订单
切记请勿在开盘时间测试！！！
```go
    // TradeTypeBuy 进行买入
    // TradeTypeSale 进行卖出
	c.SubmitTrade(model.TradeOrderForm{
			Code: "xxxxxxx",
			Name: "xxxxxxx",
			Amount:    100,
			Price:     decimal.NewFromFloat(2.856),
			TradeType: model.TradeTypeBuy,
		})
```

## 撤单
这个撤单支持批量操作，但是不建议这么操作。他的返回结果是没有状态码，只有一串字符串连在一起，形如（委托编号: 执行结果）。  
多条数据通过三个空格分割，批量操作不好判断，撤单是否执行成功。  
由于存在网络延迟的原因,委托订单生成的时间可能对应不上，因此撤单的Order需要通过 `GetRevokeList()` 获取（需要自己过滤指定的订单）。
```go
	c.RevokeOrders([]*model.Order{{
			Time:    "xxxxxx",
			OrderId: "xxxxxx",
		}
```

## 查询当日订单
东财的翻页逻辑做的实在太恶心了，需要根据某一页的数据来计算出上一页或下一页的页码，不能指定跳转某一页，所以他的查询接口，我就干脆直接将每页数据设置为100条了，这个交易频率足够我使用了。
```go
	// 当日成交的订单
	c.GetDealList()
	// 单日订单
	c.GetOrdersList()
	// 可撤销订单
	c.GetRevokeList()
```

## 查询K线数据
默认情况下只查询最近一个月的日K线数据
```go
	data, _ := api.GetKline(model.QueryKlineParam{
		Code: "xxxxx",
	})s
```

## 查询最新行情数据
其中：挂单的数据只关注买1和卖1的委托价格和委托量
```go
	api.GetQuote("xxxxx")
```

