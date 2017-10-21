package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/judwhite/go-svc/svc"
	"github.com/prism-river/killy/internal/version"
	"github.com/prism-river/killy/killyd"
)

type program struct {
	killyd *killyd.KILLYD
}

func killydFlagSet(opts *killyd.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("killyd", flag.ExitOnError)
	flagSet.Bool("version", false, "print version string")
	flagSet.String("config", "", "path to config file")
	return flagSet
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func loadmeta(configFile string) (meta killyd.Meta, err error) {
	if configFile != "" {
		_, err = toml.DecodeFile(configFile, &meta)
		if err != nil {
			return
		}
	}
	return
}

func (p *program) Start() error {
	opts := killyd.NewOptions()
	flagSet := killydFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(version.String())
		os.Exit(0)
	}

	killyd, err := killyd.New(opts)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	configFile := flagSet.Lookup("config").Value.String()
	meta, err := loadmeta(configFile)
	if err != nil {
		log.Fatalf("ERROR: failed to load config file %s - %s", configFile, err.Error())
	}
	//for k,v := range meta.Topics {
	//	fmt.Println("--")
	//	fmt.Println(k)
	//	for kk,vv := range v {
	//		fmt.Println("-")
	//		fmt.Println(kk)
	//		fmt.Println(vv.Address)
	//		fmt.Println(vv.Interval)
	//	}
	//}

	err = killyd.Loadmeta(meta)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	killyd.Main()
	p.killyd = killyd

	return nil
}

func (p *program) Stop() error {
	if p.killyd != nil {
		p.killyd.Exit()
		log.Fatalf("INFO: stop the killy")
	}
	return nil
}
