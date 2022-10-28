package rank

import (
	"testing"
)

func TestGetGMSRank(t *testing.T) {
	rank, err := GetGMSRank("MoreCashFarm")
	if err != nil {
		t.Fail()
	}

	t.Logf("Got Rank: %v", rank)
}
