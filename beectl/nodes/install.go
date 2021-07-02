package nodes

import (
	"errors"
	"os/exec"
	"runtime"

	"github.com/marvin9002/swarm-install/beectl/utils"
)

func (h *Swarm) Install() error {
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
		return h.Linux()
	}
	return nil
}

func (h *Swarm) Linux() error {
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
