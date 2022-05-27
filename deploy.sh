#!/usr/bin/env sh
# 确保脚本抛出遇到的错误
set -e
npm run build # 生成静态文件
pwd
# 复制 README.md 文件到 dist_temp 文件夹
cp -i ./README.md ./.vuepress/dist
pwd
cd ./.vuepress/dist # 进入生成的文件夹




msg="deploy"
# echo 'blog.zhequtao.com' > CNAME
if [ -z "$CODING_TOKEN" ]; then  # -z 字符串 长度为0则为true；$CODING_TOKEN来自于github仓库`Settings/Secrets`设置的私密环境变量
  msg='deploy'
  codingUrl=git@e.coding.net:cpu887/cpu887.coding.me.git
else
  codingUrl=https://13751003867:${CODING_TOKEN}@e.coding.net/cpu887/cpu887.coding.me.git
fi

path="./.git"
if [ ! -d ${path} ]; then
  echo "不存在"
else
  echo "存在"
  rm -rf ./.git
fi
git init
git add -A
git commit -m "${msg}"
git push -f $codingUrl master # 推送到coding

