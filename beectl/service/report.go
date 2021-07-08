package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/marvin9002/swarm-install/beectl/utils"

	"github.com/marvin9002/swarm-install/beectl/service/config"

	"github.com/marvin9002/swarm-install/beectl/service/nodes"
)

type Report struct {
	NodeName               string // 节点名称
	NodeIP                 string // 节点IP
	NodeStatus             string // 节点状态
	NodePort               string // 节点端口
	Peers                  int
	NodeVersion            string // 节点版本
	UncashedAmount         int    // 未兑现金额
	TotalCheque            int    // 总票数
	TotalChequeBalance     string // 总额
	ChequeAvailableBalance string //  可用余额
	TotalReceived          string
	UseDisk                string // 磁盘使用
}

// [node-name,node-ip,node-status,node-port,peers,node-version,diskavail,blance,free-blance,no-cashout,total-cashout,total-cash-blance]
func (s *Service) Report(cfg config.Config) {
	s.Println("Start Report...")
	node := nodes.NewNode(cfg)
	t := time.NewTicker(s.ReportIntval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//t.Reset(s.ReportIntval)
			totalBalance, availableBalance := node.GetBalance()
			var size int64
			size, err2 := utils.DirSizeB(cfg.Get(config.Dir).(string))
			if err2 != nil {
				return
			}

			sizeMB := size / 1024 / 1024
			s.Println("sizeMB:", sizeMB)
			report := &Report{
				NodeName:               node.GetName(),
				NodeIP:                 node.GetIp(),
				NodeStatus:             node.GetStatus(),
				NodePort:               cfg.GetID(),
				Peers:                  node.GetPeers(),
				NodeVersion:            "v1.0.0",
				UncashedAmount:         node.GetUncashedAmount(),
				TotalCheque:            node.GetTotalCheque(),
				TotalChequeBalance:     totalBalance,
				ChequeAvailableBalance: availableBalance,
				TotalReceived:          node.GetSettlements(),
				UseDisk:                fmt.Sprintf("%sMB", strconv.Itoa(int(sizeMB))),
			}
			s.Println(fmt.Sprintf("Report-Data:%#v", report))
			err := s.report(report)
			if err != nil {
				s.Println("Report-Error:", err)
				continue
			}
			s.Println("Report-Hash:", report)
		}
	}

}

func (this *Service) report(report *Report) error {
	url, err := pathJoin(this.ReportServer, "")
	this.Println("report-url:", url)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(report)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshal))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	hea := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	this.Println(string(body))
	this.Println(statuscode)
	this.Println(hea)
	return nil
}

func pathJoin(base, p string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}
