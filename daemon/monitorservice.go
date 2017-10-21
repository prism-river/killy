package daemon

import (
	"time"

	"github.com/ngaut/log"

	goutil "github.com/hawkingrei/golang_util"
)

type Collectd struct {
	wg       goutil.WaitGroupWrapper
	daemon   *Daemon
	exitChan chan int
	interval int
}

func NewCollectd(daemon *Daemon) *Collectd {
	return &Collectd{
		daemon:   daemon,
		exitChan: make(chan int),
		interval: 1,
	}
}

func (c *Collectd) Start() {
	c.wg.Wrap(func() { c.GetAllTidb() })
	c.wg.Wrap(func() { c.GetAllTikv() })
	c.wg.Wrap(func() { c.GetAllPd() })
	c.wg.Wrap(func() { c.GetPdDown() })
	c.wg.Wrap(func() { c.GetPdOffline() })
}

func (c *Collectd) Stop() {
	close(c.exitChan)
}

func (c *Collectd) GetAllTidb() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := getAllTidb(c.daemon)
			if err != nil {
				log.Error(err)
			}
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
			err := getAllTikv(c.daemon)
			if err != nil {
				log.Error(err)
			}
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
			err := getAllPd(c.daemon)
			if err != nil {
				log.Error(err)
			}
			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Collectd) GetPdOffline() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := getPdOffline(c.daemon)
			if err != nil {
				log.Error(err)
			}
			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Collectd) GetPdDown() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := getPdDown(c.daemon)
			if err != nil {
				log.Error(err)
			}
			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}
