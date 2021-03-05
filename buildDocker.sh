#!/bin/sh
#
#
# File must be created under Linux
#
#
# 默认本地打包发布
# ./buildDocker.sh 发布线上环境81
# ./buildDocker.sh dev 发布线上环境106.52.210.81
# ./buildDocker.sh master 发布线上环境106.52.210.43

set -e

#250内网地址
REMOTE_HOST="192.168.1.250"
#REMOTE_HOST="172.25.110.172"
DEV_REMOTE_HOST="106.52.210.81"
MASTER_REMOTE_HOST="106.52.210.43"

REMOTE_USERNAME="root"
START_PORT="9000"
TAG_VERSION="v1.0"
IDRSA_PATH="./tools/ssh/id_rsa"
IMAGE_REPOSITORY="oktools"

#
##HOST_IP=$(ifconfig eth0 | grep "inet" | awk '{ print $2}' |  awk -F: 'NR==1{print $1}')
##echo $HOST_IP
#
#if [[ $1 == "master" ]]; then
#
#  TAG_VERSION="${TAG_VERSION}-master"
#  DEV_REMOTE_HOST=$MASTER_REMOTE_HOST
#  echo $DEV_REMOTE_HOST
#  echo "正式环境版本发布"
#
#elif  [[ $1 == "dev" ]]; then
#
#  TAG_VERSION="${TAG_VERSION}-dev"
#  DEV_REMOTE_HOST=$DEV_REMOTE_HOST
#  echo $DEV_REMOTE_HOST
#  echo "dev环境版本发布"
#
#else
#
#  TAG_VERSION="${TAG_VERSION}-local"
#  DEV_REMOTE_HOST=$REMOTE_HOST
#  echo
#  echo "本地环境版本发布"
#
#fi
#echo  $DEV_REMOTE_HOST $TAG_VERSION
#
#
#git status
#git pull origin

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/main.go
docker build -t $IMAGE_REPOSITORY:$TAG_VERSION ./
#rm -f ./main
docker save $IMAGE_REPOSITORY:$TAG_VERSION -o $IMAGE_REPOSITORY:$TAG_VERSION

##docker images|grep none|awk '{print $3 }'|xargs docker rmi
#chmod 0600 $IDRSA_PATH

#检查文件或者文件夹是否存在，如果存在则删除
# $1--需要检查的文件
function checkFileExists(){
    if [ ! -e ${1} ]; then
      rm -f ${1}
      echo checkFileExists
    fi
}

function pushDevDocker(){
  #echo "线上打包发布"
  #checkFileExists $IMAGE_REPOSITORY:$TAG_VERSION
  # 导出镜像
  docker save $IMAGE_REPOSITORY:$TAG_VERSION -o $IMAGE_REPOSITORY:$TAG_VERSION
  # scp免密 上传镜像
  scp -i $IDRSA_PATH ./$IMAGE_REPOSITORY:$TAG_VERSION $REMOTE_USERNAME@$DEV_REMOTE_HOST:~
  # 上传docker-compose.yaml
#  scp -i $IDRSA_PATH ./docker-compose.yaml $REMOTE_USERNAME@$DEV_REMOTE_HOST:/home/finebaas/
  rm -f ./$IMAGE_REPOSITORY:$TAG_VERSION
  # ssh免密
  ssh -i $IDRSA_PATH -tt $REMOTE_USERNAME@$DEV_REMOTE_HOST << eeooff
  docker stop $IMAGE_REPOSITORY
  docker rm -f ${IMAGE_REPOSITORY} || true
  docker load <  $IMAGE_REPOSITORY:$TAG_VERSION
  rm -f ./$IMAGE_REPOSITORY:$TAG_VERSION
  cd /home/docker
  docker-compose down
  docker-compose up -d
  sleep 1.5s
  docker exec -it $IMAGE_REPOSITORY /bin/sh
  chmod +x /app/tools/tools_aes_decrypt
  chmod +x /app/tools/tools_aes_encrypt
  exit
eeooff
  echo done!
}

function localBuildDocker() {
  echo "本地打包发布"
  docker stop $IMAGE_REPOSITORY
  docker-compose down
  docker-compose up -d
  sleep 2s
  # ssh免密
  ssh -i $IDRSA_PATH -tt $REMOTE_USERNAME@$REMOTE_HOST << eeooff
  docker exec -it $IMAGE_REPOSITORY /bin/sh
  chmod +x /app/tools/tools_aes_decrypt
  chmod +x /app/tools/tools_aes_encrypt
  exit
eeooff
  echo done!
}


#pushDevDocker
#if [[ $1 == "dev" || $1 == "master" ]]; then
#  pushDevDocker
#else
#  localBuildDocker
#fi