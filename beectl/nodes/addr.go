package nodes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func (h *Swarm) ethereumAddr(port string) error {
	url := fmt.Sprintf("http://localhost:%s/addresses", port)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	var addr *address

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &addr)
	if err != nil {
		return err
	}

	h.Println(fmt.Sprintf("[ethereum address] (非常重要)钱包地址:  %s \n", addr.Ethereum))

	return nil
}

func (h *Swarm) GetAddrs() error {
	var cfg []config
	cfg, err := h.GetCfg()
	if err != nil {
		return err
	}
	if len(cfg) <= 0 {
		return nil
	}
	var addressesArr []string
	for _, conf := range cfg {
		addresses, err := h.getAddr(conf.ID)
		if err != nil {
			return err
		}
		h.Println(addresses)
		addressesArr = append(addressesArr, addresses)
	}
	h.Println("所有节点地址信息已经存入当前目录的 addresses.log 中，方便您后期查看")
	logs, err := os.OpenFile(path.Join(h.InstallPath, "/addresses.log"), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return err
	}

	defer logs.Close()
	Logger := log.New(logs, "", 0)
	Logger.Println(addressesArr)
	return nil
}

func (h *Swarm) getAddr(port string) (string, error) {
	url := fmt.Sprintf("http://localhost:%s/addresses", port)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	type address struct {
		Overlay      string   `json:"overlay"`
		Underlay     []string `json:"underlay"`
		Ethereum     string   `json:"ethereum"`
		PublicKey    string   `json:"publicKey"`
		PssPublicKey string   `json:"pssPublicKey"`
	}

	var addr *address

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &addr)
	if err != nil {
		return "", err
	}

	return addr.Ethereum, nil
}
