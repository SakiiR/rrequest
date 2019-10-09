package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/SakiiR/ReduceRequest/internal/pkg/config"
	rparser "github.com/SakiiR/ReduceRequest/internal/pkg/parser"
	"github.com/SakiiR/ReduceRequest/internal/pkg/reducer"
	"github.com/akamensky/argparse"
)

func parseProxyString(proxyString string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxyString)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return transport, nil
}

func main() {
	parser := argparse.NewParser("rrequest", "Reduce (Burp) HTTP Request File")

	requestFile := parser.File("r", "request-file", os.O_RDONLY, 0600, &argparse.Options{Required: true, Help: "Request File to reduce"})
	proxyStr := parser.String("x", "http-proxy", &argparse.Options{Required: false, Help: "HTTP proxy to send the requests through"})
	forceSSL := parser.Flag("s", "ssl", &argparse.Options{Required: false, Help: "Forces SSL"})
	k := parser.Flag("k", "disable-check-cert", &argparse.Options{Required: false, Help: "Disable SSL cert checks"})
	v := parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "Add debug output"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	logrus.SetLevel(logrus.InfoLevel)
	if *v {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Configure HTTP Client Transport configuration
	cfg := &config.Config{RequestFile: requestFile, ForceSSL: *forceSSL}
	cfg.Transport = http.DefaultTransport.(*http.Transport)
	if *proxyStr != "" {
		cfg.Transport, err = parseProxyString(*proxyStr)
		if err != nil {
			cfg.Transport = http.DefaultTransport.(*http.Transport)
		}
	}
	if *forceSSL == true {
		cfg.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: *k}
	}

	p := rparser.Parser{Config: cfg}

	logrus.Debug("Initializing ...")
	p.Init()

	reducer.ReduceRequest(&p)
}
