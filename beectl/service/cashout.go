package service

import (
	"fmt"
	"time"

	"github.com/marvin9002/swarm-install/beectl/service/config"

	"github.com/marvin9002/swarm-install/beectl/service/nodes"
)

func (s *Service) CashOut(cfg config.Config) {
	s.Println("Start CashOut...")
	node := nodes.NewNode(cfg)
	t := time.NewTicker(s.CashIntval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			cashout, err := node.Cashout()
			s.Println(fmt.Sprintf("CashOut-Data:%#v", cashout))
			if err != nil {
				s.Println(err)
				continue
			}
			s.Println("CashOut-Hash:", cashout)
		}
	}

}
