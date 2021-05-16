#!/bin/bash

COLOR="echo -e \\033[1;31m"
END="\033[m"

install_docker(){
${COLOR}"开始安装 Docker....."${END}
sleep 1

sudo yum install -y yum-utils
sudo yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo

${COLOR}"Docker开始安装:"${END}
sleep 2
sudo yum install docker-ce docker-ce-cli containerd.io
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
sudo systemctl start docker
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

${COLOR}"开始安装 Swarm Bee Server....."${END}
sleep 2


docker-compose up -d
${COLOR}"Swarm Bee Server 安装完成"${END}
sleep 2
${COLOR}"开始提取节点地址....."${END}
sleep 3
docker-compose logs -f bee-1 bee-2 bee-3 bee-4 bee-5 bee-6 bee-7 bee-8 bee-9 bee-10 bee-11 bee-12 bee-13 bee-14 bee-15 bee-16 bee-17 bee-18 bee-19 bee-20| awk -F '=' '!a[$8]++{if (length($8)!=0 && $8~/0x/) printf $8"\b \n"}'
${COLOR}"节点地址提取完成....."${END}
}

docker --version &> /dev/null && ${COLOR}"Docker已安装"${END} || install_docker

docker-compose --version &> /dev/null && ${COLOR}"Docker Compose已安装"${END} || install_docker_compose


docker-compose --version &> /dev/null && start_swarm_bee || echo "请手动执行安装"
