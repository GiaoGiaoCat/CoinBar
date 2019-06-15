package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/caseymrm/menuet"
)

func main() {
	go helloClock()
	menuet.App().RunApplication()
}

func helloClock() {
	for {
		// 如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
		resp := ContractMarketDetailMerged("EOS_NW")
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "Sell " + resp["ask"] + " - Buy " + resp["bid"],
		})
		time.Sleep(time.Second)
	}
}

func ContractMarketDetailMerged(symbol string) (result map[string]string) {
	result = map[string]string{"ask": "0", "bid": "0"}
	resp, err := http.Get("https://api.hbdm.com/market/depth?symbol=" + symbol + "&type=step0")
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result
	}

	var depthReturn DepthReturn
	err = json.Unmarshal(body, &depthReturn)
	if err != nil {
		return result
	}
	result["ask"] = fmt.Sprintf("%.3f", depthReturn.Tick.Asks[0][0])
	result["bid"] = fmt.Sprintf("%.3f", depthReturn.Tick.Bids[0][0])
	return result
}

type DepthReturn struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Tick   struct {
		Asks    [][]float64 `json:"asks"`
		Bids    [][]float64 `json:"bids"`
		Ch      string      `json:"ch"`
		ID      int         `json:"id"`
		Mrid    int         `json:"mrid"`
		Ts      int64       `json:"ts"`
		Version int         `json:"version"`
	} `json:"tick"`
	Ts int64 `json:"ts"`
}
