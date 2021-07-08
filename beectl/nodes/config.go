package nodes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func (h *Swarm) UpdateCfg(node, addr string) error {
	var cfg []config
	f, err := os.OpenFile(h.CtrlCfg, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return err
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &cfg)
		if err != nil {
			return err
		}
	}
	h.Println("old:", cfg)
	if len(cfg) > 0 {
		for _, conf := range cfg {
			if conf.ID != addr {
				cfg = append(cfg, config{
					Name: node,
					Addr: fmt.Sprintf("http://localhost:%s", addr),
					ID:   addr,
					Cfg:  path.Join(h.BeeCfgPath, fmt.Sprintf("/bee%s.yaml", addr)),
				})
			}
		}
	} else {
		cfg = append(cfg, config{
			Name: node,
			Addr: fmt.Sprintf("http://localhost:%s", addr),
			ID:   addr,
			Cfg:  path.Join(h.BeeCfgPath, fmt.Sprintf("/bee%s.yaml", addr)),
		})
	}

	h.Println(cfg)

	marshal, err := json.Marshal(cfg)
	if err != nil {
		h.Println(err)
		return err
	}

	if len(byteValue) > 0 {
		err = f.Truncate(0)
		if err != nil {
			return err
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString(string(marshal))
	if err != nil {
		h.Println(err)
		return err
	}
	f.Sync()
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
	byteValue, _ := ioutil.ReadAll(f)

	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &cfg)
		if err != nil {
			return nil, err
		}
	}
	h.Println("Decoder success")
	h.Println(cfg)
	return cfg, nil
}
