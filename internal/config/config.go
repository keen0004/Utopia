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

type SSLConfig struct {
	Cert     string `json:"cert"`
	Key      string `json:"key"`
	Password string `json:"password"`
}

type ServiceConfig struct {
	Port uint      `json:"port"`
	SSL  SSLConfig `json:"ssl"`
}

type ChainConfig struct {
	ChainListFile   string `json:"chainlist"`
	AccountListFile string `json:"accountlist"`
	Network         string `json:"network"`
	From            string `json:"from"`
}

type Configs struct {
	Server ServiceConfig `json:"service"`
	Chain  ChainConfig   `json:"chain"`
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
	err = chain.LoadChainList(config.Chain.ChainListFile)
	if err != nil {
		return err
	}

	// load account list
	err = wallet.LoadAccountList(config.Chain.AccountListFile)
	if err != nil {
		return err
	}

	return nil
}
