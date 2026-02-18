#!/bin/bash

env=production
target=crmApp
SERVICE_NAME=target
SERVICE_FILE="$target.service"
TARGET_PATH="/etc/systemd/system/$SERVICE_FILE"
src_path="$(go env GOPATH)/src"
project_path=$(cd $(dirname $0); pwd)
project_path_release=$project_path/app/
release_path=/var/local/rbac/server/
release_project_path=$release_path/$target
release_bin_path=$release_project_path
release_config_path=$release_project_path/config/
release_log_path=$release_project_path/logs
run_app_path="/var/local/rbac/server/wxMiniApp"



printEnv(){
    printf "Print Env \n"
    printf "============================================\n"
    printf "Commend Params        | %s %s \n" $1  $2
    printf "Project Path          | %s\n" $project_path
    printf "Src Path              | %s\n" $src_path
    printf "project_path_release  | %s\n" $project_path_release
    printf "Release Path          | %s\n" $release_path
    printf "Release Bin  Path     | %s\n" $release_bin_path
    printf "Release Config Path   | %s\n" $release_config_path
    printf "============================================\n\n\n"
}

cleanDir(){
    printf "Clean Release Dir \n"
    printf "============================================\n"
    cd $release_path/
    echo $release_path/
    #rm -rf ./*
    if [ $? != 0 ]; then
        printf "Clean release dir failed\n"
        exit 101
    else
        printf "Clean release dir successed\n"
    fi

    mkdir -p $release_config_path
    mkdir -p $release_bin_path
    mkdir -p $release_log_path
    printf "============================================\n\n\n"
}

buildBin(){
    printf "Build Bin \n"
    printf "============================================\n"
    cd $project_path_release/
    printf "Pull dependence  ...\n"
    go mod tidy
    go mod vendor
    if [ $? != 0 ]; then
        printf "Compiling project failed\n"
        exit 100
    fi
    printf "Pull dependence End\n"
    printf "Compiling project ...\n"

    go build -o $release_project_path/$target
    if [ $? != 0 ]; then
        printf "Compiling project failed\n"
        exit 102
    else
	    printf "Compiling project successed\n"
    fi
    cd $release_project_path
    printf "============================================\n\n\n"
}

copyConf(){
    printf "Copy Conf Files\n"
    printf "============================================\n"
    cd $project_path_release/
    cp -r config/$env/* $release_config_path/
    echo "Copying config/ into release/config"
    if [ $? != 0 ]; then
        printf "Copying conf failed\n"
        exit 103
    fi
    printf "============================================\n\n\n"
}

printRelease(){
    printf "Print Release Directory\n"
    printf "============================================\n"
    cd $project_path
    find $release_path
    printf "============================================\n\n\n"
}
startService(){
  # 检查原文件是否存在
  if [ ! -f "./$SERVICE_FILE" ]; then
      echo "错误：当前目录下找不到 $SERVICE_FILE 文件"
      exit 1
  fi

  # 检查目标文件是否存在
  if [ -f "$TARGET_PATH" ]; then
      echo "检测到 $TARGET_PATH 已存在"

      # 检查服务是否正在运行
      if systemctl is-active --quiet $SERVICE_NAME; then
          echo "服务 $SERVICE_NAME 正在运行，正在重启..."
          systemctl restart $SERVICE_NAME
      fi
  fi

  # 复制文件
  echo "正在复制 $SERVICE_FILE 到 $TARGET_PATH ..."
  cp ./$SERVICE_FILE "$TARGET_PATH"

  # 重载 systemd
  echo "正在重载 systemd 配置..."
  systemctl daemon-reload

  # 启动服务
  echo "正在启动 $SERVICE_NAME 服务..."
  systemctl start $SERVICE_NAME

  # 检查服务状态
  echo "服务状态："
  systemctl status $SERVICE_NAME --no-pager
}
runApp(){
    printf "Run crmApp Start\n"
    printf "============================================\n"
    systemctl restart crmApp
    printf "Run App successed\n"
    printf "============================================\n\n\n"
}
printEnv
cleanDir
buildBin
copyConf
printRelease
runApp
exit 0



