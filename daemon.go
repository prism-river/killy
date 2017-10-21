package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

// TCPMessage defines what a message that can be
// sent or received to/from LUA scripts
type TCPMessage struct {
	Cmd  string   `json:"cmd,omitempty"`
	Args []string `json:"args,omitempty"`
	// Id is used to associate requests & responses
	ID   int         `json:"id,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// ContainerEvent is one kind of Data that can
// be transported by a TCPMessage in the Data field.
// It describes a Docker container event. (start, stop, destroy...)
type ContainerEvent struct {
	Action    string `json:"action,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ImageRepo string `json:"imageRepo,omitempty"`
	ImageTag  string `json:"imageTag,omitempty"`
	CPU       string `json:"cpu,omitempty"`
	RAM       string `json:"ram,omitempty"`
	Running   bool   `json:"running,omitempty"`
}

// Table represents a table in TiDB
type Table struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Data    [][]string `json:"data"`
}

// Config for the daemon
type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Name     string `json:"name"`
	} `json:"database"`
	Interval int `json:"interval"` // in seconds
}

// Daemon maintains state when the dockercraft daemon is running
type Daemon struct {
	// The configuration
	Config *Config
	// Version is the version of the Docker Daemon
	Version string
	// BinaryName is the name of the Docker Binary
	BinaryName string
	// previouscpustats is a map containing the previous cpu stats we got from the
	// docker daemon through the docker remote api
	previousCPUStats map[string]*CPUStats

	// tcpMessages can be used to send bytes to the Lua
	// plugin from any go routine.
	tcpMessages chan []byte

	sync.Mutex
}

// NewDaemon returns a new instance of Daemon
func NewDaemon() *Daemon {
	return &Daemon{
		previousCPUStats: make(map[string]*CPUStats),
	}
}

// CPUStats contains the Total and System CPU stats
type CPUStats struct {
	TotalUsage  uint64
	SystemUsage uint64
}

// Init initializes a Daemon
func (d *Daemon) Init() error {
	var err error
	// load configuration
	d.Config = new(Config)
	var configFile *os.File
	configFile, err = os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(d.Config)

	d.tcpMessages = make(chan []byte)

	return nil
}

// Serve exposes a TCP server on port 25566 to handle
// connections from the LUA scripts
func (d *Daemon) Serve() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":25566")

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln("listen tcp error:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("tcp conn accept error:", err)
		}
		// no need to handle connection in a go routine
		// goproxy is used as support for one single Lua plugin.
		d.handleConn(conn)
	}
}

// StartMonitoringEvents listens for events from the
// Docker daemon and uses callback to transmit them
// to LUA scripts.
func (d *Daemon) StartMonitoringEvents() {
	log.Info("Monitoring Database Events")

	// mysql test
	go func() {
		db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v", d.Config.Database.User, d.Config.Database.Password, d.Config.Database.Address, d.Config.Database.Name))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		for {
			tables := make([]string, 0)
			rows, err := db.Query("SHOW TABLES")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				err := rows.Scan(&name)
				if err != nil {
					log.Fatal(err)
				}
				tables = append(tables, name)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}

			res := make([]Table, 0)
			for _, tableName := range tables {
				rows, err := db.Query("SELECT * FROM " + tableName)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				columns, err := rows.Columns()
				if err != nil {
					log.Fatal(err)
				}
				table := Table{
					Name:    tableName,
					Columns: columns,
					Data:    make([][]string, 0),
				}
				for rows.Next() {
					fields := make([]string, len(columns))
					pointers := make([]interface{}, len(columns))
					for i := range fields {
						pointers[i] = &fields[i]
					}
					err := rows.Scan(pointers...)
					if err != nil {
						log.Fatal(err)
					}
					table.Data = append(table.Data, fields)
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				res = append(res, table)
			}

			tcpMsg := TCPMessage{}
			tcpMsg.Cmd = "event"
			tcpMsg.Args = []string{"table"}
			tcpMsg.ID = 0
			tcpMsg.Data = &res

			data, err := json.Marshal(&tcpMsg)
			if err != nil {
				log.Println("statCallback error:", err)
				return
			}

			separator := []byte(string('\n'))

			d.tcpMessages <- append(data, separator...)

			time.Sleep(time.Duration(d.Config.Interval) * time.Second)
		}
	}()
}

