package api

import (
	"bytes"
	"eastmoneyapi/model"
	"eastmoneyapi/util"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const timeFormat = "20060102"

const month = 30 * 24 * time.Hour

func createDefaultKlineQuery(q model.QueryKlineParam) model.QueryKlineParam {
	date := time.Now().Add(-month).Format(timeFormat)
	var result = model.QueryKlineParam{
		Code:  q.Code,
		Begin: date,
		End:   "20500101",
		Type:  "101",
	}
	if q.Begin != "" {
		result.Begin = q.Begin
	}
	if q.End != "" {
		result.End = q.End
	}
	if q.Type != "" {
		result.Type = q.Type
	}
	return result
}

// GetKline 获取K线数据
// 当K线数据为1分钟K线时，Begin 和 End 不起作用，仅仅能获取最近一个交易日的数据，无法获取历史数据
func GetKline(q model.QueryKlineParam) ([]*model.Kline, error) {
	param := createDefaultKlineQuery(q)
	var client = http.Client{}
	req, _ := http.NewRequest("GET", "http://push2his.eastmoney.com/api/qt/stock/kline/get", nil)
	query := req.URL.Query()
	query.Add("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	query.Add("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	query.Add("beg", param.Begin)
	query.Add("end", param.End)
	query.Add("secid", util.GetFullSecurityCode(param.Code))
	// 日k线
	query.Add("klt", string(param.Type))
	// 前复权
	query.Add("fqt", "1")
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	return parseKline(buf.Bytes())
}

func parseKline(buf []byte) ([]*model.Kline, error) {
	var tmp struct {
		Data struct {
			Code     string   `json:"code"`
			PrePrice float64  `json:"prePrice"`
			Klines   []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.Unmarshal(buf, &tmp); err != nil {
		return nil, err
	}
	// 2022-04-06,2.904,2.911,2.921,2.889,8255950,2398340848.000,1.09,-0.41,-0.012,3.62
	// 2022-10-14,29.73,30.28,30.60,29.35,544783,1639917704.17,4.24,2.64,0.78,0.59
	var result = make([]*model.Kline, 0, len(tmp.Data.Klines))
	for _, str := range tmp.Data.Klines {
		var fields = strings.Split(str, ",")
		openPrice, _ := decimal.NewFromString(fields[1])
		closePrice, _ := decimal.NewFromString(fields[2])
		maxPrice, _ := decimal.NewFromString(fields[3])
		minPrice, _ := decimal.NewFromString(fields[4])
		volumn, _ := decimal.NewFromString(fields[5])
		amp, _ := decimal.NewFromString(fields[8])
		result = append(result, &model.Kline{
			Code:        tmp.Data.Code,
			Date:        fields[0],
			Amplitude:   amp,
			Volume:      volumn,
			OpenPrice:   openPrice,
			ClosePrice:  closePrice,
			HigestPrice: maxPrice,
			LowestPrice: minPrice,
		})
	}
	return result, nil
}
