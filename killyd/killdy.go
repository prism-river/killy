package killyd

import (
	"errors"
	"log"
	"os"
	"sync"
	"sync/atomic"

	//_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	goutil "github.com/hawkingrei/golang_util"
	"github.com/mkevac/debugcharts"
	"github.com/prism-river/killy/collectors"
	"github.com/prism-river/killy/internal/lg"
)

type KILLYD struct {
	sync.RWMutex
	daemon         *Daemon
	Meta           Meta
	opts           atomic.Value
	waitGroup      goutil.WaitGroupWrapper
	exitChan       chan int
	pushinfluxChan chan *collectors.CollectData
	topicMap       map[string]*Topic
}

func New(opts *Options) (v *KILLYD, err error) {
	if opts.Logger == nil {
		opts.Logger = log.New(os.Stderr, opts.LogPrefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	}
	v = &KILLYD{
		topicMap:       make(map[string]*Topic),
		exitChan:       make(chan int),
		pushinfluxChan: make(chan *collectors.CollectData, 100000000),
	}
	v.daemon = NewDaemon(&context{v})
	opts.logLevel, err = lg.ParseLogLevel(opts.LogLevel, opts.Verbose)
	if err != nil {
		v.logf(LOG_FATAL, "%s", err)
		os.Exit(1)
	}
	v.swapOpts(opts)
	return
}

func (v *KILLYD) getOpts() *Options {
	return v.opts.Load().(*Options)
}

func (v *KILLYD) swapOpts(opts *Options) {
	v.opts.Store(opts)
}
func (v *KILLYD) Loadmeta(meta Meta) error {
	v.Meta = meta
	for topicName, services := range meta.Topics {
		topic := v.GetTopic(topicName)
		for servername, service := range services {
			topic.GetChannel(servername, service)
		}
	}
	return nil
}
func (v *KILLYD) Main() {
	v.daemon.Init()
	v.waitGroup.Wrap(func() { v.daemon.Serve() })
	v.waitGroup.Wrap(func() { v.daemon.StartMonitoringEvents() })
	v.waitGroup.Wrap(func() { v.ToMinecraft() })
	router := gin.Default()
	debugcharts.GinDebugRouter(router)
	router.Run(":8434")
}

func (v *KILLYD) Exit() {
	v.logf(LOG_INFO, "KILLYD: closing topics")
	for _, vv := range v.topicMap {
		vv.Exit()
	}
	close(v.exitChan)
	v.daemon.Exit()
	v.waitGroup.Wait()
}

// GetTopic performs a thread safe operation
// to return a pointer to a Topic object (potentially new)
func (v *KILLYD) GetTopic(topicName string) *Topic {
	// most likely, we already have this topic, so try read lock first.
	v.RLock()
	t, ok := v.topicMap[topicName]
	v.RUnlock()
	if ok {
		return t
	}

	v.Lock()
	t, ok = v.topicMap[topicName]
	if ok {
		v.Unlock()
		return t
	}
	t = NewTopic(topicName, &context{v})
	v.topicMap[topicName] = t
	v.logf(LOG_INFO, "TOPIC(%s): created", t.name)
	v.Unlock()
	return t
}

// GetExistingTopic gets a topic only if it exists
func (v *KILLYD) GetExistingTopic(topicName string) (*Topic, error) {
	v.RLock()
	defer v.RUnlock()
	topic, ok := v.topicMap[topicName]
	if !ok {
		return nil, errors.New("topic does not exist")
	}
	return topic, nil
}

// DeleteExistingTopic removes a topic only if it exists
func (v *KILLYD) DeleteExistingTopic(topicName string) error {
	v.RLock()
	topic, ok := v.topicMap[topicName]
	if !ok {
		v.RUnlock()
		return errors.New("topic does not exist")
	}
	v.RUnlock()

	// delete empties all channels and the topic itself before closing
	// (so that we dont leave any messages around)
	//
	// we do this before removing the topic from map below (with no lock)
	// so that any incoming writes will error and not create a new topic
	// to enforce ordering
	topic.Exit()

	v.Lock()
	delete(v.topicMap, topicName)
	v.Unlock()

	return nil
}

func (v *KILLYD) ToMinecraft() {
	for {
		select {
		case c := <-v.pushinfluxChan:
			v.daemon.ConversionMinecraft(*c)
		case <-v.exitChan:
			goto exit
		}
	}
exit:
	v.logf(LOG_DEBUG, "KILLYD: ToMinecraftexit")
}
