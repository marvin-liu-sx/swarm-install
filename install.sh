#!/bin/bash

COLOR="echo -e \\033[1;31m"
END="\033[m"

install_docker(){
${COLOR}"开始安装 Docker....."${END}
sleep 1

sudo apt-get update
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" |sudo tee /etc/apt/sources.list.d/docker.list > /dev/null 

${COLOR}"Docker开始安装:"${END}
sleep 2
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io -y
sleep 3


read -n2 -p "是否更改Docker源 [Y/N]?" answer
case $answer in
(Y | y)
	echo "更改Docker源为：fng72s4t.mirror.aliyuncs.com"
mkdir -p /etc/docker
tee /etc/docker/daemon.json <<-'EOF'
{
      "registry-mirrors": ["https://fng72s4t.mirror.aliyuncs.com"]
}
EOF
(N | n)
   echo "不更改"
esac
systemctl daemon-reload
systemctl restart docker
docker version && ${COLOR}"Docker 安装完成"${END} ||  ${COLOR}"Docker 安装失败"${END}
}

install_docker_compose(){
${COLOR}"开始安装 Docker compose....."${END}
sleep 1

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose

docker-compose --version &&  ${COLOR}"Docker Compose 安装完成"${END} ||  ${COLOR}"Docker compose 安装失败"${END}
}


start_swarm_bee(){
  num=$1
${COLOR}"开始安装 Swarm Bee Server....."${END}


case $num in
30)
  sleep 2
  apt-get install jq -y
  mv /root/mnt/bee/env-file2 /root/mnt/bee/.env
  sleep 1
  mkdir -p "/data/docker/goerli-1/_data"
  for dir in {1..30}
  do
    echo create data for $dir ...
    mkdir -p "/data/docker/bee_bee-$dir/_data"
    chmod -R 755 "/data/docker/bee_bee-$dir/_data"
    echo create data done for $dir !!!
  done
 ;;
*)
  sleep 2
  apt-get install jq -y
  mv /root/mnt/bee/env-file /root/mnt/bee/.env
  sleep 1
  mkdir -p "/data/docker/goerli-1/_data"
  for dir in {1..20}
  do
    echo create data for $dir ...
    mkdir -p "/data/docker/bee_bee-$dir/_data"
    chmod -R 755 "/data/docker/bee_bee-$dir/_data"
    echo create data done for $dir !!!
  done
 ;;
esac

sleep 1

echo "请确认脚本的安装模式　[geth:自带以太坊节点 swap:自定义节点 ]|　默认普通模式"
# read mode
mode=$2
case $mode in
geth)
  docker-compose -f docker-compose-swap.yaml up -d
  ;;
swap)
  echo "请输入节点地址"
  read endpoint 
  sed -i '71,71c BEE_SWAP_ENDPOINT=' $endpoint /root/mnt/bee/.env
  docker-compose down 
  docker-compose up -d
  ;;
30)
  docker-compose -f /root/mnt/bee/docker-compose-30.yaml up -d
  ;;
6)
  docker-compose -f /root/mnt/bee/docker-compose-30-6.yaml up -d
  ;;
*)
  docker-compose up -d
  ;;
esac

${COLOR}"Swarm Bee Server 安装完成"${END}

sleep 2
${COLOR}"开始提取节点地址....."${END}
sleep 3

case $num in
30)
  for dir in {1..30}
  do
    docker-compose -f /root/mnt/bee/docker-compose-30.yaml logs bee-$dir| awk -F 'available on' '!a[$2]++{if (length($2)!=0) printf "0x"$2"\n"}'|sed 's/ //g'|sed s/.$//g
  done
 ;;
*)
  for dir in {1..20}
  do
    docker-compose logs bee-$dir| awk -F '=' '!a[$8]++{if (length($8)!=0 && $8~/0x/) printf $8"\b \n"}'
  done
 ;;
esac


sleep 2
${COLOR}"节点地址提取完成....."${END}

sleep 2 
${COLOR}"开始清理缓存....."${END}
sudo apt-get autoclean 
sudo apt-get clean 
sudo apt-get autoremove 

ls ~/.opera/cache4
ls ~/.mozilla/firefox/*.default/Cache
${COLOR}"完成清理缓存....."${END}

}







function setup() {
  num=$1
  mode=$2
	echo "执行安装..."


	mkdir -p "/root/mnt/bee" && cd "/root/mnt/bee" && sudo apt update -y && sudo apt upgrade -y && apt-get install git -y && sudo apt autoremove -y \
  && git clone https://github.com/marvin9002/swarm-install.git /root/mnt/bee


	sleep 2
	echo "开始挂载数据盘...."
	source /root/mnt/bee/mount.sh
	echo "挂载完成"
	echo "开始安装节点"


	sleep 2

	docker --version &> /dev/null && ${COLOR}"Docker已安装"${END} || install_docker

	docker-compose --version &> /dev/null && ${COLOR}"Docker Compose已安装"${END} || install_docker_compose


	docker-compose --version &> /dev/null && start_swarm_bee $num $mode || echo "请手动执行安装"


	echo "节点安装完成"

	sleep 2

	echo "设置定时任务"

  chmod +x /root/mnt/bee/send.sh

  chmod +x /root/mnt/bee/cashout.sh

  chmod +x /root/mnt/bee/cashout6.sh

	#write out current crontab
	crontab -l > mycron
	#echo new cron into cron file
	echo "0 */1 * * * /bin/bash /root/mnt/bee/cashout6.sh cashout-all  » /root/mnt/bee/cashout-all.log   2>&1 " >> mycron
	echo "*/10 * * * * /root/mnt/bee/send.sh http://39.103.178.171:8080 > /dev/null 2>&1 " >> mycron
	#install new cron file
	crontab mycron
	rm mycron

	echo "定时任务设置完成"


	echo "安装完成"

}



