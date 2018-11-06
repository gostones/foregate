package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gostones/foregate/rp"
	"github.com/gostones/foregate/tunnel"
	"github.com/gostones/foregate/util"
	"os"
	"strconv"
	"strings"
)

//
var help = `
	Usage: foregate [command] [--help]

	Commands:
		server    - server
		client    - service worker

		server [--bind $PORT]
`

func main() {

	flag.Bool("help", false, "")
	flag.Bool("h", false, "")
	flag.Usage = func() {}
	flag.Parse()

	args := flag.Args()

	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	//
	switch subcmd {
	case "server":
		server(args)
	case "client":
		client(args)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, help)
	os.Exit(1)
}

const (
	rps = `
[common]
bind_port = %v
`
	//server_port, instance, service_host, service_port, remote_port
	rpc = `
[common]
server_addr = localhost
server_port = %v
http_proxy =

[rpc%v]
type = tcp
local_ip = %v
local_port = %v
remote_port = %v
`
	bindPort = 8080

	rpsPort = 8000

	fgPort = 8088
)

func server(args []string) {
	flags := flag.NewFlagSet("server", flag.ContinueOnError)

	//tunnel port
	bind := flags.Int("bind", parseInt(os.Getenv("PORT"), bindPort), "")

	rport := flags.Int("rps", parseInt(os.Getenv("RPS_PORT"), rpsPort), "")

	flags.Parse(args)

	//
	go rp.Server(fmt.Sprintf(rps, *rport))

	tunnel.TunServer(fmt.Sprintf("%v", *bind))
}

//
func rpClient(host string, lport, rport int) error {
	hostPort := strings.Split(host, ":")
	shost := hostPort[0]

	sport, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return err
	}

	rp.Client(fmt.Sprintf(rpc, lport, rport, shost, sport, rport))

	//should never return or error
	return errors.New("failed to connect to RP server")
}

func client(args []string) {
	flags := flag.NewFlagSet("client", flag.ContinueOnError)

	port := flags.Int("port", parseInt(os.Getenv("FG_PORT"), fgPort), "")
	url := flags.String("url", os.Getenv("FG_URL"), "")
	proxy := flags.String("proxy", "", "")
	toHostPort := flags.String("hostport", "", "reverse proxy service host:port")

	flags.Parse(args)

	if *url == "" {
		usage()
	}

	if *toHostPort == "" {
		usage()
	}

	lport := util.FreePort()

	remote := fmt.Sprintf("localhost:%v:localhost:%v", lport, rpsPort)

	fmt.Fprintf(os.Stdout, "remote: %v\n", remote)

	go tunnel.TunClient(*proxy, *url, remote)

	sleep := util.BackoffDuration()

	for {
		rc := rpClient(*toHostPort, lport, *port)

		sleep(fmt.Errorf("error: %v", rc))
	}
}

func parseInt(s string, v int) int {
	if s == "" {
		return v
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		i = v
	}
	return i
}
