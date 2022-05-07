package chain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type ChainMeta struct {
	Id        uint64   // Chain id
	Name      string   // Chain name in full mode
	Currency  string   // Currency name
	IsTest    bool     // Is a test network
	RpcServer []string // The rpc url list for call
	Explorer  string   // The block chain explorer
}

var (
	ChainList       []ChainMeta
	ChainIdMap      map[uint64]int
	ChainNameMap    map[string]int
	ChainListConfig = "./chainlist.json"
)

func init() {
	ChainList = make([]ChainMeta, 0)
	ChainIdMap = make(map[uint64]int)
	ChainNameMap = make(map[string]int)

	_, err := os.Stat(ChainListConfig)
	if err != nil {
		return
	}

	err = ReloadChainList(ChainListConfig)
	if err != nil {
		return
	}

	return
}

func ChainMetaById(id uint64) (*ChainMeta, error) {
	index, ok := ChainIdMap[id]
	if !ok {
		return nil, errors.New("Chain id is not exist")
	}

	return &ChainList[index], nil
}

func ChainMetaByName(name string) (*ChainMeta, error) {
	index, ok := ChainNameMap[name]
	if !ok {
		return nil, errors.New("Chain name is not exist")
	}

	return &ChainList[index], nil
}

func AddChainMeta(id uint64, name string, currency string, isTest bool, server []string, explorer string) error {
	_, ok := ChainIdMap[id]
	if ok {
		return errors.New("Chain id is exist")
	}

	_, ok = ChainNameMap[name]
	if ok {
		return errors.New("Chain name is exist")
	}

	ChainList = append(ChainList, ChainMeta{
		Id:        id,
		Name:      name,
		Currency:  currency,
		IsTest:    isTest,
		RpcServer: server,
		Explorer:  explorer,
	})

	ChainIdMap[id] = len(ChainList) - 1
	ChainNameMap[name] = len(ChainList) - 1

	return nil
}

func ReloadChainList(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	meta := make([]ChainMeta, 0)
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return err
	}

	for _, m := range meta {
		AddChainMeta(m.Id, m.Name, m.Currency, m.IsTest, m.RpcServer, m.Explorer)
	}

	return nil
}

func SaveChainList(path string) error {
	data, err := json.MarshalIndent(ChainList, "", "    ")
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(path, data, 0666)
}
