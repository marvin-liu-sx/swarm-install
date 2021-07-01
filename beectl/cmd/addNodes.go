/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/marvin9002/swarm-install/beectl/utils"

	"github.com/spf13/cobra"
)

// addNodesCmd represents the addNodes command
var addNodesCmd = &cobra.Command{
	Use:   "addNodes",
	Short: "添加bee节点",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			return err
		}
		dir, err := cmd.Flags().GetString("data-dir")
		if err != nil {
			return err
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}

		ports, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}
		swapEndpoint, err := cmd.Flags().GetString("swap-endpoint")
		if err != nil {
			return err
		}

		return install(count, dir, password, ports, swapEndpoint)
	},
}

func init() {
	rootCmd.AddCommand(addNodesCmd)
	addNodesCmd.Flags().IntP("count", "", 1, "要添加几个节点 ")
	addNodesCmd.Flags().StringP("data-dir", "", "/root/", "data-dir (required) bee数据存放目录")
	addNodesCmd.Flags().StringP("password", "", "", "[可选] bee启动密码,不输入将随机生成")
	addNodesCmd.Flags().StringP("ports", "", "", "[可选] 必须用逗号分隔输入如: 1633,1634,1635 分别对应 api-addr/p2p-addr/debug-api-addr的端口(如输入只会创建一个节点) 注意不要端口重复,默认自动生成端口")
	addNodesCmd.Flags().StringP("swap-endpoint", "", "", "swap-endpoint address (required)")
}

func install(count int, dir, pwd, ports, endpoint string) error {
	if count <= 0 {
		return errors.New("节点数错误")
	}
	if dir == "" {
		return errors.New("数据目录错误")
	}
	if endpoint == "" {
		return errors.New("endpoint 错误")
	}

	f, err := os.OpenFile("/root/bee/result.log", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return err
	}

	defer f.Close()
	handle := &InstallSwarm{
		Logger:   log.New(f, "[InstallSwarm]: ", log.Ldate|log.Ltime),
		Count:    count,
		Ports:    ports,
		Pwd:      pwd,
		EndPoint: endpoint,
		Dir:      dir,
	}
	err = handle.InstallBee()
	if err != nil {
		return err
	}

	handle.buildCfg()

	return nil
}

type InstallSwarm struct {
	Logger   *log.Logger
	Count    int
	Dir      string
	Ports    string
	Pwd      string
	EndPoint string
}

func (h *InstallSwarm) InstallBee() error {
	h.Println("InstallBee ....")
	_, err := exec.LookPath("/usr/bin/bee")
	if err == nil {
		h.Println(err)
		return nil
	}
	h.Println("bee install ....")
	switch runtime.GOOS {
	case "darwin":
	case "windows":
	case "linux":
		return h.InstallBeeLinux()
	}
	return nil
}

func (h *InstallSwarm) InstallBeeLinux() error {
	h.Println(runtime.GOARCH)
	switch runtime.GOARCH {
	case "amd64":
		h.Println("InstallBeeLinux ......")
		if !utils.Exist("./bee_1.0.0_amd64.deb") {
			if utils.ExecCommand("wget", []string{"https://github.com/ethersphere/bee/releases/download/v1.0.0/bee_1.0.0_amd64.deb"}) {
				utils.ExecCommand("sudo dpkg -i bee_1.0.0_amd64.deb", nil)
			} else {
				h.Println("InstallBeeLinux  err......")
				return errors.New("Install bee error")
			}
		}
		if !utils.ExecCommand("dpkg", []string{"-i", "bee_1.0.0_amd64.deb"}) {
			h.Println("install error")
		}
	case "armv7":
	case "arm64":
	}
	h.Println("install bee success")
	return nil
}

type BeeCfg struct {
	ApiAddr      string
	P2PAddr      string
	DEBUGApiAddr string
	Pwd          string
	EndPoint     string
	Dir          string
}

