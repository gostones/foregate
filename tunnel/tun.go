package tunnel

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/jpillora/chisel/client"
	"time"
)

//var commonHelp = `
//    --pid Generate pid file in current directory
//
//    -v, Enable verbose logging
//
//    --help, This help text
//
//  Version:
//    ` /*+ chshare.BuildVersion */ + `
//
//  Read more:
//    https://github.com/jpillora/chisel
//
//`

func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	if err := ioutil.WriteFile("chisel.pid", pid, 0644); err != nil {
		log.Fatal(err)
	}
}

//var serverHelp = `
//  Usage: chisel server [options]
//
//  Options:
//
//    --host, Defines the HTTP listening host – the network interface
//    (defaults the environment variable HOST and falls back to 0.0.0.0).
//
//    --port, -p, Defines the HTTP listening port (defaults to the environment
//    variable PORT and fallsback to port 8080).
//
//    --key, An optional string to seed the generation of a ECDSA public
//    and private key pair. All commications will be secured using this
//    key pair. Share the subsequent fingerprint with clients to enable detection
//    of man-in-the-middle attacks (defaults to the CHISEL_KEY environment
//    variable, otherwise a new key is generate each run).
//
//    --authfile, An optional path to a users.json file. This file should
//    be an object with users defined like:
//      {
//        "<user:pass>": ["<addr-regex>","<addr-regex>"]
//      }
//    when <user> connects, their <pass> will be verified and then
//    each of the remote addresses will be compared against the list
//    of address regular expressions for a match. Addresses will
//    always come in the form "<host/ip>:<port>".
//
//    --auth, An optional string representing a single user with full
//    access, in the form of <user:pass>. This is equivalent to creating an
//    authfile with {"<user:pass>": [""]}.
//
//    --proxy, Specifies another HTTP server to proxy requests to when
//    chisel receives a normal HTTP request. Useful for hiding chisel in
//    plain sight.
//
//    --socks5, Allows client to access the internal SOCKS5 proxy. See
//    chisel client --help for more information.
//` + commonHelp

func TunServer(port int, proxy string) {
	//flags := flag.NewFlagSet("server", flag.ContinueOnError)
	//
	//host := flags.String("host", "", "")
	//p := flags.String("p", "", "")
	//port := flags.String("port", "", "")
	//key := flags.String("key", "", "")
	//authfile := flags.String("authfile", "", "")
	//auth := flags.String("auth", "", "")
	//proxy := flags.String("proxy", "", "")
	//socks5 := flags.Bool("socks5", false, "")
	//pid := flags.Bool("pid", false, "")
	//verbose := flags.Bool("v", false, "")
	//
	//flags.Usage = func() {
	//	fmt.Print(serverHelp)
	//	os.Exit(1)
	//}
	//flags.Parse(args)
	//
	//if *host == "" {
	//	*host = os.Getenv("HOST")
	//}
	//if *host == "" {
	//	*host = "0.0.0.0"
	//}
	//if *port == "" {
	//	*port = *p
	//}
	//if *port == "" {
	//	*port = os.Getenv("PORT")
	//}
	//if *port == "" {
	//	*port = "8080"
	//}
	//if *key == "" {
	//	*key = os.Getenv("CHISEL_KEY")
	//}
	//s, err := NewServer(&Config{
	//	KeySeed:  *key,
	//	AuthFile: *authfile,
	//	Auth:     *auth,
	//	Proxy:    *proxy,
	//	Socks5:   *socks5,
	//})

	s, err := NewServer(&Config{
		KeySeed:  "",
		AuthFile: "",
		Auth:     "",
		Proxy:    proxy,
		Socks5:   false,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.Debug = true
	//if *pid {
	//	generatePidFile()
	//}
	if err = s.Run("0.0.0.0", fmt.Sprintf("%v", port)); err != nil {
		log.Fatal(err)
	}
}

//var clientHelp = `
//  Usage: chisel client [options] <server> <remote> [remote] [remote] ...
//
//  <server> is the URL to the chisel server.
//
//  <remote>s are remote connections tunnelled through the server, each of
//  which come in the form:
//
//    <local-host>:<local-port>:<remote-host>:<remote-port>
//
//    ■ local-host defaults to 0.0.0.0 (all interfaces).
//    ■ local-port defaults to remote-port.
//    ■ remote-port is required*.
//    ■ remote-host defaults to 0.0.0.0 (server localhost).
//
//    example remotes
//
//      3000
//      example.com:3000
//      3000:google.com:80
//      192.168.0.5:3000:google.com:80
//      socks
//      5000:socks
//
//    *When the chisel server has --socks5 enabled, remotes can
//    specify "socks" in place of remote-host and remote-port.
//    The default local host and port for a "socks" remote is
//    127.0.0.1:1080. Connections to this remote will terminate
//    at the server's internal SOCKS5 proxy.
//
//  Options:
//
//    --fingerprint, A *strongly recommended* fingerprint string
//    to perform host-key validation against the server's public key.
//    You may provide just a prefix of the key or the entire string.
//    Fingerprint mismatches will close the connection.
//
//    --auth, An optional username and password (client authentication)
//    in the form: "<user>:<pass>". These credentials are compared to
//    the credentials inside the server's --authfile. defaults to the
//    AUTH environment variable.
//
//    --keepalive, An optional keepalive interval. Since the underlying
//    transport is HTTP, in many instances we'll be traversing through
//    proxies, often these proxies will close idle connections. You must
//    specify a time with a unit, for example '30s' or '2m'. Defaults
//    to '0s' (disabled).
//
//    --proxy, An optional HTTP CONNECT proxy which will be used reach
//    the chisel server. Authentication can be specified inside the URL.
//    For example, http://admin:password@my-server.com:8081
//` + commonHelp

func TunClient(proxy string, url string, remote string) {
	log.Printf("proxy: %v url: %v remote: %v", proxy, url, remote)
	//flags := flag.NewFlagSet("client", flag.ContinueOnError)
	//
	//fingerprint := flags.String("fingerprint", "", "")
	//auth := flags.String("auth", "", "")
	//keepalive := flags.Duration("keepalive", 0, "")
	//proxy := flags.String("proxy", "", "")
	//pid := flags.Bool("pid", false, "")
	//verbose := flags.Bool("v", false, "")
	//flags.Usage = func() {
	//	fmt.Print(clientHelp)
	//	os.Exit(1)
	//}
	//flags.Parse(args)
	////pull out options, put back remaining args
	//args = flags.Args()
	//if len(args) < 2 {
	//	log.Fatalf("A server and least one remote is required")
	//}
	//
	//if *auth == "" {
	//	*auth = os.Getenv("AUTH")
	//}
	//
	//c, err := chclient.NewClient(&chclient.Config{
	//	Fingerprint: *fingerprint,
	//	Auth:        *auth,
	//	KeepAlive:   *keepalive,
	//	HTTPProxy:   *proxy,
	//	Server:      args[0],
	//	Remotes:     args[1:],
	//})

	keepalive := time.Duration(12 * time.Second)

	c, err := chclient.NewClient(&chclient.Config{
		Fingerprint: "",
		Auth:        "",
		KeepAlive:   keepalive,
		HTTPProxy:   proxy,
		Server:      url,
		Remotes:     []string{remote},
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Debug = true
	//if *pid {
	//	generatePidFile()
	//}
	if err = c.Run(); err != nil {
		log.Fatal(err)
	}
}
