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

<<<<<<< HEAD
func NewCollectd(daomen *Daemon) Collectd {
	return Collectd{
		daomen:   daomen,
=======
func NewCollectd(daemon *Daemon) *Collectd {
	return &Collectd{
		daemon:   daemon,
>>>>>>> cad1ad3b4a7894e78fcb9c19a3a212e8ee416765
		exitChan: make(chan int),
		interval: 1,
	}
}

func (c *Collectd) Start() {
	c.wg.Wrap(func() { c.GetAllTidb() })
	c.wg.Wrap(func() { c.GetAllTikv() })
	c.wg.Wrap(func() { c.GetAllPd() })
	return nil
}

func (c *Collectd) Stop() {
	close(c.exitChan)
}

func (c *Collectd) GetAllTidb() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := getAllTidb(c.daomen)
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
			err := getAllTikv(c.daomen)
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
			err := getAllPd(c.daomen)
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
