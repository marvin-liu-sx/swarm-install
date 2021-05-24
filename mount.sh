#!/usr/bin/env bash

# parted脚本自动挂载分区磁盘
#  1.parted 核心命令
# yum install -y parted                                   # 安装 parted 分区工具包
# parted -s /dev/sdb mklabel msdos                # parted -s 选择磁盘 /dev/sdb
#                                                                         # 格式化磁盘 /dev/sdb为gpt 动态分区
#                                                                         # label  [ˈlebəl] 标签
# parted -s /dev/sdb mkpart primary 0 100%        # 新建主分区，全部空间
#                                                                         # part [pɑrt] 分开，分区
#                                                                         # primary [ˈpraɪˌmɛri] 主分区
# # parted -s a mkpart entended 3G 5G                     # 第一个扩展分区:从3G 到5G
# # parted -s b mkpart logic 5G 100%                              # 第二个扩展分区:从5G到100%
#                                                                                         # logic [ˈlɒdʒɪk] 逻辑，分区
# parted -s /dev/sdb print                                        # 选择磁盘，并打印信息；可以看到 Number [ˈnʌmbər] 编号 1
# mkfs -t ext4 /dev/sdb1                                  # 格式化分区为 ext4  ；也可以格式化为ext3 等其他分区
# mkdir /www                                                      # 新建挂载分区 /www
# mount /dev/sdb1 /www                                    # 挂载分区到到文件夹 /www
# echo "
# /dev/sdb1                               /www                    ext4    defaults        0 0
# " >> /etc/fatab                                         # 写入开机启动配置文件
#                                                                         #  default [diˈfɔ:lts] 默认
# df -h                                                   # 重启服务器并查看挂载的分区
# 2. parted 脚本自动分区

azparted=$(apt list --installed | grep parted)
# 定义一个名称为azparted的变量，值为：
# 查看已安装的包，grep 匹配parted 名称
cdazparted=$(echo ${#azparted})
# 定义一个新变量：打印$azparted 变量的字符串长度  echo $(#azparted)
 
if [ $cdazparted -lt 1 ]
# 判断变量字符串长度小于1
then                                            # 然后
	echo "您还没有安装parted，正在为您安装，请稍后："
    sudo apt-get install  parted
else                                            # 其他情况
    echo "parted 已经安装，无需操作,已经安装的信息为：
    $azparted"
fi                                              # if的结束标记 fi
#----------------------------------- 检测安装工具向上包结束  -------------------------------
echo "未分区的磁盘有: "
parted -l | grep -w "unrecognised disk label"
echo "read 接收窗口命令界面输入的字符串;-p 加文字说明；需要分区的磁盘变量 a请输入你要进行分区的磁盘; 磁盘格式为：/dev/dev 如果输入错误字符或者闪跳，请Ctrl +c  退出重新输入"
read -p "请输入要分区的磁盘："  a               ; 
echo "磁盘挂载目录格式为：   /www     将会清空该目录下的文件，重复执行脚本，仅仅更改挂载目录，不会更改硬盘数据"
read -p "请输入您要挂载到那个目录："  m         ; 
#----------------------------------- 用户输入信息向上结束  -------------------------------
gsh=$(parted -s $a print | grep primary  )

# 定义一个变量 gsh  ； sed s/[[:space:]]//g 删除空格
# parted -s $a print     打印要分区的磁盘信息$a为 /dev/xxx
# 打印信息后去匹配磁盘分区的类型：
# primary [ˈprʌɪm(ə)ri] 主分区，类型；logical  [ˈlɒdʒɪk(ə)l] 逻辑分区
#　如果主分区和逻辑分区都不存在，那么就是空盘
###### 重点：
# parted -s $a print                            # 选择磁盘，并打印：打印选择的磁盘信息
#  awk -F " " '$5=="primary" {print zi++ $6} '
# -F " " 分隔符为空格,  后面有单引号引起来的为坐标：先行后列
# $5 第5列的内容包含有 primary 字符  就被选中
# print 打印第6列，  有多个内容用zi++ 方法来换行打印所有
# parted -s /dev/sda print | awk -F " " '$5=="primary" {print c++ $6} '
if [ ${#gsh}  -eq 0 ]                                   # ${#gsh} 变量长度， -eq 小等于 0
then
    echo "磁盘 $a 为空盘，没有找到主分区和逻辑分区，可以进行进行格式化并分区操作，请等待。。。  "
else 
   echo "  #　parted -s /dev/xxx   rm 1 磁盘名称： $a 磁盘存在数据，是否继续格式化操作"
fi
echo "请确认脚本的继续执行ｙ继续　|　任意字符退出"
read gshy
case $gshy in
y) echo "您确认了继续格式化操作,脚本将继续执行"
;;
*) echo "您否定了格式化操作，不会影响磁盘数据。即将退出脚本。"
exit                                                                    # 直接退出脚本
;;
esac
sleep 1
parted -s $a mklabel gpt                      # 格式化为gpt 动态分区
sleep 1
parted -s $a unit s
 # msdos 其他类型
sleep 1
parted -s $a mkpart data 2048s 100%      # 分区 全部
#parted -s $a mkpart entended 3G 5G     # 第一个扩展分区:从3G 到5G
#parted -s $a mkpart logic 5G 100%      # 第二个扩展分区:从5G到100%
# logic [ˈlɒdʒɪk] 逻辑，分区
#-------------------------------------挂载------------------------------------
sleep 1
b=$(echo $a"1")
sleep 1
mkfs -t ext4 $b                                 # 格式化分区
rm -rf $m
mkdir $m
mount  $b $m                                            # 挂载分区到/www 目录下
# 如果需要挂载到指定目录请创建文件夹后，再将此处的/mnt修改
bd=$(blkid $b | awk -F '=' '{print substr($2, 2, length($2)-7)}')
# 打印变量，awk 字段处理，-F 指定分隔符为  / 
# 坐标：第一行，第3列
sed -i "/"eiscparted"/d" /etc/fstab             #先清除启动挂载
sed -i "/^$/d" /etc/fstab                # 清除空行
echo "UUID=$bd           $m                            ext4      defaults        0 2                #eiscparted">>/etc/fstab     
sudo mount -a                              # 开机自动挂载，字符单独一行，才会保留格式
echo "再次查看挂载的磁盘，如果之前有挂载过此硬盘，重启生效挂载到新目录
"
df -h
# umount /dev/sdb*                      # 取消挂载所有分区
# -------------------------------------删除---------------------------------
# parted -s /dev/sdb rm 5                       # rm删除sdb磁盘编号5的分区 
# parted -s /dev/sdb print                      # 查看分区
# parted -s /dev/sdb mklabel msdos      # 清除分区表，方便其他工具进行分区
# 脚本执行：
# rm -rf parted.sh ; yum install -y wget ; wget eisc.cn/file/shell/parted.sh ; chmod +x parted.sh ; ./parted.sh