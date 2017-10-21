package killyd

import (
	"fmt"
	"sync"
	"time"

	"github.com/prism-river/killy/collectors"
)

type Consumer interface {
	UnPause()
	Pause()
	Close() error
	TimedOutMessage()
	Empty()
}

type Channel struct {
	topicName string
	name      string
	ctx       *context
	meta      ChannelsMeta

	exitFlag  int32
	exitChan  chan int
	exitMutex sync.RWMutex
}

func NewChannel(topicName string, channelName string, channelsMeta ChannelsMeta, ctx *context) *Channel {
	fmt.Println("fuck your newchannel")
	c := &Channel{
		exitChan:  make(chan int),
		topicName: topicName,
		name:      channelName,
		meta:      channelsMeta,
		ctx:       ctx,
	}
	c.ctx.killyd.waitGroup.Wrap(func() { c.StartChannel() })
	return c
}

func (c *Channel) getconn() (mc collectors.Collectd, err error) {
	switch c.topicName {
	case "tidb":
		c.ctx.killyd.logf(LOG_INFO, "work(%s,%s): start", "tidb", c.name)
		mc, err = collectors.GetTidbConn(c.meta.Address, c.name)
	case "pd":
		c.ctx.killyd.logf(LOG_INFO, "work(%s,%s): start", "pd", c.name)
		mc, err = collectors.GetPdConn(c.meta.Address, c.name)
	}
	return mc, err
}

func (c *Channel) getTikvConn() (mc collectors.Collectd, err error) {
	c.ctx.killyd.logf(LOG_INFO, "work(%s,%s): start", "tikv", c.name)
	mc, err = collectors.GetPdTikvConn(c.meta.Address, c.name)
	return mc, err
}

func (c *Channel) StartChannel() {
	var mc collectors.Collectd
	mc, err := c.getconn()
	if c.topicName == "pd" {
		c.ctx.killyd.waitGroup.Wrap(func() { c.StartTikvChannel() })
	}
	if err != nil {
		c.ctx.killyd.logf(LOG_ERROR, "work(%s): fail to get connect : %s", c.name, err.Error())
		c.Close()
		return
	}
	ticker := time.NewTicker(time.Duration(c.meta.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			data, _ := mc.Start()
			c.ctx.killyd.pushinfluxChan <- &data
			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Channel) StartTikvChannel() {
	var mc collectors.Collectd
	mc, err := c.getTikvConn()
	if err != nil {
		c.ctx.killyd.logf(LOG_ERROR, "work(%s): fail to get connect : %s", c.name, err.Error())
		c.Close()
		return
	}
	//c.ctx.killyd.logf(LOG_DEBUG, "work(%s,%s): create to collect data", c.meta.Address, c.name)
	ticker := time.NewTicker(time.Duration(c.meta.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			data, err := mc.Start()
			//c.ctx.killyd.logf(LOG_DEBUG, "work(%s,%s): start to collect data", c.meta.Address, c.name)
			if err != nil {
				c.ctx.killyd.logf(LOG_ERROR, "work(%s,%s): fail to collect data : %s", c.meta.Address, c.name, err.Error())
				continue
			}
			c.ctx.killyd.pushinfluxChan <- &data
			continue
		case <-c.exitChan:
			goto exit
		}
	}
exit:
	ticker.Stop()
}

func (c *Channel) Close() error {
	c.ctx.killyd.logf(LOG_INFO, "CHANNEL(%s): closing", c.name)
	close(c.exitChan)
	return nil
}
