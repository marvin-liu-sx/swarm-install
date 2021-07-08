package main

import (
	"flag"
	"time"

	"github.com/marvin9002/swarm-install/beectl/service"
)

func main() {
	repotAddr := flag.String("repot-addr", "http://localhost", "监控上报地址")
	cfgPath := flag.String("ctrl-cfg", "", "Bee Nodes 管理工具 配置文件路径")

	ReportIntval := flag.String("report-intval", "1s", "监控上报间隔时间")
	CashIntval := flag.String("cash-intval", "1m", "支票兑换间隔时间")
	flag.Parse()

	CashIntvalTime, err := time.ParseDuration(*CashIntval)
	if err != nil {
		return
	}
	ReportIntvalTime, err := time.ParseDuration(*ReportIntval)
	if err != nil {
		return
	}

	srv := service.NewService(*repotAddr, *cfgPath, ReportIntvalTime, CashIntvalTime)
	srv.Run()
}
