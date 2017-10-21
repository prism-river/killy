package daemon

import (
	"time"

	goutil "github.com/hawkingrei/golang_util"
)

type Collectd struct {
	wg       goutil.WaitGroupWrapper
	daomen   *Daemon
	exitChan chan int
	interval int
}

func NewCollectd(daomen *Daomen) collectd {
	return &Collectd{
		daomen:   daomen,
		exitChan: make(chan int),
		interval: 1,
	}

}

func (c *Collectd) Start() error {
	c.wg.Wrap(func() { c.GetAllTidb() })
	c.wg.Wrap(func() { c.GetAllTikv() })
	c.wg.Wrap(func() { c.GetAllPd() })
}

func (c *Collectd) Stop() {
	close(c.exitChan)
}

func (c *Collectd) GetAllTidb() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:

			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Collectd) GetAllTikv() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:

			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Collectd) GetAllPd() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:

			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}
