package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

//func init() {
//	fPath, _ := os.Getwd()
//
//	fPath = path.Join(fPath, "config")
//
//	configPath := flag.String("c", fPath, "config file path")
//
//	flag.Parse()
//
//	fmt.Println(*configPath)
//
//	err := LoadConfigInformation(*configPath)
//
//	if err != nil {
//		panic(err)
//	}
//}

func LoadConfigInformation(configPath string) (err error) {

	var (
		filePath string

		wr string
	)

	if configPath == "" {

		wr, _ = os.Getwd()

		wr = path.Join(wr, "conf")

	} else {

		wr = configPath

	}

	filePath = path.Join(wr, "config_"+Environment+".yml")

	configData, err := ioutil.ReadFile(filePath)

	if err != nil {

		fmt.Printf(" config file read failed: %s", err)

		os.Exit(-1)

	}

	err = yaml.Unmarshal(configData, &ConfigInfo)

	if err != nil {

		fmt.Printf(" config parse failed: %s", err)

		os.Exit(-1)

	}

	return nil

}
