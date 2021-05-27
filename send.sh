#!/bin/bash
port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352")

function makejson(){
  url=$1
  ip=$(curl -s api.infoip.io/ip)
  cpu_us=`vmstat | awk '{print $13}' | sed -n '$p'`
  cpu_sy=`vmstat | awk '{print $14}' | sed -n '$p'`
  cpu_id=`vmstat | awk '{print $15}' | sed -n '$p'`
  cpu_sum=$(($cpu_us+$cpu_sy))
  diskavail=$(df -P | grep '/dev/vdb1' | awk {'print $3'})
  diskfree=$(df -P . | awk 'NR==2{print $4}')

  ramusage=$(free | awk '/Mem/{printf("RAM Usage: %.2f\n"), $3/$2*100}'| awk '{print $3}')

  for port in ${port_arr[@]}
  do
	  peers=$(curl -s http://localhost:$port/peers | jq '.peers | length')
	  
	  cheque=$(curl -s http://localhost:$port/chequebook/cheque | jq '.lastcheques | length')
	  json='{"name":"'"$ip:$port"'","peers":'$peers',"diskavail":'$diskavail',"diskfree":'$diskfree',"cheque":'$cheque',"cpu_sum":'$cpu_sum',"memory_usage":'$ramusage'}'
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