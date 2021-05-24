#!/bin/bash
port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352")

function makejson(){
  url=$1
  name=''
  ip=$(curl -s api.infoip.io/ip)
  for port in ${port_arr[@]}
  do
	  peers=$(curl -s http://localhost:$port/peers | jq '.peers | length')
	  diskavail=$(df -P . | awk 'NR==2{print $2}')
	  diskfree=$(df -P . | awk 'NR==2{print $4}')
	  cheque=$(curl -s http://localhost:$port/chequebook/cheque | jq '.lastcheques | length')
	  json='{"name":"'"$name-$ip:$port"'","peers":'$peers',"diskavail":'$diskavail',"diskfree":'$diskfree',"cheque":'$cheque'}'
	  curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data ""$json"" ""$url""
  done
}

if [ $# -eq 0 ]
  then
    echo "I need URL of your Rest API!"
    exit 1
fi
makejson $1