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


	for dir in {1..20}
	do
		mkdir -p "/data/bee-keys/bee_bee-$dir"
		mkdir -p "/data/bee-bee_clef/bee_clef-$dir"
		cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/keystore" "/data/bee-bee_clef/bee_clef-${dir}"
		cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/password" "/data/bee-bee_clef/bee_clef-${dir}"
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
	tar -zcvf clef.tar.gz /data/bee-bee_clef
	tar -zcvf keys.tar.gz /data/bee-keys
}
Install_Main