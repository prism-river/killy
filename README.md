<h1 align="center">
        <br>
        <img width="300" src="presentation/images/logo.png" alt="killy">
        <br>
        <h4 align="center">Play TiDB in Minecraft!</h4>
        <br>
</h1>

<p align="center">
    <a href="https://travis-ci.org/prism-river/killy"><img src="https://travis-ci.org/prism-river/killy.svg?branch=master"></a>
	<a href="https://goreportcard.com/report/github.com/prism-river/killy"><img src="https://goreportcard.com/badge/github.com/prism-river/killy"></a>
	<a href="https://godoc.org/github.com/prism-river/killy"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
	<a href="https://libraries.io/github/prism-river/killy"><img src="https://img.shields.io/librariesio/github/prism-river/killy.svg"></a>
</p>

## Features Preview
 
### TiDB Cluster Status

<div align="center">
	<img src="./presentation/images/status.gif" alt="" width="600">
</div>

### Table Status

<div align="center">
	<img src="./presentation/images/table.png" alt="" width="600">
</div>

### Query in Minecraft

<div align="center">
	<img src="./presentation/images/query.png" alt="" width="600">
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

## Architecture

<div align="center">
	<img src="./presentation/images/arch.png" alt="" width="600">
</div>

## Acknowledgments

* Thanks [github.com/docker/dockercraft](https://github.com/docker/dockercraft) for its awesome idea.
