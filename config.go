package flowride

import (
	"encoding/json"
	"os"
	"io/ioutil"
)

type FlConf struct {
	FlCapConf `json:"flcap"`
	FlMkConf `json:"flmk"`
}

type FlCapConf struct {
	DpdkInPort uint `json:"dpdkInPort"`
	DpdkOutPort uint `json:"dpdkOutPort"`
	CpuList     string `json:"cpuList"`
	DPDKArgs    []string `json:"dpdkArgs"`
}

type FlMkConf struct {
	IP string `json:"ip"`
	Port uint32 `json:"port"`
} 

func ConfigFromFileName(a string) (FlConf, error) {
	var ret FlConf
	file, err := os.Open(a) 
	if err != nil {
		return ret, err
	}
	confBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(confBytes, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
