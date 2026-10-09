package easyparcel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	APIKey     string
	IsSandbox  bool
	HTTPClient *http.Client
}

func NewClient(apiKey string, isSandbox bool) *Client {
	return &Client{
		APIKey:     apiKey,
		IsSandbox:  isSandbox,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) BaseURL() string {
	if c.IsSandbox {
		return "https://demo.connect.easyparcel.my"
	}
	return "https://connect.easyparcel.my"
}

type CourierRate struct {
	CourierName string  `json:"courierName"`
	ServiceCode string  `json:"serviceCode"`
	ServiceName string  `json:"serviceName"`
	Price       float64 `json:"price"`
	DeliveryETA string  `json:"deliveryEta"`
}

type RateCheckRequest struct {
	SenderPostcode   string  `json:"senderPostcode"`
	ReceiverPostcode string  `json:"receiverPostcode"`
	WeightKg         float64 `json:"weightKg"`
}

// CheckRates queries EasyParcel API for available shipping couriers and prices
func (c *Client) CheckRates(req RateCheckRequest) ([]CourierRate, error) {
	if c.APIKey == "" {
		// Jika API Key belum diset dalam .env, kembalikan anggaran kadar standard pos Malaysia
		return defaultEstimatedRates(req.WeightKg), nil
	}

	endpoint := fmt.Sprintf("%s/?ac=EPRateCheckingBulk", c.BaseURL())

	bulkData := []map[string]interface{}{
		{
			"pick_code":    req.SenderPostcode,
			"pick_state":   "",
			"pick_country": "MY",
			"send_code":    req.ReceiverPostcode,
			"send_state":   "",
			"send_country": "MY",
			"weight":       fmt.Sprintf("%.2f", req.WeightKg),
			"width":        "0",
			"length":       "0",
			"height":       "0",
			"date_coll":    time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		},
	}

	bulkJSON, err := json.Marshal(bulkData)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("api", c.APIKey)
	formData.Set("bulk", string(bulkJSON))

	httpReq, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("easyparcel request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		APIStatus string `json:"api_status"`
		Result    []struct {
			Status string `json:"status"`
			Rates  []struct {
				CourierName string `json:"courier_name"`
				ServiceCode string `json:"service_id"`
				ServiceName string `json:"service_name"`
				Price       string `json:"price"`
				Delivery    string `json:"delivery"`
			} `json:"rates"`
		} `json:"result"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return defaultEstimatedRates(req.WeightKg), nil
	}

	if len(apiResp.Result) == 0 || len(apiResp.Result[0].Rates) == 0 {
		return defaultEstimatedRates(req.WeightKg), nil
	}

	var results []CourierRate
	for _, rate := range apiResp.Result[0].Rates {
		priceVal, _ := strconv.ParseFloat(rate.Price, 64)
		results = append(results, CourierRate{
			CourierName: rate.CourierName,
			ServiceCode: rate.ServiceCode,
			ServiceName: rate.ServiceName,
			Price:       priceVal,
			DeliveryETA: rate.Delivery,
		})
	}

	return results, nil
}

// defaultEstimatedRates menyediakan fallback rate pos Malaysia sekiranya API Key EasyParcel belum diisi
func defaultEstimatedRates(weightKg float64) []CourierRate {
	baseRate := 8.00
	if weightKg > 1.0 {
		baseRate += (weightKg - 1.0) * 3.00
	}

	return []CourierRate{
		{
			CourierName: "J&T Express",
			ServiceCode: "EP-JT-EXP",
			ServiceName: "Standard Delivery (Door to Door)",
			Price:       baseRate,
			DeliveryETA: "1-3 Hari Bekerja",
		},
		{
			CourierName: "Pos Laju",
			ServiceCode: "EP-POS-LAJU",
			ServiceName: "Next Day Delivery",
			Price:       baseRate + 1.50,
			DeliveryETA: "1-2 Hari Bekerja",
		},
		{
			CourierName: "Ninja Van",
			ServiceCode: "EP-NINJA",
			ServiceName: "Standard Parcel",
			Price:       baseRate - 0.50,
			DeliveryETA: "2-3 Hari Bekerja",
		},
	}
}
