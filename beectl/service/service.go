package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/marvin9002/swarm-install/beectl/service/config"
)

type Service struct {
	beeCtrlcfgPath string
	Logger         *log.Logger
	InstallPath    string
	ReportIntval   time.Duration
	ReportServer   string
	CashIntval     time.Duration
}

func NewService(reportAddr, cfgPath string, ReportIntval, CashIntval time.Duration) *Service {
	getwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return &Service{
		beeCtrlcfgPath: cfgPath,
		Logger:         log.New(os.Stdout, "[Agent]: ", log.Ldate|log.Ltime),
		InstallPath:    getwd,
		ReportIntval:   ReportIntval,
		ReportServer:   reportAddr,
		CashIntval:     CashIntval,
	}
}

func (s *Service) Run() {
	var cfg []config.Config
	f, err := os.Open(s.beeCtrlcfgPath)
	if err != nil {
		s.Println("create file: ", err)
		return
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &cfg)
		if err != nil {
			return
		}
	}
	s.Println("Decoder success")
	s.Println(cfg)
	if len(cfg) > 0 {
		for _, conf := range cfg {
			go s.Report(conf)
			go s.CashOut(conf)
		}
	}
	select {}
}

func (s *Service) Printf(format string, v ...interface{}) {
	f, err := os.OpenFile(path.Join(s.InstallPath+"/debug.log"), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	defer f.Close()
	lg := log.New(f, "[Agent]: ", log.Ldate|log.Ltime)
	lg.Printf(format, v)
	s.Logger.Printf(format, v)
}

func (s *Service) Println(v ...interface{}) {
	f, err := os.OpenFile(path.Join(s.InstallPath+"/debug.log"), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	defer f.Close()
	lg := log.New(f, "[Agent]: ", log.Ldate|log.Ltime)
	lg.Println(v)
	s.Logger.Println(v)
}
