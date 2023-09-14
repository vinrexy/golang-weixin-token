
# 微信 accesstoken 中控服

## 架构
  基于echo，支持多实例分布式部署
  
## 部署
  创建bin目录 并将doc下的conf目录拷贝到bin目录，修改对应的 RedisUrl 和 Port

## 构建
  本地执行 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wxtkserver src/main.go` 生成目标机器的执行文件，复制到linux机器上运行即可（尽量不要用源码跑，三方库导入太麻烦）

## 启动
  nohup ./wxtkserver &

## 使用手册
  doc/wxtoken-guide

