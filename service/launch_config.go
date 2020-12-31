package adash

import (
	"bufio"
	"strings"
)

type LaunchConfig struct {
}

func Parse(buf []byte, lc *LaunchConfig) error {

	scanner := bufio.NewReader(strings.NewReader(string(buf)))
	for {
		_, err := scanner.ReadString('\n')
		if err != nil {

		}
	}

	return nil
}

type LaunchConfigSingleAlgoType struct {
	AlgoType  string // BUOYANCY
	Marker    string // EUROPE:300
	CSVHeader string //
	CSVRow    []string
}

type InstanceInputs struct {
	AlgoProductCode        string `json:"algoProductCode"`
	DerivativeBrokerCode   string `json:"derivativeBrokerCode"`
	EquityBrokerCode       string `json:"equityBrokerCode"`
	EquityAccount          string `json:"equityAccount"`
	DerivativeAccount      string `json:"derivativeAccount"`
	ProductID              string `json:"productId"`
	AllowedShots           string `json:"allowedShots"`
	PositionTwap           string `json:"positionTwap"`
	PositionCostTwap       string `json:"positionCostTwap"`
	StartTime              string `json:"startTime"`
	EndTime                string `json:"endTime"`
	TwapTimeRange          string `json:"twapTimeRange"`
	RandomInterval         string `json:"randomInterval"`
	Threshold              string `json:"threshold"`
	PositionOption         string `json:"positionOption"`
	OptionPositionCost     string `json:"optionPositionCost"`
	PositionHedge          string `json:"positionHedge"`
	StopLossLimit          string `json:"stopLossLimit"`
	StopProfitLimit        string `json:"stopProfitLimit"`
	PositionCostHedge      string `json:"positionCostHedge"`
	InterestRate           string `json:"interestRate"`
	Strike                 string `json:"strike"`
	ShortCodeShare         string `json:"shortCodeShare"`
	CountryCode            string `json:"countryCode"`
	ExpiryDate             string `json:"expiryDate"`
	OverriddenOpeningPrice string `json:"overriddenOpeningPrice"`
	SectorIndex            string `json:"sectorIndex"`
	CurrencyCode           string `json:"currencyCode"`
	OpenTime               string `json:"openTime"`
	OptionCode4Threshold   string `json:"optionCode4Threshold"`
}
