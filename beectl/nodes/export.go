package nodes

import (
	"fmt"
	"os"
	"path"

	"github.com/marvin9002/swarm-install/beectl/utils"
)

func (h *Swarm) ExportSwarmKey(addr string) error {
	h.Println(fmt.Sprintf("可通过命令查看此bee节点的日志(非常重要): journalctl -f -u bee%s \n", addr))

	dest, err := h.backUp(addr)
	if err != nil {
		return err
	}

	h.Printf("已经打包备份钱包相关信息在:%s \n", dest)
	h.Println("!!! 请注意备份钱包相关信息目录文件, 请将备份后的文件保存到其他电脑上以防止本机硬盘损坏以后钱包丢失!!!")
	h.Println("当前输出信息已经存入当前目录的 result.log 中，方便您后期查看")
	return nil
}

func (h *Swarm) ExportAllSwarmKey() error {
	var cfg []config
	cfg, err := h.GetCfg()
	if err != nil {
		return err
	}
	if len(cfg) <= 0 {

	}
	for _, cf := range cfg {
		dest, err := h.backUp(cf.ID)
		if err != nil {
			return err
		}
		h.Printf("已经打包备份bee%s节点钱包相关信息在:%s \n", cf.ID, dest)
	}

	h.Println("!!! 请注意备份钱包相关信息目录文件, 请将备份后的文件保存到其他电脑上以防止本机硬盘损坏以后钱包丢失!!!")
	h.Println("当前输出信息已经存入当前目录的 result.log 中，方便您后期查看")
	return nil
}

func (h *Swarm) backUp(nodeId string) (string, error) {
	f1, err := os.Open(path.Join(h.Dir, fmt.Sprintf("/bee%s", nodeId)))
	if err != nil {
		h.Println(err)
		return "", err
	}
	defer f1.Close()
	f2, err := os.Open(path.Join(h.BeeCfgPath, fmt.Sprintf("/bee%s.yaml", nodeId)))
	if err != nil {
		h.Println(err)
		return "", err
	}
	defer f2.Close()

	var files = []*os.File{f1, f2}
	dest := path.Join(path.Dir(h.InstallPath), fmt.Sprintf("/bkup_bee%s_keys.zip", nodeId))
	err = utils.Compress(files, dest)
	if err != nil {
		h.Println(err)
		return "", err
	}
	return dest, nil
}
