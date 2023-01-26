package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const SZMarket = "SA"
const SHMarket = "HA"

func GetMarket(code string) string {
	if strings.HasPrefix(code, "60") ||
		strings.HasPrefix(code, "51") ||
		strings.HasPrefix(code, "68") {
		return SHMarket
	}
	if strings.HasPrefix(code, "00") ||
		strings.HasPrefix(code, "15") ||
		strings.HasPrefix(code, "30") {
		return SZMarket
	}
	panic("暂未支持的证券代码")
}

func Retry(maxTimes int, fn func() error) error {
	var err error
	for count := 0; count < maxTimes; count++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 20 * time.Duration(count))
	}
	return fmt.Errorf("重试 %d 后，执行依旧失败，最后一次失败的原因：%s", maxTimes, err.Error())
}

// 上证基金代码以50、51、52开头，
// 深证基金代码以15、16、18开头。
func GetFullSecurityCode(code string) string {
	if strings.HasPrefix(code, "51") ||
		strings.HasPrefix(code, "60") ||
		strings.HasPrefix(code, "68") {
		return "1." + code
	}
	if strings.HasPrefix(code, "1") ||
		strings.HasPrefix(code, "0") ||
		strings.HasPrefix(code, "3") {
		return "0." + code
	}
	panic("暂未支持的证券代码")
}

// 对于ETF而言，价格的倍率是1000，对于个股而言，倍率是100
func GetPriceMagnification(code string) float64 {
	if strings.HasPrefix(code, "5") || strings.HasPrefix(code, "1") {
		return 1000.0
	}
	return 100.0
}

// 获取股票所在的板块
func GetCodeMarket(code string) string {
	if strings.HasPrefix(code, "5") || strings.HasPrefix(code, "6") {
		return "1"
	}
	if strings.HasPrefix(code, "1") || strings.HasPrefix(code, "0") || strings.HasPrefix(code, "3") {
		return "0"
	}
	panic("unsupport security code")
}

// 获取最接近的整百数
func GetNearestHundredfoldInt(num float64) int64 {
	return (decimal.NewFromFloat(num).IntPart() / 100) * 100
}
