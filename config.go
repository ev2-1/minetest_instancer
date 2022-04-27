package main

import (
	"encoding/json"
	"errors"
	"os"

	"fmt"
)

var srvNetwork = "ev21_defaul"
var minetestContainer = "minetest_5.4.1"
var worldPath = "/home/derz/Documents/ev2.1/worlds/"
var worldMount = "/mount/worlds/"
var configPath = "/home/derz/Documents/ev2.1/config/"
var gamePath = "/home/derz/Documents/ev2.1/games/"

type conf struct {
	SrvNetwork        string
	MinetestContainer string
	WorldPath         string
	WorldMount        string
	ConfigPath        string
	GamePath          string
}

var config *conf

func ReadConfig() error {
	f, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.Size() == 0 {
		f.WriteString("{\n\t\n}\n")
		f.Seek(0, os.SEEK_SET)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	// check if config is valid:
	fmt.Println(config.GamePath)
	if len(config.GamePath) == 0 {
		return errors.New("config field without value")
	}

	fmt.Println("[CONFIG] load")

	return nil
}
