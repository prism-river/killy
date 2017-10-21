package collectors

// CollectData is all the Metric data
type CollectData struct {
	Name string
	Type string
	Data map[string]interface{}
	Fail bool
}

// Collectd is the interface of all the collector
type Collectd interface {
	Start() (CollectData, error)
}
