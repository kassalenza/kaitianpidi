#!/bin/bash
# 检查参数
if [ -z "$1" ]; then
  echo "提交信息不能为空!"
  exit 1
fi


# 检查是否在一个 git 仓库中
if [ ! -d ".git" ]; then
  echo "未找到 .git 目录，正在初始化 Git 仓库..."
  git init
  if [ $? -ne 0 ]; then
    echo "初始化 Git 仓库失败！"
    exit 1
  fi
fi


# git config
# git config --global user.name "kass_alenza"
# git config --global user.email "kassalenza@gamil.com"
git config --global user.name "张航"
git config --global user.email "zh01771702@alibaba-inc.com"


# 获取远程仓库（推送地址）
remote_repository=$(git remote -v | grep push | cut -f1)
echo $remote_repository
# exit 0
# 如果没有远程仓库，则添加
if [ -z "$remote_repository" ]; then
  # remote_url="git@github.com:kassalenzaaa/radic.git"
  # remote_url="git@code.alibaba-inc.com:zh01771702/cmder_agt.git"
  remote_url="git@code.alibaba-inc.com:zh01771702/check_ebs_write.git"
  git remote add origin "$remote_url"
  if [ $? -ne 0 ]; then
    echo "添加远程仓库失败！"
    exit 1
  fi
fi


# 获取当前分支
current_branch=$(git branch | grep '*' | cut -d' ' -f 2)
# echo $current_branch
# exit 0


# 添加所有更改
git add .
# 同步本地更新到远程仓库 ps: 这个建议手动执行（确定仓库内多余的文件无需回滚再同步！）
# git add -u


# 提交更改，提交信息为第一个参数
if ! git commit -m "$1"; then
  echo "提交失败！"
  exit 1
fi


# 推送更改到远程仓库的当前分支
if ! git push $remote_repository $current_branch; then
  echo "推送失败！"
  exit 1
fi


echo "提交成功！"