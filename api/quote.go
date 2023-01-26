package api

import (
	"eastmoneyapi/model"
	"eastmoneyapi/util"
	"encoding/json"
	"net/http"
	"time"
)

// GetQuote 获取最新的行情数据
func GetQuote(code string) (*model.Stockquote, error) {
	req, _ := http.NewRequest("GET", "http://push2.eastmoney.com/api/qt/stock/get", nil)
	query := req.URL.Query()
	query.Add("fields", "f58,f734,f107,f57,f43,f59,f169,f170,f152,f177,f111,f46,f60,f44,f45,f47,f260,f48,f261,f279,f277,f278,f288,f19,f17,f531,f15,f13,f11,f20,f18,f16,f14,f12,f39,f37,f35,f33,f31,f40,f38,f36,f34,f32,f211,f212,f213,f214,f215,f210,f209,f208,f207,f206,f161,f49,f171,f50,f86,f84,f85,f168,f108,f116,f167,f164,f162,f163,f92,f71,f117,f292,f51,f52,f191,f192,f262,f294,f295,f269,f270,f256,f257,f285,f286")
	// 证券编号
	query.Add("secid", util.GetFullSecurityCode(code))
	req.URL.RawQuery = query.Encode()
	resp, err := (&http.Client{
		Timeout: time.Second * 3,
	}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 目前只关注买1价和卖1价
	var d struct {
		Data model.Stockquote `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}
	d.Data.GetActualPrice(util.GetPriceMagnification(d.Data.Code))
	return &d.Data, nil
}
