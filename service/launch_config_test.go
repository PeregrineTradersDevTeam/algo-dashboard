package adash

import (
	"testing"
)

func TestParse(t *testing.T) {

	var file1 = []byte(
		`
BUOYANCY
EUROPE:300
algoProductCode;derivativeBrokerCode;equityBrokerCode;equityAccount;derivativeAccount;productId;allowedShots;positionTwap;positionCostTwap;startTime;endTime;twapTimeRange;randomInterval;threshold;positionOption;optionPositionCost;positionHedge;stopLossLimit;stopProfitLimit;positionCostHedge;interestRate;strike;shortCodeShare;countryCode;expiryDate;overriddenOpeningPrice;sectorIndex;currencyCode;openTime;optionCode4Threshold
T001;BMTB;BMTB;EUR:PEREGRINE;EUR:PEREGRINE;ABN NA Equity;750;0;0;05-12-2019 09:30:00;05-12-2019 17:30:00;1;20;0.0278081270658481;0;0;0;-100;10;0;0;15.5;ABN ;NA;20-12-2019 17:30:00;-1;BEBANKS Index;EUR;05-12-2019 09:00:00;ABN NA 12/20/19 C15.5  Equity
T002;BMTB;BMTB;EUR:PEREGRINE;EUR:PEREGRINE;AD NA Equity;750;0;0;05-12-2019 09:30:00;05-12-2019 17:30:00;1;20;0.0429183490430992;0;0;0;-100;10;0;0;23.5;AH ;NA;20-12-2019 17:30:00;-1;BEFOODR Index;EUR;05-12-2019 09:00:00;AH NA 12/20/19 C23.5  Equity
T003;BMTB;BMTB;EUR:PEREGRINE;EUR:PEREGRINE;INGA NA Equity;750;0;0;05-12-2019 09:30:00;05-12-2019 17:30:00;1;20;0.0194574912090999;0;0;0;-100;10;0;0;10.5;ING ;NA;20-12-2019 17:30:00;-1;BEBANKS Index;EUR;05-12-2019 09:00:00;ING NA 12/20/19 C10.5  Equity
T004;BMTB;BMTB;EUR:PEREGRINE;EUR:PEREGRINE;MT NA Equity;750;0;0;05-12-2019 09:30:00;05-12-2019 17:30:00;1;20;0.0234977854782009;0;0;0;-100;10;0;0;15.5;MT ;NA;20-12-2019 17:30:00;-1;BWIRON Index;EUR;05-12-2019 09:00:00;MT NA 12/20/19 C15.5  Equity
T005;BMTB;BMTB;EUR:PEREGRINE;EUR:PEREGRINE;NN NA Equity;750;0;0;05-12-2019 09:30:00;05-12-2019 17:30:00;1;20;0.0642094119546773;0;0;0;-100;10;0;0;35;NN ;NA;20-12-2019 17:30:00;-1;BIELINSP Index;EUR;05-12-2019 09:00:00;NN NA 12/20/19 C35  Equity
END
`)
	var lc LaunchConfig

	err := Parse(file1, &lc)
	if err != nil {
		t.Fatal("Parse failed:", err.Error())
	}
}
