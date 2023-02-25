#!/bin/bash

env=$1
target=$2
src_path="$(go env GOPATH)/src"
project_path=$(cd $(dirname $0); pwd)
project_path_release=$project_path/app/
release_path=/Users/zhangxiaodong/temp
release_project_path=$release_path/$target
release_bin_path=$release_project_path/
release_config_path=$release_project_path/config/
release_log_path=$release_project_path/log



printEnv(){
    printf "Print Env \n"
    printf "============================================\n"
    printf "Commond Params        | %s %s \n" $1  $2
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
printEnv
cleanDir
buildBin
copyConf
printRelease
exit 0
}


