package repository

import (
	"io"
	"log"
	"net/http"

	"n1h41/bybit-wallet-tui/configs"
	"n1h41/bybit-wallet-tui/dto"
	"n1h41/bybit-wallet-tui/utils"
)

type BybitRepository interface {
	GetAllCoinBalance()
	GetWalletBalance() dto.GetWalletBalanceResp
}

type bybitRepository struct {
	config configs.BybitConfig
	client *http.Client
}

// GetWalletBalance implements BybitRepository.
func (b *bybitRepository) GetWalletBalance() dto.GetWalletBalanceResp {
	endPoint := "/v5/account/wallet-balance"
	// params := "accountType=UNIFIED&coin=INJ"
	params := "accountType=UNIFIED"
	time_stamp := utils.GetTimestamp()
	signature := utils.CreateSignature(b.config.ApiSecret, b.config.ApiKey, b.config.RecvWindow, params)
	request, error := http.NewRequest("GET", b.config.Url+endPoint+"?"+params, nil)
	utils.AddAllHeaders(request, b.config.ApiKey, signature, time_stamp, b.config.RecvWindow)
	response, error := b.client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	getWalletBalanceResp, err := dto.UnmarshalGetWalletBalanceResp(body)
	if err != nil {
		log.Fatal(err)
	}
	return getWalletBalanceResp
}

// GetAllCoinBalance implements BybitRepository.
func (b *bybitRepository) GetAllCoinBalance() {
	endPoint := "/v5/asset/transfer/query-account-coins-balance"
	time_stamp := utils.GetTimestamp()
	// params := "accountType=UNIFIED&coin=USDT"
	params := "accountType=UNIFIED"
	signature := utils.CreateSignature(b.config.ApiSecret, b.config.ApiKey, b.config.RecvWindow, params)
	request, error := http.NewRequest("GET", b.config.Url+endPoint+"?"+params, nil)
	utils.AddAllHeaders(request, b.config.ApiKey, signature, time_stamp, b.config.RecvWindow)
	// reqDump, err := httputil.DumpRequestOut(request, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("Request Dump:\n%s", string(reqDump))
	response, error := b.client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	// log.Println("response Status:", response.Status)
	// log.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	log.Println("response Body:", string(body))
}

func NewBybitRepo(config configs.BybitConfig, httpClient *http.Client) BybitRepository {
	return &bybitRepository{
		config: config,
		client: httpClient,
	}
}
