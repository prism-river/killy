<h1 align="center">
        <br>
        <img width="300" src="presentation/images/logo.png" alt="killy">
        <br>
        <h4 align="center">Play TiDB in Minecraft!</h4>
        <br>
</h1>

[![Build Status](https://travis-ci.org/prism-river/killy.svg?branch=master)](https://travis-ci.org/prism-river/killy)
[![Go Report Card](https://goreportcard.com/badge/github.com/prism-river/killy)](https://goreportcard.com/report/github.com/prism-river/killy)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/prism-river/killy)
[![](https://img.shields.io/badge/docker-supported-blue.svg)](https://godoc.org/github.com/prism-river/killy)
[![Libraries.io for GitHub](https://img.shields.io/librariesio/github/prism-river/killy.svg)](https://libraries.io/github/prism-river/killy)

## Features Preview
 
### TiDB Cluster Status

<div align="center">
	<img src="./presentation/images/status.png" alt="" width="500">
</div>

### Table Status

<div align="center">
	<img src="./presentation/images/table.png" alt="" width="500">
</div>

### Query in Minecraft

<div align="center">
	<img src="./presentation/images/querys.png" alt="" width="500">
</div>

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

## Acknowledgments

* Thanks [github.com/docker/dockercraft](https://github.com/docker/dockercraft) for its awesome idea.
