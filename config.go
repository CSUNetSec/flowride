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
	DpdkInPorts []string `json:"dpdkInPorts"`
	DpdkOutPorts []string `json:"dpdkOutPorts"`
	CpuList     string `json:"cpuList"`
	DPDKArgs    []string `json:"dpdkArgs"`
	Profiler    bool `json:"profiler"`
}

type FlMkConf struct {
	IP string `json:"ip"`
	Port uint32 `json:"port"`
} 

func ConfigFromFileName(a string) (FlConf, error) {
	var ret FlConf
	// by default keep the profiler off.
	ret.Profiler = false
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
