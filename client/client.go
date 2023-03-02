package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"eastmoneyapi/config"
	"eastmoneyapi/model"
	"eastmoneyapi/util"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	math_rand "math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	logrus "github.com/sirupsen/logrus"
)

var verifyCodeImgFile = "./verifyCode.jpg"
var baseUrl = "https://jywg.18.cn"

var client *EastMoneyClient
var intiClientOnce sync.Once

type EastMoneyClient struct {
	c           *http.Client
	closeCh     chan struct{}
	validateKey string
}

func NewEastMoneyClient() *EastMoneyClient {
	intiClientOnce.Do(func() {
		jar, _ := cookiejar.New(nil)
		client = &EastMoneyClient{
			c: &http.Client{
				Timeout: 3 * time.Second,
				Jar:     jar,
			},
			closeCh: make(chan struct{}, 1),
		}
		if err := client.login(config.GetConfg().User.Account, config.GetConfg().User.Password); err != nil {
			// 第一次登录失败，说明账号密码可能是错误的，直接panic
			panic("账号登录失败," + err.Error())
		}
		go func() {
			for {
				time.Sleep(time.Minute * 10)
				select {
				case <-client.closeCh:
					return
				default:
					client.login(config.GetConfg().User.Account, config.GetConfg().User.Password)
				}
			}
		}()
	})

	return client
}

// login 登录接口
func (e *EastMoneyClient) login(userId string, pwd string) error {
	var loginFn = func() error {
		randNumber := decimal.NewFromFloat(math_rand.Float64())
		if err := getVeriyCodeImg(randNumber.String()); err != nil {
			return errors.New("获取验证码失败: " + err.Error())
		}
		verifyCode, err := util.ImgOCR(verifyCodeImgFile)
		if err != nil {
			return errors.New("验证码识别失败: " + err.Error())
		}
		// 东方财富的验证码全是数字，如果识别出字母说明出错,不需要再往下执行了
		if _, err := strconv.Atoi(verifyCode); err != nil {
			return errors.New("验证码识别出错")
		}

		if err != nil {
			return errors.New("验证码安全加密识别失败: " + err.Error())
		}
		return e.doLogin(loginReq{
			userId:       userId,
			password:     pwd,
			verifyCode:   verifyCode,
			randNumber:   randNumber.String(),
			securityInfo: "",
		})
	}
	return util.Retry(5, loginFn)

}

type loginReq struct {
	userId       string
	password     string
	verifyCode   string
	randNumber   string
	securityInfo string
}

func (e *EastMoneyClient) doLogin(param loginReq) error {
	var formData = make(url.Values, 0)
	formData.Add("userId", param.userId)
	formData.Add("randNumber", param.randNumber)
	formData.Add("identifyCode", param.verifyCode)
	formData.Add("secInfo", param.securityInfo)
	formData.Add("password", encrypt(param.password))

	formData.Add("duration", "15")
	formData.Add("type", "Z")
	formData.Add("authCode", "")

	body := strings.NewReader(formData.Encode())
	req, _ := createRequestWithBaseHeader("POST", baseUrl+"/Login/Authentication", body)

	resp, err := e.c.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	var result = struct {
		Status  interface{} `json:"Status"`
		ErrCode interface{} `json:"Errcode"`
		Message string      `json:"Message"`
	}{}

	if err := bindJson(resp.Body, &result); err != nil {
		return err
	}
	if s, ok := result.Status.(float64); !ok || s != 0 {
		return errors.New(result.Message)
	}

	return e.getValidateKey()
}

// 这个ValidateKey隐藏在html中，随机访问一个页面，解析出来即可
func (e *EastMoneyClient) getValidateKey() error {
	req, _ := createRequestWithBaseHeader("GET", baseUrl+"/Search/Position", nil)
	resp, err := e.c.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	defer resp.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	target := doc.Find("#em_validatekey")
	if len(target.Nodes) != 1 {
		return errors.New("无法找到目标节点")
	}
	attrs := target.Nodes[0].Attr
	for i := range attrs {
		if attrs[i].Key == "value" {
			e.validateKey = attrs[i].Val
			return nil
		}
	}
	return errors.New("目标节点，没有value属性")
}

// SubmitTrade 提交订单交易
func (e *EastMoneyClient) SubmitTrade(order model.TradeOrderForm) (string, error) {
	var formData = make(url.Values, 0)
	formData.Add("stockCode", order.Code)
	formData.Add("zqmc", order.Name)
	formData.Add("amount", strconv.Itoa(order.Amount))
	formData.Add("tradeType", string(order.TradeType))
	formData.Add("market", util.GetMarket(order.Code))
	if util.IsEFT(order.Code) {
		order.Price = order.Price.Round(3)
	} else {
		order.Price = order.Price.Round(2)
	}
	formData.Add("price", order.Price.String())

	req, _ := createRequestWithBaseHeader(
		"POST",
		baseUrl+"/Trade/SubmitTradeV2?validatekey="+e.validateKey,
		strings.NewReader(formData.Encode()),
	)
	resp, err := e.c.Do(req)
	if err != nil {
		return "", errors.New(err.Error())
	}

	defer resp.Body.Close()
	var decoder = json.NewDecoder(resp.Body)
	var result = struct {
		Status  int    `json:"status"`
		Message string `json:"Message"`
		Data    []struct {
			OrderId string `json:"Wtbh"`
		} `json:"Data"`
	}{}
	if err := decoder.Decode(&result); err != nil {
		return "", errors.New(err.Error())
	}
	if result.Status != 0 {
		return "", errors.New(result.Message)
	}
	if len(result.Data) != 1 {
		return "", errors.New("未知情况发生，委托编号不是唯一")
	}
	msg := fmt.Sprintf(
		"\n订单委托成功:\n"+
			"\t委托编号: %s\n"+
			"\t委托时间: %s\n"+
			"\t代码: %s\n"+
			"\t名称: %s\n"+
			"\t委托数量: %d\n"+
			"\t委托价格: %s\n"+
			"\t委托方向: %s\n",
		result.Data[0].OrderId,
		time.Now().Format("2006-01-02 15:04:05"),
		order.Code,
		order.Name,
		order.Amount,
		order.Price.String(),
		order.TradeType)
	log.Println(msg)
	return result.Data[0].OrderId, nil
}

