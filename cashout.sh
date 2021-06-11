#!/usr/bin/env bash
[ -z ${DEBUG_API+x} ] && DEBUG_API=http://localhost:
[ -z ${MIN_AMOUNT+x} ] && MIN_AMOUNT=100

port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352" "18353" "18354" "18355" "18356" "18357" "18358" "18359" "19351" "19352" "19353")

function getPeers() {
  local port=$1
  curl -s "$DEBUG_API$port/chequebook/cheque" | jq -r '.lastcheques | .[].peer'
}

function getCumulativePayout() {
  local peer=$1
  local port=$2
  local cumulativePayout=$(curl -s "$DEBUG_API$port/chequebook/cheque/$peer" | jq '.lastreceived.payout')
  if [ $cumulativePayout == null ]
  then
    echo 0
  else
    echo $cumulativePayout
  fi
}

function getLastCashedPayout() {
  local peer=$1
  local port=$2
  local cashout=$(curl -s "$DEBUG_API$port/chequebook/cashout/$peer" | jq '.cumulativePayout')
  if [ $cashout == null ]
  then
    echo 0
  else
    echo $cashout
  fi
}

function getUncashedAmount() {
  local peer=$1
  local port=$2
  local cumulativePayout=$(getCumulativePayout $peer $port)
  if [ $cumulativePayout == 0 ]
  then
    echo 0
    return
  fi

  cashedPayout=$(getLastCashedPayout $peer $port)
  let uncashedAmount=$cumulativePayout-$cashedPayout
  echo $uncashedAmount
}

function cashout() {
  local peer=$1
  local port=$2
  local response=$(curl -s -XPOST "$DEBUG_API$port/chequebook/cashout/$peer")
  local txHash=$(echo "$response" | jq -r .transactionHash)
  local addr=$(curl -s $DEBUG_API$port/addresses | jq '.ethereum')
  local peer_num=$(curl -s $DEBUG_API$port/peers | jq '.peers | length')
  echo "节点: $port 节点地址: $addr  节点链接数: $peer_num | result: $result | txHash: $txHash\r\n" > /root/mnt/bee/success_$port.log 2>&1
}


function cashoutAll() {
  local minAmount=$1
  
  for port in ${port_arr[@]}
  do
    for peer in $(getPeers $port)
    do
      local uncashedAmount=$(getUncashedAmount $peer $port)
      if (( "$uncashedAmount" > $minAmount ))
      then
        echo "uncashed cheque for $peer ($uncashedAmount uncashed)" >&2
        cashout $peer $port
      fi
    done
  done
}

function listAllUncashed() {
  counts=0
  no_cash=0
  for port in ${port_arr[@]}
  do
    echo "list $port"
    for peer in $(getPeers $port)
    do
      echo $peer
      let counts+=1
      local uncashedAmount=$(getUncashedAmount $peer $port)
      if (( "$uncashedAmount" > 0 ))
      then
        let no_cash+=1
        echo $peer $uncashedAmount
      fi
    done
  done
  echo chequebook total: $counts
  echo chequebook no cash total: $no_cash
}





case $1 in
cashout)
  cashout $2
  ;;
cashout-all)
  cashoutAll $MIN_AMOUNT
  ;;
list-uncashed|*)
  listAllUncashed
  ;;
esac















