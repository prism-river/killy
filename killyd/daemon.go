package killyd

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // driver for database/sql
	"github.com/prism-river/killy/collectors"
)

// Table represents a table in TiDB
type Table struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Data    [][]string `json:"data"`
}

// TCPMessage defines what a message that can be
// sent or received to/from LUA scripts
type TCPMessage struct {
	Cmd  string   `json:"cmd,omitempty"`
	Args []string `json:"args,omitempty"`
	// Id is used to associate requests & responses
	ID   int         `json:"id,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// TidbEvent is one kind of Data that can
// be transported by a TCPMessage in the Data field.
type TidbEvent struct {
	TidbConnections  string
	TidbAvailHosts   []string
	TikvAvailHosts   []string
	PdAvailHosts     []string
	TidbUnavailHosts []string
	TikvUnavailHosts []string
	PdUnavailHosts   []string
	TidbNum          string
	TikvNum          string
	PdNum            string
	Totalcap         string
	Totalavail       string
	EveryTikvStatus  interface{}
	sync.RWMutex
}

type Daemon struct {
	SendData TidbEvent
	// The configuration
	ctx    *context
	Config *Config

	// tcpMessages can be used to send bytes to the Lua
	// plugin from any go routine.
	tcpMessages chan []byte
	ExitChan    chan int

	sync.Mutex
}

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Name     string `json:"name"`
	} `json:"database"`
	Interval int `json:"interval"` // in seconds
}

func NewDaemon(ctx *context) *Daemon {
	return &Daemon{
		ctx: ctx,
	}
}

func (d *Daemon) Init() {
	// load configuration
	d.Config = new(Config)
	for topicName, services := range d.ctx.killyd.Meta.Topics {
		if topicName == "tidb" {
			for _, service := range services {
				d.Config.Database.Address = service.MysqlAddress
				d.Config.Interval = service.MysqlInterval
				d.Config.Database.Name = service.Db
				d.Config.Database.Password = service.Password
				d.Config.Database.User = service.Username
				goto Exit
			}
		}
	}
Exit:
	d.tcpMessages = make(chan []byte)
	d.ExitChan = make(chan int)
}

func (d *Daemon) Exit() {
	close(d.ExitChan)
}

func (d *Daemon) Serve() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":25566")
	ln, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		d.ctx.killyd.logf(LOG_FATAL, "listen tcp error: %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			d.ctx.killyd.logf(LOG_FATAL, "tcp conn accept error: %v", err)
		}
		// no need to handle connection in a go routine
		// goproxy is used as support for one single Lua plugin.
		//d.ctx.killyd.waitGroup.Wrap(func() { d.AlwaySend(conn) })
		d.ctx.killyd.waitGroup.Wrap(func() { d.handleConn(conn) })
		select {
		case <-d.ExitChan:
			break
		default:
			continue
		}

	}
}

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
				d.ctx.killyd.logf(LOG_FATAL, "conn read error: %v", err)
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
			// if one complete message and separator is found
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
		d.ctx.killyd.logf(LOG_DEBUG, "tcpMessage: %v", string(tcpMessage))
		_, err := conn.Write(tcpMessage)
		if err != nil {
			d.ctx.killyd.logf(LOG_ERROR, "conn write error: %v", err)
			break
		}
	}
}

func (d *Daemon) handleMessage(message []byte) {
	var tcpMsg TCPMessage

	err := json.Unmarshal(message, &tcpMsg)
	if err != nil {
		d.ctx.killyd.logf(LOG_ERROR, "json unmarshal error: %v", err)
		return
	}
	d.ctx.killyd.logf(LOG_DEBUG, "handleMessage: %#v \n", tcpMsg)

	switch tcpMsg.Cmd {

	case "query":
		query := tcpMsg.Data.(string)
		db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v", d.Config.Database.User, d.Config.Database.Password, d.Config.Database.Address, d.Config.Database.Name))
		var msg TCPMessage
		if err != nil {
			msg = TCPMessage{
				Cmd:  "event",
				Args: []string{"error"},
				ID:   0,
				Data: err.Error(),
			}
		} else {
			defer db.Close()
			res, err := sqlQuery(db, query)
			if err != nil {
				msg = TCPMessage{
					Cmd:  "event",
					Args: []string{"error"},
					ID:   0,
					Data: err.Error(),
				}
			} else {
				msg = TCPMessage{
					Cmd:  "event",
					Args: []string{"result"},
					ID:   0,
					Data: res,
				}
			}
		}
		data, err := json.Marshal(&msg)
		if err != nil {
			log.Println("query error:", err)
			return
		}
		separator := []byte(string('\n'))
		d.tcpMessages <- append(data, separator...)
	}
}

// help function
func sqlQuery(db *sql.DB, query string) (*Table, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	table := Table{
		Name:    "result",
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
			return nil, err
		}
		table.Data = append(table.Data, fields)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func add(elements []string, element string) (result []string) {
	for _, e := range elements {
		if string(e) == element {
			result = elements
			return
		} else {
			result = append(result, string(e))
		}
	}
	result = append(result, element)
	return result
}

func remove(elements []string, element string) (result []string) {
	for _, e := range elements {
		if string(e) == element {
			continue
		} else {
			result = append(result, string(e))
		}
	}
	return result
}

func (d *Daemon) ConversionMinecraft(data collectors.CollectData) {
	d.SendData.Lock()
	defer d.SendData.Unlock()
	switch data.Type {
	case "tidb":
		if data.Fail {
			d.SendData.TidbAvailHosts = remove(d.SendData.TidbAvailHosts, data.Name)
			d.SendData.TidbUnavailHosts = add(d.SendData.TidbUnavailHosts, data.Name)
		} else {
			d.SendData.TidbAvailHosts = add(d.SendData.TidbAvailHosts, data.Name)
			d.SendData.TidbUnavailHosts = remove(d.SendData.TidbUnavailHosts, data.Name)
		}
		d.SendData.TidbNum = string(len(d.SendData.TidbAvailHosts))
		d.SendData.TidbConnections = string(data.Data["connections"].(int))
	case "pdtikv":
		if data.Fail {
			d.SendData.PdAvailHosts = remove(d.SendData.PdAvailHosts, data.Name)
			d.SendData.PdUnavailHosts = add(d.SendData.PdUnavailHosts, data.Name)
		} else {
			d.SendData.PdAvailHosts = add(d.SendData.PdAvailHosts, data.Name)
			d.SendData.PdUnavailHosts = remove(d.SendData.PdUnavailHosts, data.Name)
		}
		d.SendData.TikvAvailHosts = data.Data["availAddress"].([]string)
		d.SendData.TikvUnavailHosts = data.Data["unavailAddress"].([]string)
		d.SendData.TikvNum = string(len(d.SendData.TikvAvailHosts))
		d.SendData.Totalavail = string(data.Data["totalAvail"].(int))
		d.SendData.Totalcap = string(data.Data["totalcap"].(int))
		d.SendData.EveryTikvStatus = data.Data["EveryTikvStatus"]
	case "pd":
		if data.Fail {
			d.SendData.PdAvailHosts = remove(d.SendData.PdAvailHosts, data.Name)
			d.SendData.PdUnavailHosts = add(d.SendData.PdUnavailHosts, data.Name)
		} else {
			d.SendData.PdAvailHosts = add(d.SendData.PdAvailHosts, data.Name)
			d.SendData.PdUnavailHosts = remove(d.SendData.PdUnavailHosts, data.Name)
		}
	}
	tcpMsg := TCPMessage{}
	tcpMsg.Cmd = "monitor"
	tcpMsg.Args = []string{"all"}
	tcpMsg.ID = 0
	tcpMsg.Data = &d.SendData
	sendD, err := json.Marshal(&tcpMsg)
	if err != nil {
		d.ctx.killyd.logf(LOG_ERROR, "statCallback error: %v", err)

	}
	separator := []byte(string('\n'))
	d.tcpMessages <- append(sendD, separator...)
}

func (d *Daemon) StartMonitoringEvents() {
	d.ctx.killyd.logf(LOG_INFO, "Monitoring TiDB")
	// monitor table
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v", d.Config.Database.User, d.Config.Database.Password, d.Config.Database.Address, d.Config.Database.Name))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticker := time.NewTicker(time.Duration(d.Config.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
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
				table, err := sqlQuery(db, "SELECT * FROM "+tableName)
				if err != nil {
					log.Fatal(err)
				}
				table.Name = tableName
				res = append(res, *table)
			}

			tcpMsg := TCPMessage{}
			tcpMsg.Cmd = "event"
			tcpMsg.Args = []string{"table"}
			tcpMsg.ID = 0
			tcpMsg.Data = &res

			data, err := json.Marshal(&tcpMsg)
			if err != nil {
				log.Println("table monitor error:", err)
				return
			}

			separator := []byte(string('\n'))

			d.tcpMessages <- append(data, separator...)
		case <-d.ExitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}