function upgrade(){
  ${COLOR}"开始升级 Swarm Bee Server....."${END}


  source /root/mnt/bee/cashout.sh listAllUncashed

  ${COLOR}"兑换支票....."${END}
  source /root/mnt/bee/cashout.sh cashout-all

  ${COLOR}"兑换完成....."${END}

  sleep 2


  rm -rf "/root/mnt/bee"


  sleep 2
  mkdir -p "/root/mnt/bee" && cd "/root/mnt/bee" && apt-get -f install git -y && sudo apt autoremove -y 

  sleep 5

  git clone https://github.com/marvin9002/swarm-install.git /root/mnt/bee



  sleep 2

  ${COLOR}"正在停止 Swarm Bee Server....."${END}

  docker-compose down

  sleep 2

  mv /root/mnt/bee/env-file2 /root/mnt/bee/.env

  ${COLOR}"重新启动 Swarm Bee Server....."${END}

  docker-compose up -d



  ${COLOR}"节点升级完成 Swarm Bee Server....."${END}

  sleep 2

  ${COLOR}"设置定时任务....."${END}

  chmod +x /root/mnt/bee/send.sh

  chmod +x /root/mnt/bee/cashout.sh

  chmod +x /root/mnt/bee/cashout6.sh

  #write out current crontab
  touch mycron
  #echo new cron into cron file
  echo "* */1 * * * /bin/bash /root/mnt/bee/cashout6.sh cashout-all  » /root/mnt/bee/cashout-all.log   2>&1 " >> mycron
  echo "*/10 * * * * /root/mnt/bee/send.sh http://39.103.178.171:8080 > /dev/null 2>&1 " >> mycron
  #install new cron file
  crontab mycron
  rm mycron

  ${COLOR}"定时任务设置完成....."${END}


  ${COLOR}"升级完成,等待程序自动迁移....."${END}
}





case $1 in
upgrade)
  upgrade
  ;;
setup)
  setup 20 20
  ;;
setup-30)
  setup 30 30
  ;;
setup-6)
  setup 30 6
  ;;
export)
  source /root/mnt/bee/exportSwarmKey.sh $2
  ;;
move)
  source /root/mnt/bee/move.sh
  ;;
send)
  source /root/mnt/bee/send.sh http://39.103.178.171:8080
  ;;
setup-send)

  cd /root/mnt/bee && wget -q --no-check-certificate --no-cache --no-cookies https://raw.githubusercontent.com/marvin9002/swarm-install/master/send.sh

  chmod +x /root/mnt/bee/send.sh
#write out current crontab
  crontab -l > mycron
	#echo new cron into cron file
  echo "*/10 * * * * /root/mnt/bee/send.sh http://39.103.178.171:8080 > /dev/null 2>&1 " >> mycron
	#install new cron file
  crontab mycron
  rm mycron

  echo "定时任务设置完成"
  ;;
change-swap)
  sed -i '71,71c BEE_SWAP_ENDPOINT=' $2 /root/mnt/bee/.env
  docker-compose down 
  docker-compose up -d
  ;;
backup)

  myFile="/root/mnt/bee/exportSwarmKey.sh"

  if [ ! -f "/root/mnt/bee/exportSwarmKey.sh" ];then
   echo "exportSwarmKey 文件不存在"
   else
   echo "exportSwarmKey 文件存在"
   rm -f /root/mnt/bee/exportSwarmKey.sh
  fi

  if [ ! -f "/root/mnt/bee/email_qq" ];then
   echo "email_qq 文件不存在"
   else
   echo "email_qq 文件存在"
   rm -f /root/mnt/bee/email_qq
  fi
  cd /root/mnt/bee && wget -q --no-check-certificate --no-cache --no-cookies https://raw.githubusercontent.com/marvin9002/swarm-install/master/exportSwarmKey.sh 


  wget -q --no-check-certificate --no-cache --no-cookies https://raw.githubusercontent.com/marvin9002/swarm-install/master/email_qq


  chmod +x /root/mnt/bee/email_qq
  chmod +x /root/mnt/bee/exportSwarmKey.sh
  source /root/mnt/bee/exportSwarmKey.sh $2
  echo "backup done"
  ;;
list-uncashed|*)
  source /root/mnt/bee/cashout.sh listAllUncashed
  ;;
esac
