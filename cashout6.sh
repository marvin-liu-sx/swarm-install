#!/usr/bin/env bash
[ -z ${DEBUG_API+x} ] && DEBUG_API=http://localhost:  
[ -z ${MIN_AMOUNT+x} ] && MIN_AMOUNT=100
port_arr=("16351" "16352" "16353" "16354" "16355" "16356" "16357" "16358" "16359" "17351" "17352" "17353" "17354" "17355" "17356" "17357" "17358" "17359" "18351" "18352" "18353" "18354" "18355" "18356" "18357" "18358" "18359" "19351" "19352" "19353")

# cashout script for bee >= 0.6.0
# note this is a simple bash script which might not work well or at all on some platforms
# for a more robust interface take a look at https://github.com/ethersphere/swarm-cli

function getPeers() { 
  local port=$1
  curl -s "$DEBUG_API$port/chequebook/cheque" | jq -r '.lastcheques | .[].peer'  
}

function getUncashedAmount() {  
  curl -s "$DEBUG_API$2/chequebook/cashout/$1" | jq '.uncashedAmount'  
}

function cashout() {
  local peer=$1
  txHash=$(curl -s -XPOST "$DEBUG_API$2/chequebook/cashout/$peer" | jq -r .transactionHash)
  echo cashing out cheque for $peer in transaction $txHash >&2
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
      # local uncashedAmount=$(getUncashedAmount $peer)
      # if (( "$uncashedAmount" > 0 ))
      # then
      #   echo $peer $uncashedAmount
      # fi
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
uncashed-for-peer)
  getUncashedAmount $2
  ;;
list-uncashed|*)
  listAllUncashed
  ;;
esac
