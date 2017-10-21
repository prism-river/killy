package collectors

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PdConnection struct {
	conn string
	Name string
}

var _ Collectd = &PdConnection{}

type pdDate struct {
	Members []struct {
		Name       string   `json:"name"`
		MemberID   int64    `json:"member_id"`
		PeerUrls   []string `json:"peer_urls"`
		ClientUrls []string `json:"client_urls"`
	} `json:"members"`
}

func GetPdConn(address string, name string) (conn *PdConnection, err error) {
	return &PdConnection{
		conn: address,
		Name: name,
	}, err
}

func (pc *PdConnection) stats() (result []byte, err error) {
	resp, err := http.Get("http://" + pc.conn + "/pd/api/v1/members")
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

func (pc *PdConnection) convertCollectData(data []byte, fail bool) (cd CollectData) {
	cd.Name = pc.Name
	cd.Type = "pd"
	cd.Fail = fail

	cd.Data = make(map[string]interface{})
	if fail {
		return
	}
	var da pdDate
	var availName []string
	json.Unmarshal(data, &da)
	for _, d := range da.Members {
		name := d.Name
		availName = append(availName, name)
	}
	cd.Data["availName"] = availName
	return
}

func (pc *PdConnection) Start() (data CollectData, err error) {
	result, err := pc.stats()
	if err != nil {
		return pc.convertCollectData(result, true), err
	}
	//fmt.Println(string(result))
	return pc.convertCollectData(result, false), err
}
