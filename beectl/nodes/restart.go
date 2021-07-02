package nodes

import (
	"errors"
	"fmt"

	"github.com/marvin9002/swarm-install/beectl/utils"
)

func (h *Swarm) Restart(nodeId string) error {
	if nodeId == "" {
		var cfg []config
		cfg, err := h.GetCfg()
		if err != nil {
			return err
		}
		for _, conf := range cfg {
			err := h.restart(conf.ID)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		err := h.restart(nodeId)
		if err != nil {
			return err
		}
	}
	return nil

}

func (h *Swarm) restart(nodeId string) error {

	if !utils.ExecCommand("systemctl", []string{"restart", fmt.Sprintf("bee%s", nodeId)}) {
		h.Println("restart bee error")
		return errors.New("systemctl daemon-reload error: exec: \"systemctl\": executable file not found in $PATH")
	}
	h.Println(fmt.Sprintf("restart bee%s success ", nodeId))
	return nil

}
