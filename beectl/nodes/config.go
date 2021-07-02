package nodes

import (
	"encoding/json"
	"fmt"
	"os"
)

func (h *Swarm) UpdateCfg(node, addr string) error {

	var cfg []config
	f, err := os.OpenFile(h.CtrlCfg, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return err
	}
	defer f.Close()

	// 创建json解码器
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		h.Println("Decoder failed", err.Error())
	} else {
		h.Println("Decoder success")
		h.Println(cfg)
	}
	cfg = append(cfg, config{
		Name: node,
		Addr: fmt.Sprintf("http://localhost:%s", addr),
		ID:   addr,
	})

	marshal, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(marshal))
	if err != nil {
		return err
	}
	return nil
}

func (h *Swarm) GetCfg() ([]config, error) {
	var cfg []config
	f, err := os.OpenFile(h.CtrlCfg, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return nil, err
	}
	defer f.Close()

	// 创建json解码器
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		h.Println("Decoder failed", err.Error())
		return nil, err
	}
	h.Println("Decoder success")
	h.Println(cfg)
	return cfg, nil
}
