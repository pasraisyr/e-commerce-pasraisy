package toyyibpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	SecretKey    string
	CategoryCode string
	IsSandbox    bool
	HTTPClient   *http.Client
}

func NewClient(secretKey, categoryCode string, isSandbox bool) *Client {
	return &Client{
		SecretKey:    secretKey,
		CategoryCode: categoryCode,
		IsSandbox:    isSandbox,
		HTTPClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) BaseURL() string {
	if c.IsSandbox {
		return "https://dev.toyyibpay.com"
	}
	return "https://toyyibpay.com"
}

type CreateBillParams struct {
	BillName            string
	BillDescription     string
	AmountRM            float64
	ExternalReferenceNo string
	CustomerName        string
	CustomerEmail       string
	CustomerPhone       string
	ReturnURL           string
	CallbackURL         string
}

type CreateBillResponse struct {
	BillCode string `json:"BillCode"`
}

// CreateBill requests ToyyibPay to generate a payment bill
func (c *Client) CreateBill(params CreateBillParams) (billCode string, paymentURL string, err error) {
	if c.SecretKey == "" || c.CategoryCode == "" {
		return "", "", errors.New("toyyibpay API Key atau Category Code belum ditetapkan dalam .env")
	}

	endpoint := fmt.Sprintf("%s/index.php/api/createBill", c.BaseURL())

	// ToyyibPay requires amount in Cents (Sen) -> RM 50.00 = 5000
	amountInSen := int(params.AmountRM * 100)

	formData := url.Values{}
	formData.Set("userSecretKey", c.SecretKey)
	formData.Set("categoryCode", c.CategoryCode)
	formData.Set("billName", params.BillName)
	formData.Set("billDescription", params.BillDescription)
	formData.Set("billPriceSetting", "1") // 1 = Fixed Amount
	formData.Set("billPayorInfo", "1")    // 1 = Required
	formData.Set("billAmount", strconv.Itoa(amountInSen))
	formData.Set("billReturnUrl", params.ReturnURL)
	formData.Set("billCallbackUrl", params.CallbackURL)
	formData.Set("billExternalReferenceNo", params.ExternalReferenceNo)
	formData.Set("billTo", params.CustomerName)
	formData.Set("billEmail", params.CustomerEmail)
	formData.Set("billPhone", params.CustomerPhone)
	formData.Set("billSplitPayment", "0")
	formData.Set("billPaymentChannel", "0") // 0 = Both FPX & Card

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("toyyibpay request error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("gagal membaca respon toyyibpay: %w", err)
	}

	// Respon ToyyibPay biasanya dalam bentuk JSON array: [{"BillCode":"xyz"}]
	var billResponses []CreateBillResponse
	if err := json.Unmarshal(bodyBytes, &billResponses); err != nil {
		return "", "", fmt.Errorf("respon toyyibpay tidak sah (%s): %w", string(bodyBytes), err)
	}

	if len(billResponses) == 0 || billResponses[0].BillCode == "" {
		return "", "", fmt.Errorf("gagal menjana bil ToyyibPay: %s", string(bodyBytes))
	}

	generatedCode := billResponses[0].BillCode
	payURL := fmt.Sprintf("%s/%s", c.BaseURL(), generatedCode)

	return generatedCode, payURL, nil
}
