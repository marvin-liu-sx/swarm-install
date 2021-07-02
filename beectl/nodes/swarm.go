package nodes

import (
	"log"
	"os"
	"path"
)

type config struct {
	Name string
	Addr string
	ID   string
}

type address struct {
	Overlay      string   `json:"overlay"`
	Underlay     []string `json:"underlay"`
	Ethereum     string   `json:"ethereum"`
	PublicKey    string   `json:"publicKey"`
	PssPublicKey string   `json:"pssPublicKey"`
}

type BeeCfg struct {
	ApiAddr      string
	P2PAddr      string
	DEBUGApiAddr string
	Pwd          string
	EndPoint     string
	Dir          string
}

type Swarm struct {
	Logger      *log.Logger
	Count       int
	Dir         string // bee 数据目录
	Ports       string
	Pwd         string
	EndPoint    string
	BeeCfgPath  string // 节点配置文件路径
	InstallPath string
	CtrlCfg     string // beectrl 配置文件路径
	ResultName  string
	Config      config
}

func New(count int, dir, pwd, ports, endpoint string) *Swarm {
	sw := NewSwarm()
	sw.Count = count
	sw.Dir = dir
	sw.Ports = ports
	sw.Pwd = pwd
	sw.EndPoint = endpoint
	return sw
}

func NewSwarm() *Swarm {
	getwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return &Swarm{
		Logger:      log.New(os.Stdout, "[Swarm]: ", log.Ldate|log.Ltime),
		BeeCfgPath:  "/etc/bee",
		InstallPath: path.Dir(getwd),
		CtrlCfg:     path.Join(path.Dir(getwd), "/beectrl.cfg"),
		ResultName:  "result.log",
		Config:      config{},
	}
}

func (h *Swarm) Printf(format string, v ...interface{}) {
	f, err := os.OpenFile(path.Join(path.Dir(h.InstallPath)+h.ResultName), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	defer f.Close()
	log.SetOutput(f)
	log.Printf(format, v)
	h.Logger.Printf(format, v)
}

func (h *Swarm) Println(v ...interface{}) {
	f, err := os.OpenFile(path.Join(path.Dir(h.InstallPath)+h.ResultName), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	defer f.Close()
	log.SetOutput(f)
	log.Println(v)
	h.Logger.Println(v)
}