// GetOrdersList 获取当日的所有订单信息
func (e *EastMoneyClient) GetOrdersList() ([]*model.Order, error) {
	return e.getOrders(baseUrl + "/Search/GetOrdersData?validatekey=" + e.validateKey)
}

// GetDealList 获取当日成交信息
func (e *EastMoneyClient) GetDealList() ([]*model.Order, error) {
	return e.getOrders(baseUrl + "/Search/GetDealData?validatekey=" + e.validateKey)
}

// GetRevokeList 获取可撤单的订单信息
func (e *EastMoneyClient) GetRevokeList() ([]*model.Order, error) {
	return e.getOrders(baseUrl + "/Trade/GetRevokeList?validatekey=" + e.validateKey)
}

func (e *EastMoneyClient) getOrders(u string) ([]*model.Order, error) {
	var form = make(url.Values, 0)
	form.Add("qqhs", "100")
	req, _ := createRequestWithBaseHeader("POST", u, strings.NewReader(form.Encode()))
	resp, err := e.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result = struct {
		Data    []*model.Order `json:"Data"`
		Status  int            `json:"Status"`
		Message string         `json:"Message"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// RevokeOrders 撤单，支持批量撤单，但是不建议使用，返回一串的字符串，需要自行判断有没有撤单成功。
// 格式为： 委托编号: 消息
func (e *EastMoneyClient) RevokeOrders(list []*model.Order) (string, error) {
	if len(list) == 0 {
		return "没有需要撤单的交易", nil
	}

	var revokes = ""
	for i := range list {
		revokes += list[i].Date + "_" + list[i].OrderId + ","
	}
	revokes = revokes[:len(revokes)-1]
	var form = make(url.Values)
	form.Add("revokes", revokes)

	var req, _ = createRequestWithBaseHeader(
		"Post",
		baseUrl+"/Trade/RevokeOrders?validatekey="+e.validateKey,
		strings.NewReader(form.Encode()),
	)
	resp, err := e.c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

// GetStockList 查询当前的持仓情况
func (e *EastMoneyClient) GetStockList() ([]*model.PositionDetail, error) {
	var formData = make(url.Values, 0)
	formData.Add("qqhs", "10")
	req, _ := createRequestWithBaseHeader("POST", baseUrl+"/Search/GetStockList", strings.NewReader(formData.Encode()))
	resp, err := e.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var decoder = json.NewDecoder(resp.Body)
	var result = struct {
		Message string                  `json:"Message"`
		Data    []*model.PositionDetail `json:"Data"`
	}{}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, err
}

// QueryAssetAndPosition 查询账户资产和持仓情况
func (e *EastMoneyClient) QueryAssetAndPosition() (*model.AccountDetail, error) {
	var form = make(url.Values, 0)
	form.Add("moneyType", "RMB")
	req, _ := createRequestWithBaseHeader(
		"Post",
		baseUrl+"/Com/queryAssetAndPositionV1?validatekey="+e.validateKey,
		strings.NewReader(form.Encode()))
	resp, err := e.c.Do(req)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data []model.AccountDetail `json:"Data"`
	}
	if err := bindJson(resp.Body, &result); err != nil {
		return nil, err
	}
	if len(result.Data) != 1 {
		return nil, errors.New("仅支持查询一个账户详情")
	}
	return &result.Data[0], nil
}

func createRequestWithBaseHeader(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("sec-ch-ua-platform", "Windows")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

func (e *EastMoneyClient) getSecurityInfo(code string) (string, error) {
	resp, err := http.Get("http://127.0.0.1:18888/api/verifyUserInfo?" + code)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var decoder = json.NewDecoder(resp.Body)
	var data = struct {
		UserInfo string `json:"userInfo"`
	}{}
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}
	return data.UserInfo, nil
}

// 获取验证码图片, 需要传入一个数字绑定图片
func getVeriyCodeImg(randNum string) error {
	resp, err := http.Get(baseUrl + "/Login/YZM?randNum=" + randNum)
	if err != nil {
		return errors.New(err.Error())
	}
	defer resp.Body.Close()
	f, err := os.OpenFile(verifyCodeImgFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return errors.New(err.Error())
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
func bindJson(r io.ReadCloser, t interface{}) error {
	defer r.Close()
	var decoder = json.NewDecoder(r)
	return decoder.Decode(t)
}

const pubPEM = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHdsyxT66pDG4p73yope7jxA92
c0AT4qIJ/xtbBcHkFPK77upnsfDTJiVEuQDH+MiMeb+XhCLNKZGp0yaUU6GlxZdp
+nLW8b7Kmijr3iepaDhcbVTsYBWchaWUXauj9Lrhz58/6AE/NF0aMolxIGpsi+ST
2hSHPu3GSXMdhPCkWQIDAQAB
-----END PUBLIC KEY-----`

func encrypt(value string) string {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(value))
	if err != nil {
		panic("encrypt failed: " + err.Error())
	}
	enc_str := base64.StdEncoding.EncodeToString([]byte(ciphertext))
	return enc_str
}
