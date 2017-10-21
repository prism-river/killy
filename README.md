# TiDB craft

[![Build Status](https://travis-ci.org/prism-river/killy.svg?branch=master)](https://travis-ci.org/prism-river/killy)

## 需求

1. 监控数据的显示
2. 执行 SQL
3. 数据库表的显示

## Config

Copy the config.example.json to config.json and edit it.
## Instructions

```bash
cp -r config/* Server/
cp -r Killy Server/Plugins/
mkdir bin
ln -s /usr/bin/docker bin/docker-${DOCKER_VERSION}-ce
go build
./killy &
cd ./Server
./Cuberite
```
