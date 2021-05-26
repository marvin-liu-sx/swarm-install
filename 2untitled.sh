https://goerli.infura.io/v3/595a77530f414f16b321e4b898c816d6

http://47.243.21.250:8545/


http://8.209.241.89:8545


docker-compose down




docker-compose -f docker-compose-move.yaml up -d




docker-compose -f docker-compose-move.yaml logs -f bee-1


curl -s http://localhost:16351/peers | jq '.peers | length'


curl localhost:16351/chequebook/balance | jq





select time,metric,node,cheque from (SELECT   timeof AS "time",   host AS metric,node  ,cheque FROM nodes WHERE   timeof BETWEEN FROM_UNIXTIME(1621839039) AND FROM_UNIXTIME(1621924839) group by host,node ORDER BY timeof);



mkdir -p mnt/bee && cd mnt/bee



sed -i '23,23c BEE_CORS_ALLOWED_ORIGINS=*' /root/mnt/bee/.env
cd /root/mnt/bee &&  docker-compose down && docker-compose up -d

cd /root/mnt/bee && docker-compose -f docker-compose-move.yaml up -d




mkdir -p "/root/mnt/bee" && cd "/root/mnt/bee" && git clone https://github.com/marvin9002/swarm-install.git /root/mnt/bee   && /bin/bash install.sh setup 


wget -q --no-check-certificate --no-cache --no-cookies https://github.com/marvin9002/swarm-install/raw/master/install.sh && /bin/bash install.sh setup 