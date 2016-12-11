package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/fagongzi/gateway/pkg/conf"
	"github.com/fagongzi/gateway/pkg/util"
	"github.com/fagongzi/gateway/proxy"
)

var (
	cpus       = flag.Int("cpus", 1, "use cpu nums")
	logFile    = flag.String("log-file", "", "which file to record log, if not set stdout to use.")
	logLevel   = flag.String("log-level", "info", "log level.")
	configFile = flag.String("config", "", "config file")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*cpus)

	util.InitLog(*logFile)
	level := util.SetLogLevel(*logLevel)

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.PanicErrorf(err, "read config file <%s> failure.", *configFile)
	}

	cnf := &conf.Conf{}
	err = json.Unmarshal(data, cnf)
	if err != nil {
		log.PanicErrorf(err, "parse config file <%s> failure.", *configFile)
	}

	cnf.LogLevel = level

	if cnf.EnablePPROF {
		go func() {
			log.Println(http.ListenAndServe(cnf.PPROFAddr, nil))
		}()
	}

	log.Infof("conf:<%+v>", cnf)

	server := proxy.NewProxy(cnf)

	server.Start()
}