func (h *InstallSwarm) buildCfg() {
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
		err := os.MkdirAll("/etc/bee", os.ModePerm)
		if err != nil {
			return
		}
		f, err := os.OpenFile(path.Join(fmt.Sprintf("/etc/bee/bee%s.yaml", strconv.Itoa(debug))), os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			h.Println("create file: ", err)
			return
		}
		err = t.Execute(f, cfg)
		if err != nil {
			h.Println("execute: ", err)
			return
		}
		f.Close()

		h.Println(fmt.Sprintf("add node cfg success, debug-api-addr: 127.0.0.1:%s  config: /etc/bee/bee%s.yaml ", strconv.Itoa(debug), strconv.Itoa(debug)))
		time.Sleep(2 * time.Second)
		err = h.Start(fmt.Sprintf("/etc/bee/bee%s.yaml", strconv.Itoa(debug)), strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return
		}
		h.Println(fmt.Sprintf("bee%s start success ", strconv.Itoa(debug)))
		// 提取地址
		err = h.ethereumAddr(strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return
		}
		time.Sleep(2 * time.Second)
		// 备份钱包

		err = h.ExportSwarmKey(strconv.Itoa(debug))
		if err != nil {
			h.Println(err)
			return
		}

		// 批量提取地址
		time.Sleep(5 * time.Second)
		h.UpCFG(fmt.Sprintf("bee%s", strconv.Itoa(debug)), strconv.Itoa(debug))
		start += 1
	}
	time.Sleep(3 * time.Second)
	h.GetAddrs()
}

func (h *InstallSwarm) Start(cfg, addr string) error {
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
EnvironmentFile=-/etc/default/bee
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
	utils.ExecCommand("systemctl", []string{"start", fmt.Sprintf("bee%s", addr)})
	h.Println("start bee success ")
	return nil

}

func (h *InstallSwarm) ethereumAddr(port string) error {
	url := fmt.Sprintf("http://localhost:%s/addresses", port)

	resp, err := http.Get(url)
	if err != nil {
		return err
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
		return err
	}

	err = json.Unmarshal(body, &addr)
	if err != nil {
		return err
	}

	h.Println(fmt.Sprintf("[ethereum address] (非常重要)钱包地址:  %s \n", addr.Ethereum))

	return nil
}

func (h *InstallSwarm) ExportSwarmKey(addr string) error {
	h.Println(fmt.Sprintf("可通过命令查看此bee节点的日志(非常重要): journalctl -f -u bee%s \n", addr))

	f1, err := os.Open(path.Join(h.Dir, fmt.Sprintf("/bee%s", addr)))
	if err != nil {
		h.Println(err)
		return err
	}
	defer f1.Close()
	f2, err := os.Open(path.Join("/etc/bee", fmt.Sprintf("/bee%s.yaml", addr)))
	if err != nil {
		h.Println(err)
		return err
	}
	defer f2.Close()

	var files = []*os.File{f1, f2}
	dest := path.Join(h.Dir, fmt.Sprintf("/bkup_bee%ss_keys.zip", addr))
	err = utils.Compress(files, dest)
	if err != nil {
		h.Println(err)
		return err
	}

	h.Printf("已经打包备份钱包相关信息在:%s \n", dest)
	h.Println("!!! 请注意备份钱包相关信息目录文件, 请将备份后的文件保存到其他电脑上以防止本机硬盘损坏以后钱包丢失!!!")
	h.Println("当前输出信息已经存入当前目录的 result.log 中，方便您后期查看")
	return nil
}

type Cfg struct {
	node string
	addr string
}

func (h *InstallSwarm) GetAddrs() {
	var cfg []Cfg
	f, err := os.OpenFile("/root/bee/beectrl.cfg", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return
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
	var addressesArr []string
	for _, conf := range cfg {
		addresses, err := h.getAddr(conf.addr)
		if err != nil {
			return
		}
		addressesArr = append(addressesArr, addresses)
	}
	h.Println("所有节点地址信息已经存入当前目录的 addresses.log 中，方便您后期查看")
	logs, err := os.OpenFile("/root/bee/addresses.log", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	defer logs.Close()
	Logger := log.New(logs, "", 0)
	Logger.Println(addressesArr)
}

func (h *InstallSwarm) getAddr(port string) (string, error) {
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

func (h *InstallSwarm) UpCFG(node, addr string) {

	var cfg []Cfg
	f, err := os.OpenFile("/root/bee/beectrl.cfg", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		h.Println("create file: ", err)
		return
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
	cfg = append(cfg, Cfg{
		node: node,
		addr: addr,
	})

	marshal, err := json.Marshal(cfg)
	if err != nil {
		return
	}

	_, err = f.WriteString(string(marshal))
	if err != nil {
		return
	}
}

func (h *InstallSwarm) getStart() int {
	var cfg []Cfg
	var start = 16
	f, err := os.OpenFile("/root/bee/beectrl.cfg", os.O_RDWR|os.O_CREATE, 0666)
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

func (h *InstallSwarm) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
	h.Logger.Printf(format, v)
}

func (h *InstallSwarm) Println(v ...interface{}) {
	log.Println(v)
	h.Logger.Println(v)
}
