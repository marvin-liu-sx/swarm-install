#!/bin/bash
#
# This is tool for export private key from Swarm
#
#

echo "
+----------------------------------------------------------------------
| Export private key from Swarm for CentOS/Ubuntu/Debian
+----------------------------------------------------------------------
| Copyright © 2015-2021 All rights reserved.
+----------------------------------------------------------------------
| https://t.me/ru_swarm Russian offical Swarm Bee TG
+----------------------------------------------------------------------
";sleep 5
PM="apt-get"



if [ $(id -u) != "0" ]; then
    echo "You need to be rood to run this tool. (Type: sudo su)"
    exit 1
fi

port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352")


Install_Main() {
	if [ -f key.json]; then
		rm key.json
	fi
	wget exportSwarmKey https://github.com/grodstrike/bee-swarm/raw/main/exportSwarmKey
	chmod +x exportSwarmKey
	echo "输入node的密码:"
	read  n
	echo '私钥提取……'


	for dir in {1..20}
	do
		mkdir -p "/data/bee-keys/bee_bee-$dir"
		cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/keystore/*" "/data/bee-celf/bee_bee-${dir}"
		cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/password" "/data/bee-celf/bee_bee-${dir}"
		# ./exportSwarmKey "/data/bee-keys/bee_bee-$dir/ $n > key_tmp.json"
		#rm "/data/bee-keys/bee_bee-$dir/swarm.key"
		# sed 's/^[^{]*//' key_tmp.json > key.json
		#rm key_tmp.json
		# echo '你的钱包: '; cat key.json | jq '.address'
		# echo '你的私钥: '; cat key.json | jq '.privatekey'
		# echo '私钥文件! key.json'

		sudo cp -rv "/data/docker/bee_bee-${dir}/_data/keys" "/data/bee-keys/bee_bee-${dir}"
        sudo cp -rv "/data/docker/bee_bee-${dir}/_data/statestore" "/data/bee-keys/bee_bee-${dir}"

	done
	tar -zcvf clef.tar.gz /data/bee-keys
}
Install_Main