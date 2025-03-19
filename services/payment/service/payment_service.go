package service

import (
	"fmt"
	"log"
	"time"
)

// PaymentService 處理支付相關業務邏輯
type PaymentService struct {
	// 支付閘道整合
	ApiKey     string
	MerchantID string
	IsTestMode bool
}

// NewPaymentService 創建支付服務實例
func NewPaymentService(apiKey, merchantID string, isTestMode bool) *PaymentService {
	return &PaymentService{
		ApiKey:     apiKey,
		MerchantID: merchantID,
		IsTestMode: isTestMode,
	}
}

// ProcessPayment 處理支付請求
func ProcessPayment(amount float64) (bool, error) {
	// 支付處理邏輯
	log.Printf("處理支付請求: $%.2f", amount)
	
	// 模擬支付處理延遲
	time.Sleep(500 * time.Millisecond)
	
	// 簡單模擬 - 金額大於 0 且小於 10000 則支付成功
	if amount > 0 && amount < 10000 {
		log.Println("支付成功")
		return true, nil
	}
	
	log.Println("支付失敗")
	return false, fmt.Errorf("支付失敗")
}
