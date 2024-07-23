// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getWalletBalanceResp, err := UnmarshalGetWalletBalanceResp(bytes)
//    bytes, err = getWalletBalanceResp.Marshal()

package dto

import (
	"encoding/json"
)

func UnmarshalGetWalletBalanceResp(data []byte) (GetWalletBalanceResp, error) {
	var r GetWalletBalanceResp
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetWalletBalanceResp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetWalletBalanceResp struct {
	RetCode    int64      `json:"retCode"`
	RetMsg     string     `json:"retMsg"`
	Result     Result     `json:"result"`
	RetEXTInfo RetEXTInfo `json:"retExtInfo"`
	Time       int64      `json:"time"`
}

type Result struct {
	List []List `json:"list"`
}

type List struct {
	TotalEquity            string `json:"totalEquity"`
	AccountIMRate          string `json:"accountIMRate"`
	TotalMarginBalance     string `json:"totalMarginBalance"`
	TotalInitialMargin     string `json:"totalInitialMargin"`
	AccountType            string `json:"accountType"`
	TotalAvailableBalance  string `json:"totalAvailableBalance"`
	AccountMMRate          string `json:"accountMMRate"`
	TotalPerpUPL           string `json:"totalPerpUPL"`
	TotalWalletBalance     string `json:"totalWalletBalance"`
	AccountLTV             string `json:"accountLTV"`
	TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
	Coin                   []Coin `json:"coin"`
}

type Coin struct {
	AvailableToBorrow   string `json:"availableToBorrow"`
	Bonus               string `json:"bonus"`
	AccruedInterest     string `json:"accruedInterest"`
	AvailableToWithdraw string `json:"availableToWithdraw"`
	TotalOrderIM        string `json:"totalOrderIM"`
	Equity              string `json:"equity"`
	TotalPositionMM     string `json:"totalPositionMM"`
	UsdValue            string `json:"usdValue"`
	SpotHedgingQty      string `json:"spotHedgingQty"`
	UnrealisedPnl       string `json:"unrealisedPnl"`
	CollateralSwitch    bool   `json:"collateralSwitch"`
	BorrowAmount        string `json:"borrowAmount"`
	TotalPositionIM     string `json:"totalPositionIM"`
	WalletBalance       string `json:"walletBalance"`
	CumRealisedPnl      string `json:"cumRealisedPnl"`
	Locked              string `json:"locked"`
	MarginCollateral    bool   `json:"marginCollateral"`
	Coin                string `json:"coin"`
}

type RetEXTInfo struct{}
