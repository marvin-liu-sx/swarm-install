package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/marvin9002/swarm-install/beectl/utils"
)

func (h *Swarm) Start() error {
	t := template.New("cfg")
	t = template.Must(t.Parse(
		`api-addr: 127.0.0.1:{{.ApiAddr}}
config: /etc/bee/bee{{.DEBUGApiAddr}}.yaml
data-dir: {{.Dir}}/bee{{.DEBUGApiAddr}}
debug-api-addr: 127.0.0.1:{{.DEBUGApiAddr}}
debug-api-enable: true
db-open-files-limit: 10000
p2p-addr: :{{.P2PAddr}}
full-node: true
clef-signer-enable: false
network-id: 1
mainnet: true
password: {{.Pwd}}
swap-endpoint: {{.EndPoint}}
swap-initial-deposit: 0
welcome-message: "本节点由Marvin.Pool 部署，本矿场专业提供bzz主机及批量部署系统，价格优惠 | 请联系wx:yongyuan900218"`))

	start := h.getStart()
	h.Println("Starting node ID :", start)
	for i := 0; i < h.Count; i++ {
		api := 33
		p2p := 34
		debug := 35
		number := start
		api = number*100 + api
		p2p = number*100 + p2p
		debug = number*100 + debug

		h.Println("debug:", debug)
		cfg := BeeCfg{
			ApiAddr:      strconv.Itoa(api),
			P2PAddr:      strconv.Itoa(p2p),
			DEBUGApiAddr: strconv.Itoa(debug),
			Pwd:          h.Pwd,
			EndPoint:     h.EndPoint,
			Dir:          h.Dir,
		}
		err := os.MkdirAll(h.BeeCfgPath, os.ModePerm)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path.Join(fmt.Sprintf(h.BeeCfgPath+"/bee%s.yaml", strconv.Itoa(debug))), os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			h.Println("create file: ", err)
			return err
		}
		err = t.Execute(f, cfg)
		if err != nil {
			h.Println("execute: ", err)
			return err
		}
		f.Close()

		h.Println(fmt.Sprintf("add node cfg success, debug-api-addr: 127.0.0.1:%s  config: /etc/bee/bee%s.yaml ", strconv.Itoa(debug), strconv.Itoa(debug)))
		time.Sleep(2 * time.Second)
		err = h.Run(fmt.Sprintf(h.BeeCfgPath+"/bee%s.yaml", strconv.Itoa(debug)), strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return err
		}
		h.Println(fmt.Sprintf("bee%s start success ", strconv.Itoa(debug)))
		// 提取地址
		err = h.ethereumAddr(strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return err
		}
		time.Sleep(2 * time.Second)
		// 备份钱包

		err = h.ExportSwarmKey(strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return err
		}

		// 批量提取地址
		time.Sleep(5 * time.Second)
		err = h.UpdateCfg(fmt.Sprintf("bee%s", strconv.Itoa(debug)), strconv.Itoa(debug))
		if err != nil {
			return err
		}
		start += 1
	}
	time.Sleep(3 * time.Second)
	err := h.GetAddrs()
	if err != nil {
		return err
	}
	return nil
}

func (h *Swarm) Run(cfg, addr string) error {
	h.Println(fmt.Sprintf("add node cfg success,   config:%s ", cfg))
	type Srv struct {
		Desc            string
		Start           string
		TimeoutStartSec int
		TimeoutStopSec  int
		TimeoutSec      int
	}

	t := template.New("srv")
	t = template.Must(t.Parse(
		`[Unit]
Description={{.Desc}}
Documentation=https://docs.ethswarm.org
After=network.target

[Service]
NoNewPrivileges=true
LimitNOFILE=655350

User=root
Group=root

ExecStart={{.Start}}
Restart=on-failure
RestartSec=2s

[Install]
WantedBy=multi-user.target
`))

	service := Srv{
		Desc:            "Bee - Ethereum Swarm node",
		Start:           "/usr/bin/bee start --config " + cfg,
		TimeoutStartSec: 120,
		TimeoutStopSec:  120,
		TimeoutSec:      120,
	}

	f, err := os.OpenFile(path.Join(fmt.Sprintf("/lib/systemd/system/bee%s.service", addr)), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return err
	}
	err = t.Execute(f, service)
	if err != nil {
		h.Println("execute: ", err)
		return err
	}
	defer f.Close()

	if !utils.ExecCommand("systemctl", []string{"daemon-reload"}) {
		h.Println("start bee error")
		return errors.New("systemctl daemon-reload error: exec: \"systemctl\": executable file not found in $PATH")
	}
	time.Sleep(1 * time.Second)
	if utils.ExecCommand("systemctl", []string{"start", fmt.Sprintf("bee%s", addr)}) {
		h.Println("start bee error")
		return errors.New("systemctl start bee error: exec: \"systemctl\": executable file not found in $PATH")
	}
	h.Println("start bee success ")
	return nil

}

func (h *Swarm) getStart() int {
	var cfg []config
	var start = 16
	f, err := os.OpenFile(h.CtrlCfg, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return start
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
	h.Println("Cfg len :", len(cfg))
	return start + len(cfg)
}
