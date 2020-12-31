package adash

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

//
var MarkAsLostAfter = time.Minute * 2

type Instance struct {
	ID       string `json:"id"`
	AlgoType string `json:"atype"`
	Value
	Status
	CurrencyCode  string `json:"cur" redis:"currencyCode"`
	Country       string `json:"cntr" redis:"countryCode"`
	LogItemsCount int    `json:"lic"`
	Group         string `json:"g" redis:"g"`
}

// Value holds instance's last metrics.
type Value struct {
	Pnl                     string  `json:"-" redis:"pnl"`     //float64
	PositionTwap            string  `json:"pt" redis:"pt"`     //int
	PositionHedge           string  `json:"ph" redis:"ph"`     //int
	PositionOption          string  `json:"po" redis:"po"`     //int
	Bid                     string  `json:"b" redis:"b"`       //float64
	Ask                     string  `json:"a" redis:"a"`       //float64
	Last                    string  `json:"l" redis:"l"`       //float64
	OptionBid               string  `json:"ob" redis:"ob"`     //float64
	OptionAsk               string  `json:"oa" redis:"oa"`     //float64
	ShareAverageHedgePrice  string  `json:"sahp" redis:"sahp"` //float64
	TargetQuantityPerPeriod string  `json:"tqpp" redis:"tqpp"` //int
	TargetPositionFlowTwap  string  `json:"tpft" redis:"tpft"` //int
	RemainingShots          string  `json:"rs" redis:"rs"`     //int
	PnlFloat                float64 `json:"pnl" redis:"-"`
	PositionDelta           string  `json:"pd" redis:"pd"` // pd - global
	NumOfSecurities         string  `json:"ns" redis:"ns"` // ns -
	PositionDeltaLimit      string  `json:"pl" redis:"pl"` // pl - global
	PositionDeltaCash       string  `json:"pdc" redis:"pdc"`
	PositionDeltaLimitCash  string  `json:"plc" redis:"plc"` //plc
}

func reduceValueLength(s string, cnt int) string {
	dpot := strings.Index(s, ".")
	if dpot < 0 {
		return s
	}

	if len(s)-dpot > cnt {
		s = s[:dpot+cnt+1]
	}
	return s
}

func (v *Value) ReduceLength() {
	if v.Pnl == "" {
		v.PnlFloat = 0
	} else {
		var err error
		v.PnlFloat, err = strconv.ParseFloat(reduceValueLength(v.Pnl, 2), 64)
		if err != nil {
			// TODO: записать в лог.
			v.PnlFloat = math.NaN()
		}
	}
	v.PositionTwap = reduceValueLength(v.PositionTwap, 2)
	v.PositionHedge = reduceValueLength(v.PositionHedge, 2)
	v.PositionOption = reduceValueLength(v.PositionOption, 2)
	v.Bid = reduceValueLength(v.Bid, 3)
	v.Ask = reduceValueLength(v.Ask, 3)
	v.Last = reduceValueLength(v.Last, 3)
	v.OptionBid = reduceValueLength(v.OptionBid, 3)
	v.OptionAsk = reduceValueLength(v.OptionAsk, 3)
	v.ShareAverageHedgePrice = reduceValueLength(v.ShareAverageHedgePrice, 2)
	v.TargetQuantityPerPeriod = reduceValueLength(v.TargetQuantityPerPeriod, 2)
	v.TargetPositionFlowTwap = reduceValueLength(v.TargetPositionFlowTwap, 2)
	v.RemainingShots = reduceValueLength(v.RemainingShots, 2)
}

type Status struct {
	At     int64 `json:"at" redis:"at"`
	Status int8  `json:"status" redis:"st"`
}

/*type TradingPlatformConnectorState int8

const (
	INSTANTIATED TradingPlatformConnectorState = 0
	STARTED      TradingPlatformConnectorState = 1
	READY        TradingPlatformConnectorState = 2
	DOWN         TradingPlatformConnectorState = 3
	KILLED       TradingPlatformConnectorState = 4
	ERROR        TradingPlatformConnectorState = 5
)
*/

// InstanceStatus defines allowed list of instance statuses.
type InstanceStatus int

const (
	// READY to start running (waiting market open).
	READY int8 = 0
	// RUNNING normally.
	RUNNING int8 = 1

	// CLOSED manually by user or by algorithm.
	CLOSED int8 = 2

	// FAILED execution.
	FAILED int8 = 3

	// LOST says no status updates last 2 minutes.
	LOST int8 = 4
)

//
type Servicer interface {
	MatrixI() (map[string]map[string]string, error)
	Instances(result *[]InstanceEx) error
	PublishLaunchConfig(action, fname string, content string) (bool, error)
	Publish(code string) (bool, error)
	ListLen(key string) (int, error)
	ListRange(key string, from, to int) ([]string, error)
	GetObject(key string, result interface{}) error
	Config() map[string]string
	RefreshConfig() error
	PlatformStatus() (*PlatformStatus, error)
	FileStream(fname string) (io.ReadCloser, error)
	SaveFileToDisk(fname string, buf bytes.Buffer) error
}

// Storer is and interface devinition to data persistent storage.
type Storer interface {
	Scan(pattern string, result *[]string) error
	Keys(pattern string, result *[]string) error
	Get(key string, val interface{}) error
	MGet(keys []string, result *[]string) error
	Set(key string, val string) error
	Publish(code string) (bool, error)
	ListLen(key string) (int, error)
	ListRange(key string, from, to int) ([]string, error)
	GetObject(key string, result interface{}) error
	HGet(key, field string) (string, error)
}

const CMND = "CMND"

const CommandStopInstance = "STOP:"

const CommandStartInstance = "START:"

// InstanceLogItem describes single log entry to show in the Intance Console.
type InstanceLogItem struct {
	At       int64  `json:"at"`
	Severity int    `json:"severity"`
	Message  string `json:"msg"`
}

const KeyPrefixInput = "I:"
const KeyPrefixModel = "M:"
const KeyPrefixStatus = "S:"
const KeyPrefixLog = "L:"
const KeyPrefixLaunchConfig = "FILE:"
const KeyPortfolioPnL = "PNL:"
const ConnectorPrefix = "CS:"

//
const ScanInputKeys = KeyPrefixInput + "*"

type Connector struct {
	Name     string `json:"name"`
	Status   int    `json:"status"`
	StatusAt uint64 `json:"status_at,omitempty"`
}

type PlatformStatus struct {
	Connectors    []Connector `json:"connectors"`
	LogItemsCount int         `json:"log_items_count"`

	// Status holds general platform status
	Status             int8    `json:"status"`
	StatusAt           int64   `json:"status_at,omitempty"`
	PnL                float64 `json:"pnl"`
	PositionDelta      float64 `json:"pd"`
	PositionDeltaLimit float64 `json:"pl"`
	ServerTime         string  `json:"server_time"`
}
