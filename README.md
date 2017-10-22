# TiDB craft

[![Build Status](https://travis-ci.org/prism-river/killy.svg?branch=master)](https://travis-ci.org/prism-river/killy)
[![Go Report Card](https://goreportcard.com/badge/github.com/prism-river/killy)](https://goreportcard.com/report/github.com/prism-river/killy)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/prism-river/killy)
[![](https://img.shields.io/badge/docker-supported-blue.svg)](https://godoc.org/github.com/prism-river/killy)
[![Libraries.io for GitHub](https://img.shields.io/librariesio/github/prism-river/killy.svg)](https://libraries.io/github/prism-river/killy)

## 需求

1. 监控数据的显示
2. 执行 SQL
3. 数据库表的显示

## Config

Copy the config.example.json to config.json and edit it.

## Instructions

```bash
cp -r config/* Server/
git clone https://github.com/prism-river/killy-plugin Killy
cp -r Killy Server/Plugins/
make
./build/killyd -config=example.toml
cd ./Server
./Cuberite
```

## API Specification

### TCP Messages

```go
// TCPMessage defines what a message that can be
// sent or received to/from LUA scripts
type TCPMessage struct {
	Cmd  string   `json:"cmd,omitempty"`
	Args []string `json:"args,omitempty"`
	// Id is used to associate requests & responses
	ID   int         `json:"id,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
```

#### 监控

cmd == 'monitor'

#### 数据库

cmd == 'event' and args == ['table']
