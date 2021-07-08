package nodes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/marvin9002/swarm-install/beectl/service/config"
)

type BeeNode struct {
	ApiAddr         string
	P2PAddr         string
	DEBUGApiAddr    string
	Pwd             string
	EndPoint        string
	Dir             string
	BeeCfgPath      string // 节点配置文件路径
	NodeName        string
	NodeIp          string
	NodeStatus      int
	NodeVersion     string
	DebugApiPort    string
	sign            string
	serverUrl       string
	serverToken     string
	disableRegister bool
	pushInterval    time.Duration
	srv             *http.Server
	quit            chan struct{}
}

func NewNode(cfg config.Config) *BeeNode {
	return &BeeNode{
		ApiAddr:      cfg.Get(config.ApiAddr).(string),
		P2PAddr:      cfg.Get(config.P2PAddr).(string),
		DEBUGApiAddr: cfg.Get(config.DEBUGApiAddr).(string),
		DebugApiPort: strings.Split(cfg.Get(config.DEBUGApiAddr).(string), ":")[1],
		Pwd:          cfg.Get(config.Pwd).(string),
		EndPoint:     cfg.Get(config.EndPoint).(string),
		Dir:          cfg.Get(config.Dir).(string),
		BeeCfgPath:   cfg.GetCfg(),
		NodeName:     cfg.GetName(),
		NodeIp:       cfg.GetAddr(),
		NodeStatus:   0,
		NodeVersion:  "",
	}
}

// [node-name,node-ip,node-status,node-port,peers,node-version,diskavail,blance,free-blance,no-cashout,total-cashout,total-cash-blance]

func (this *BeeNode) GetName() string {
	return this.NodeName
}

func (this *BeeNode) GetIp() string {
	curl := exec.Command("curl", "-s", "api.infoip.io/ip") // 修改了此行
	out, err := curl.Output()
	if err != nil {
		fmt.Println("erorr", err)
		return ""
	}
	return string(out)
}

func (this *BeeNode) GetVersion() string {
	return this.NodeVersion
}

func (this *BeeNode) GetPort() string {
	return this.DEBUGApiAddr
}

func (this *BeeNode) debugApiReply(path string) ([]byte, error) {
	url, err := pathJoin(this.DEBUGApiAddr, path)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (this *BeeNode) debugApiPost(path string) ([]byte, error) {
	url, err := pathJoin(this.DEBUGApiAddr, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	hea := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(statuscode)
	fmt.Println(hea)
	return body, nil
}

func pathJoin(base, p string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}
