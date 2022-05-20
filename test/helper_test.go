package tests

import (
	"math/big"
	"testing"
	"utopia/internal/helper"

	"github.com/ethereum/go-ethereum/common"
)

func TestHelper(t *testing.T) {
	data := common.FromHex("0x1234567890abcdef")
	expect := []byte{18, 52, 86, 120, 144, 171, 205, 239}

	if len(data) != len(expect) {
		t.Errorf("Length not match")
		return
	}

	for i := 0; i < len(expect); i++ {
		if data[i] != expect[i] {
			t.Errorf("Data %v not match %v", data, expect)
			return
		}
	}

	value := helper.WeiToEth(big.NewInt(1234678954320000000))
	if value != float32(1.234679) {
		t.Errorf("Expect %f but %f", 1.234679, value)
		return
	}

	wei := helper.EthToWei(value)
	if wei.Cmp(big.NewInt(1234679000000000000)) != 0 {
		t.Errorf("Expect 1234678954320000000 but %d", wei.Uint64())
		return
	}
}
