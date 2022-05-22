package config

import (
	"encoding/json"
	"io/ioutil"
	"utopia/internal/chain"
	"utopia/internal/wallet"
)

var (
	DEFAULT_CONFIG_FILE = "../../configs/utopia.json"
	Config              Configs
)

type Configs struct {
	ChainListFile   string `json:"chainlist"`
	AccountListFile string `json:"accountlist"`
	Network         string `json:"network"`
	From            string `json:"from"`
}

func (config *Configs) LoadConfig(path string) error {
	// load configs
	if path == "" {
		path = DEFAULT_CONFIG_FILE
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	// load network list
	err = chain.LoadChainList(config.ChainListFile)
	if err != nil {
		return err
	}

	// load account list
	err = wallet.LoadAccountList(config.AccountListFile)
	if err != nil {
		return err
	}

	return nil
}
