package collectors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TidbConnection struct {
	conn string
	Name string
}

var _ Collectd = &TidbConnection{}

type tidbDate struct {
	Connections int `json:"connections"`
}

func GetTidbConn(address string, name string) (conn *TidbConnection, err error) {
	return &TidbConnection{
		conn: address,
		Name: name,
	}, err
}

func (tc *TidbConnection) stats() (result []byte, err error) {
	resp, err := http.Get("http://" + tc.conn + "/status")
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	return body, err
}

func (tc *TidbConnection) convertCollectData(data []byte, fail bool) (cd CollectData) {
	cd.Name = tc.Name
	cd.Type = "tidb"
	cd.Fail = fail

	cd.Data = make(map[string]interface{})
	if fail {
		return
	}
	var da tidbDate
	json.Unmarshal(data, &da)
	cd.Data["connections"] = da.Connections
	return
}

func (tc *TidbConnection) Start() (data CollectData, err error) {
	result, err := tc.stats()
	if err != nil {
		return tc.convertCollectData(result, true), err
	}
	//fmt.Println(string(result))
	return tc.convertCollectData(result, false), err
}
