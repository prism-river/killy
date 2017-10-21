package collectors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type PdTikvConnection struct {
	conn string
	Name string
}

var _ Collectd = &PdTikvConnection{}

type TikvStore struct {
	Address   string `json:"address"`
	Capacity  string `json:"capacity"`
	Available string `json:"available"`
}

type pdTikvDate struct {
	Count  int `json:"count"`
	Stores []struct {
		Store struct {
			ID        int    `json:"id"`
			Address   string `json:"address"`
			State     int    `json:"state"`
			StateName string `json:"state_name"`
		} `json:"store"`
		Status struct {
			Capacity        string    `json:"capacity"`
			Available       string    `json:"available"`
			LeaderWeight    int       `json:"leader_weight"`
			RegionCount     int       `json:"region_count"`
			RegionWeight    int       `json:"region_weight"`
			RegionScore     int       `json:"region_score"`
			StartTs         time.Time `json:"start_ts"`
			LastHeartbeatTs time.Time `json:"last_heartbeat_ts"`
			Uptime          string    `json:"uptime"`
		} `json:"status"`
	} `json:"stores"`
}

func GetPdTikvConn(address string, name string) (conn *PdTikvConnection, err error) {
	return &PdTikvConnection{
		conn: address,
		Name: name,
	}, err
}

func (pc *PdTikvConnection) stats() (result []byte, err error) {
	resp, err := http.Get("http://" + pc.conn + "/pd/api/v1/stores")
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

func (pc *PdTikvConnection) convertCollectData(data []byte, fail bool) (cd CollectData) {
	cd.Name = pc.Name
	cd.Type = "pdtikv"
	cd.Fail = fail

	cd.Data = make(map[string]interface{})
	if fail {
		return
	}
	var da pdTikvDate

	json.Unmarshal(data, &da)

	var availAddress []string
	var EveryTikvStatus []TikvStore
	totalAvail := 0
	totalcap := 0
	for _, d := range da.Stores {
		var kv TikvStore
		name := d.Store.Address
		kv.Address = d.Store.Address
		state := d.Store.State
		if state == 0 {
			availAddress = append(availAddress, name)
		}
		cap := d.Status.Capacity
		kv.Capacity = cap
		avail := d.Status.Available
		kv.Available = avail
		numavail, _ := strconv.Atoi(avail)
		totalAvail = totalAvail + numavail
		numcap, _ := strconv.Atoi(cap)
		totalcap = totalcap + numcap
		EveryTikvStatus = append(EveryTikvStatus, kv)
	}
	cd.Data["count"] = da.Count
	cd.Data["totalAvail"] = totalAvail
	cd.Data["totalcap"] = totalcap
	cd.Data["availAddress"] = availAddress
	cd.Data["EveryTikvStatus"] = EveryTikvStatus
	return
}

func (pc *PdTikvConnection) Start() (data CollectData, err error) {
	result, err := pc.stats()
	if err != nil {
		return pc.convertCollectData(result, true), err
	}
	//fmt.Println(string(result))
	return pc.convertCollectData(result, false), err
}
