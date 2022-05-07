package tests

import (
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"
	"utopia/internal/chain"
	"utopia/internal/helper"
)

var (
	blockNumber = uint64(14690688)
	account     = "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8"
	contract    = "0xba30e5f9bb24caa003e9f2f0497ad287fdf95623"
	blockHash   = "0xa2745988669bf11c8282ca6705ef92a177b9a397f5025c52e6140ce85465d2c0"
	txHash      = "0x1aa696883331e7dcbba571f74259a6a183ab63f827dda8d5c14156f024850506"
)

func connectChain() (chain.Chain, error) {
	c := chain.NewChain(1, "ETH", "eth")
	if c == nil {
		return nil, errors.New("New chain failed")
	}

	return c, c.Connect([]string{"https://rpc.ankr.com/eth"}, true)
}

func TestChainID(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	id, err := c.ChainId()
	if id.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Chain id expect %d but %d", 1, id.Int64())
	}
}

func TestGasPrice(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	gas, err := c.GasPrice()
	if gas.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("Gas price expect >0 but %d", gas.Int64())
	}
}

func TestBlockNumber(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	num, err := c.BlockNumber()
	if num <= 0 {
		t.Errorf("Block number expect >0 but %d", num)
	}
}

func TestBlockByNumber(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	b, err := c.BlockByNumber(blockNumber)
	if b.Number().Uint64() != blockNumber {
		t.Errorf("Block by number expect %d but %d", blockNumber, b.Number().Uint64())
	}
}

func TestBlockByHash(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	b, err := c.BlockByHash(helper.Str2bytes(blockHash))
	if b.Hash().Hex() != blockHash {
		t.Errorf("Block by hash expect %s but %s", blockHash, b.Hash().Hex())
	}
}

func TestTransaction(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	tx, _ := c.Transaction(helper.Str2bytes(txHash))
	if tx.Hash().Hex() != txHash {
		t.Errorf("Transaction expect %s but %s", txHash, tx.Hash().Hex())
	}
}

func TestReceipt(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	r, err := c.Receipt(helper.Str2bytes(txHash))
	if r.TxHash.Hex() != txHash {
		t.Errorf("Transaction expect %s but %s", txHash, r.TxHash.Hex())
	}
}

func TestSendTransaction(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	// todo
}

func TestEstimateGas(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	tx, _ := c.Transaction(helper.Str2bytes(txHash))
	if tx.Hash().Hex() != txHash {
		t.Errorf("Transaction expect %s but %s", txHash, tx.Hash().Hex())
	}

	r, err := c.Receipt(helper.Str2bytes(txHash))
	if r.TxHash.Hex() != txHash {
		t.Errorf("Transaction expect %s but %s", txHash, r.TxHash.Hex())
	}

	gas, err := c.EstimateGas(tx)
	if gas <= 0 && !strings.HasPrefix(err.Error(), "execution reverted") {
		t.Errorf("Estimate gas expect >0 but %d %v", gas, err)
	}
}

func TestBalance(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	balance, err := c.Balance(account)
	if helper.WeiToEth(&balance) <= 0 {
		t.Errorf("Balance expect >0 but %f", helper.WeiToEth(&balance))
	}
}

func TestNonce(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	nonce, err := c.Nonce(account)
	if nonce <= 0 {
		t.Errorf("Nonce expect >0 but %d", nonce)
	}
}

func TestCode(t *testing.T) {
	c, err := connectChain()
	if err != nil {
		t.Errorf("New eth chain failed with error: %v", err)
		return
	}
	defer c.DisConnect()

	code, err := c.Code(contract)
	if len(code) <= 0 {
		t.Errorf("Expect code size >0 but 0")
	}
}

func TestMeta(t *testing.T) {
	err := chain.AddChainMeta(chain.ETH_MAINNET, "eth", "ETH", false, []string{"https://rpc.ankr.com/eth"}, "https://etherscan.io")
	if err != nil {
		t.Errorf("Add chain meta info failed with err %v", err)
		return
	}

	err = chain.AddChainMeta(chain.BSC_MAINNET, "bsc", "BNB", false, []string{"https://rpc.ankr.com/bsc"}, "https://bscscan.com")
	if err != nil {
		t.Errorf("Add chain meta info failed with err %v", err)
		return
	}

	path := "./chainlist.json"
	err = chain.SaveChainList(path)
	if err != nil {
		t.Errorf("Save chain meta info failed with err %v", err)
		return
	}
	defer os.Remove(path)

	err = chain.ReloadChainList(path)
	if err != nil {
		t.Errorf("Load chain meta info failed with err %v", err)
		return
	}

	meta, err := chain.ChainMetaById(chain.ETH_MAINNET)
	if meta.Name != "eth" {
		t.Errorf("Expect eth but %s", meta.Name)
		return
	}

	meta, err = chain.ChainMetaByName("bsc")
	if meta.Id != 56 {
		t.Errorf("Expect 56 but %d", meta.Id)
		return
	}
}
