package main

import (
	"encoding/json"
	"io/ioutil"
)

type LM_T struct {
	From string `json:"from"`
	To4  string `json:"to"`
	To6  string `json:"to6"`
}

type Config struct {
	Router    []*LM_T `json:"router"`
	Protocol  string  `json:"p"`
	Mark      int     `json:"mark"`
	V         int     `json:"v"`
	Listeners int     `json:"listeners"`
	//CloseAfter int     `json:"close-after"`
}

func LoadConfig(filapath string, cfg *Config) error {
	content, err := ioutil.ReadFile(filapath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return err
	}
	return nil
}
