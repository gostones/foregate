package main

import (
	"flag"
	"fmt"
	"github.com/gostones/foregate/rp"
	"github.com/gostones/foregate/tunnel"
	"github.com/gostones/foregate/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//
var help = `
	Usage: foregate [command] [--help]

	Commands:
		server [--port $PORT] [--domain $DOMAIN]
		client --url $URL --hostport local_host:local_port --port remote [--proxy proxy_url]
		connect --url $URL --ports local:remote [--proxy proxy_url]
		
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
	case "connect":
		connect(args)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, help)
	os.Exit(1)
}

const (
	rpsIni = `
[common]
bind_addr = 0.0.0.0
bind_port = 7000

#
vhost_http_port = 1080
vhost_https_port = 1443

#
dashboard_addr = 0.0.0.0
dashboard_port = 7500
dashboard_user = admin
dashboard_pwd = password

# trace, debug, info, warn, error
log_level = info
log_max_days = 1

#
#subdomain_host =
tcp_mux = true
#
`
	//server_port, instance, service_host, service_port, remote_port
	rpcIni = `
[common]
server_addr = localhost
server_port = %v
http_proxy =

[ssh_random]
type = tcp
local_ip = %v
local_port = %v
remote_port = %v
`
	//
	listenPort = 8080

	rpsPort = 7000
)

func server(args []string) {
	flags := flag.NewFlagSet("server", flag.ContinueOnError)

	//tunnel port
	listen := flags.Int("port", parseInt(os.Getenv("PORT"), listenPort), "server listening port")

	web := flags.String("web", os.Getenv("FG_WEB"), "web url")

	flags.Parse(args)

	if *web == "" {
		//default to dashboard
		*web = "http://localhost:7500"
	}

	//
	go rp.Server(rpsIni)

	port := util.FreePort()
	proxy := fmt.Sprintf("http://localhost:%v", port)
	go serve(port, *web)

	tunnel.TunServer(*listen, proxy)
}

//
func client(args []string) {
	flags := flag.NewFlagSet("client", flag.ContinueOnError)

	//
	url := flags.String("url", os.Getenv("FG_URL"), "tunnel url")
	proxy := flags.String("proxy", "", "http proxy")
	hostPort := flags.String("hostport", "", "reverse proxy service host:port")
	port := flags.Int("port", -1, "remote reverse proxy port")

	flags.Parse(args)

	if *url == "" {
		usage()
	}

	if *hostPort == "" {
		usage()
	}

	if *port == -1 {
		usage()
	}

	lport := util.FreePort()

	remote := fmt.Sprintf("localhost:%v:localhost:%v", lport, rpsPort)

	fmt.Fprintf(os.Stdout, "remote: %v\n", remote)

	go tunnel.TunClient(*proxy, *url, remote)

	sleep := util.BackoffDuration()

	for {
		hp := strings.Split(*hostPort, ":")
		shost := hp[0]

		sport, err := strconv.Atoi(hp[1])
		if err != nil {
			panic(err)
		}

		rp.Client(fmt.Sprintf(rpcIni, lport, shost, sport, *port))

		//should never return or error
		sleep(fmt.Errorf("Reverse proxy error"))
	}
}

//
func connect(args []string) {
	flags := flag.NewFlagSet("connect", flag.ContinueOnError)

	//
	url := flags.String("url", os.Getenv("FG_URL"), "tunnel url")
	proxy := flags.String("proxy", "", "http proxy")
	ports := flags.String("ports", "", "local_port:remote_port")

	flags.Parse(args)

	if *url == "" {
		usage()
	}

	if *ports == "" {
		usage()
	}

	pa := strings.Split(*ports, ":")

	remote := fmt.Sprintf("localhost:%v:localhost:%v", pa[0], pa[1])

	fmt.Fprintf(os.Stdout, "remote: %v\n", remote)

	tunnel.TunClient(*proxy, *url, remote)
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

func serve(port int, target string) {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	scheme := os.Getenv("PG_VSCHEME")
	vhost := os.Getenv("PG_VHOST")
	if scheme == "" {
		scheme = "http"
	}
	if vhost == "" {
		vhost = "localhost:8080"
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		req.URL.Host = vhost
		req.URL.Scheme = scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = vhost

		proxy.ServeHTTP(res, req)
	})

	log.Printf("serve port: %v target: %v\n", port, target)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Println(err)
	}
}
