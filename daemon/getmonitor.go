package daemon

import (
	"encoding/json"
	"fmt"
)

type common struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				ExportedInstance string `json:"exported_instance"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type allPd struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name             string `json:"__name__"`
				ExportedInstance string `json:"exported_instance"`
				ExportedJob      string `json:"exported_job"`
				Instance         string `json:"instance"`
				Job              string `json:"job"`
				Type             string `json:"type"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func getAllTidb(d *Daemon) (err error) {
	rawdata, err := d.Client.Query("count(tidb_domain_load_schema_total) by (exported_instance)")
	if err != nil {
		return
	}
	var data common
	var hosts []string
	var count int
	json.Unmarshal(rawdata, &data)
	for _, host := range data.Data.Result {
		hosts = append(hosts, host.Metric.ExportedInstance)
		count++
	}
	d.SendData.Lock()
	d.SendData.Tidbhosts = hosts
	d.SendData.TidbNum = count
	fmt.Println(d.SendData.Tidbhosts)
	fmt.Println(d.SendData.TidbNum)
	d.SendData.Unlock()
	return
}

func getAllTikv(d *Daemon) (err error) {
	rawdata, err := d.Client.Query("count(tikv_engine_num_subcompaction_scheduled) by (exported_instance)")
	if err != nil {
		return
	}
	var data common
	var hosts []string
	var count int
	json.Unmarshal(rawdata, &data)
	for _, host := range data.Data.Result {
		hosts = append(hosts, host.Metric.ExportedInstance)
		count++
	}
	d.SendData.Lock()
	d.SendData.Tikvhosts = hosts
	d.SendData.TikvNum = count
	fmt.Println(d.SendData.Tikvhosts)
	fmt.Println(d.SendData.TikvNum)
	d.SendData.Unlock()
	return
}

func getAllPd(d *Daemon) (err error) {
	rawdata, err := d.Client.Query("pd_cluster_status{type=\"store_up_count\"}")
	if err != nil {
		return
	}
	var data allPd
	json.Unmarshal(rawdata, &data)
	d.SendData.Lock()
	d.SendData.PdNum = data.Data.Result[0].Value[1].(string)
	fmt.Println(d.SendData.PdNum)
	d.SendData.Unlock()
	return
}
