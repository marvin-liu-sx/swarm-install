https://goerli.infura.io/v3/595a77530f414f16b321e4b898c816d6

http://47.243.21.250:8545/


http://8.209.241.89:8545


docker-compose down




docker-compose -f docker-compose-move.yaml up -d




docker-compose -f docker-compose-move.yaml logs -f bee-1


curl -s http://localhost:16351/peers | jq '.peers | length'


curl localhost:16351/chequebook/balance | jq


curl localhost:16351/chequebook/cheque | jq

BEE_BOOTNODE=/ip4/104.131.161.236/tcp/1634/p2p/16Uiu2HAm1eaF1oKBWewu9NHUte4HKRnjp5AFBXKu8H4JpNqcSvCH


  
select time,metric,node,cheque from (SELECT   timeof AS "time",   host AS metric,node  ,cheque FROM nodes WHERE   timeof BETWEEN FROM_UNIXTIME(1621839039) AND FROM_UNIXTIME(1621924839) group by host,node ORDER BY timeof);



mkdir -p mnt/bee && cd mnt/bee



sed -i '23,23c BEE_CORS_ALLOWED_ORIGINS=*' /root/mnt/bee/.env
cd /root/mnt/bee &&  docker-compose down && docker-compose up -d

cd /root/mnt/bee && docker-compose -f docker-compose-move.yaml up -d




mkdir -p "/root/mnt/bee" && cd "/root/mnt/bee" && git clone https://github.com/marvin9002/swarm-install.git /root/mnt/bee   && /bin/bash install.sh setup 

wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh && chmod +x /root/install.sh && /bin/bash install.sh setup-6 

rm /root/mnt/bee/exportSwarmKey && rm /root/mnt/bee/clef.tar.gz && rm /root/mnt/bee/clef.tar.gz &&
rm install.sh && wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh  && chmod +x install.sh && /bin/bash install.sh export 20


rm install.sh && wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh  && chmod +x install.sh && /bin/bash install.sh change-swap ws://47.243.21.250:8546




if [ ! -f "/root/install.sh" ];then
     echo "install 文件不存在"
     else
     echo "install 文件存在"
     rm -f /root/install.sh
    fi
    wget -d --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh  && chmod +x install.sh && /bin/bash install.sh backup 20 





wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh && chmod +x /root/install.sh && /bin/bash install.sh upgrade 



#!/usr/bin/env bash
port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352" "18353" "18354" "18355" "18356" "18357" "18358" "18359" "19351" "19352" "19353")

function test(){
 for port in ${port_arr[@]}
 do
  echo $(curl -s http://localhost:$port/addresses | jq -r '.ethereum')
 done
}

test



cd /root/mnt/bee

    docker-compose down
    sleep 2
    for dir in {1..20}
    do
     rm -rf "/data/docker/bee_bee-$dir/_data/localstore"
    done
    sleep 2
    docker-compose up -d



ln -s /data/bee /root/mnt/