// handleConn handles a TCP connection
// with a Dockercraft Lua plugin.
func (d *Daemon) handleConn(conn net.Conn) {

	go func() {
		separator := []byte(string('\n'))

		buf := make([]byte, 256)
		cursor := 0
		for {
			// resize buf if needed
			if len(buf)-cursor < 256 {
				buf = append(buf, make([]byte, 256-(len(buf)-cursor))...)
			}
			n, err := conn.Read(buf[cursor:])
			if err != nil && err != io.EOF {
				log.Fatalln("conn read error: ", err)
			}
			cursor += n

			// TODO(aduermael): check cNetwork plugin implementation
			// conn.Read doesn't seem to be blocking if there's nothing
			// to read. Maybe the broken pipe is due to an implementation
			// problem on cNetwork plugin side
			if cursor == 0 {
				<-time.After(500 * time.Millisecond)
				continue
			}
			// log.Println("TCP data read:", string(buf[:cursor]), "cursor:", cursor)

			// see if there's a complete json message in buf.
			// messages are separated with \n characters
			messages := bytes.Split(buf[:cursor], separator)
			// if one complete message and seperator is found
			// then we should have len(messages) > 1, the
			// last entry being an incomplete message or empty array.
			if len(messages) > 1 {
				shiftLen := 0
				for i := 0; i < len(messages)-1; i++ {
					// log.Println(string(messages[i]))

					msgCopy := make([]byte, len(messages[i]))
					copy(msgCopy, messages[i])

					go d.handleMessage(msgCopy)
					shiftLen += len(messages[i]) + 1
				}
				copy(buf, buf[shiftLen:])
				cursor -= shiftLen
			}
		}
	}()

	for {
		tcpMessage := <-d.tcpMessages
		log.Debug("tcpMessage:", string(tcpMessage))
		_, err := conn.Write(tcpMessage)
		if err != nil {
			log.Fatal("conn write error:", err)
		}
	}
}

// handleMessage handles a message read
// from TCP connection
func (d *Daemon) handleMessage(message []byte) {

	var tcpMsg TCPMessage

	err := json.Unmarshal(message, &tcpMsg)
	if err != nil {
		log.Println("json unmarshal error:", err)
		return
	}

	log.Debugf("handleMessage: %#v \n", tcpMsg)

	switch tcpMsg.Cmd {
	}
}

// execDockerCmd handles Docker commands
func (d *Daemon) execDockerCmd(args []string) {
	if len(args) > 0 {
		log.Debugln("execDockerCmd:", d.BinaryName, args)
		cmd := exec.Command(d.BinaryName, args...)
		err := cmd.Run() // will wait for command to return
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}
}

// Utility functions
func splitRepoAndTag(repoTag string) (string, string) {

	repo := ""
	tag := ""

	repoAndTag := strings.Split(repoTag, ":")

	if len(repoAndTag) > 0 {
		repo = repoAndTag[0]
	}

	if len(repoAndTag) > 1 {
		tag = repoAndTag[1]
	}

	return repo, tag
}

func containerEventToTCPMsg(containerEvent ContainerEvent) ([]byte, error) {
	tcpMsg := TCPMessage{}
	tcpMsg.Cmd = "event"
	tcpMsg.Args = []string{"containers"}
	tcpMsg.ID = 0
	tcpMsg.Data = &containerEvent
	data, err := json.Marshal(&tcpMsg)
	if err != nil {
		return nil, errors.New("containerEventToTCPMsg error: " + err.Error())
	}
	return data, nil
}
