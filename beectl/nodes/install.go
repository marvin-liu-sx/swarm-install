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

		if utils.ExecCommand("cat", []string{"/etc/redhat-release"}) {
			h.Println("InstallBeeLinux Centos......")
			if !utils.Exist("./bee_1.0.0_amd64.rpm") {
				if utils.ExecCommand("wget", []string{"https://github.com/ethersphere/bee/releases/download/v1.0.0/bee_1.0.0_amd64.rpm"}) {
					utils.ExecCommand("rpm", []string{"-i", "bee_1.0.0_amd64.rpm"})
				} else {
					h.Println("InstallBeeLinux Centos  err......")
					return errors.New("Install bee error")
				}
			}
			if !utils.ExecCommand("dpkg", []string{"-i", "bee_1.0.0_amd64.deb"}) {
				h.Println("install error")
			}
		}

		if utils.ExecCommand("cat", []string{"/etc/lsb-release"}) {
			h.Println("InstallBeeLinux Ubuntu......")
			if !utils.Exist("./bee_1.0.0_amd64.deb") {
				if utils.ExecCommand("wget", []string{"https://github.com/ethersphere/bee/releases/download/v1.0.0/bee_1.0.0_amd64.deb"}) {
					utils.ExecCommand("dpkg", []string{"-i", "bee_1.0.0_amd64.deb"})
				} else {
					h.Println("InstallBeeLinux Ubuntu  err......")
					return errors.New("Install bee error")
				}
			}
			if !utils.ExecCommand("dpkg", []string{"-i", "bee_1.0.0_amd64.deb"}) {
				h.Println("install error")
			}
		}

	case "armv7":
	case "arm64":
	}
	h.Println("install bee success")
	return nil
}
