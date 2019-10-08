package parser

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/SakiiR/ReduceRequest/internal/pkg/config"
)

// Parser parser representation
type Parser struct {
	Config          *config.Config
	Request         *http.Request
	Client          *http.Client
	InitialResponse *http.Response
}

// Init parse the parser request file and store the request
func (parser *Parser) Init() error {
	buf := bufio.NewReader(parser.Config.RequestFile)

	req, err := http.ReadRequest(buf)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to read request file: '%s'", err))
		return err
	}

	req.RequestURI = ""
	req.URL.Scheme = "http"
	req.URL.Host = req.Host
	if parser.Config.ForceSSL {
		req.URL.Scheme = "https"
	}

	parser.Request = req

	if parser.Config.Transport != nil {
		parser.Client = &http.Client{Transport: parser.Config.Transport}
	}
	parser.InitialResponse, err = parser.Do(parser.Request)
	if err != nil {
		return err
	}

	return nil
}

// Do sends the request via the parser HTTP Client and return the response
func (parser *Parser) Do(request *http.Request) (*http.Response, error) {
	resp, err := parser.Client.Do(request)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to communicate with the server: '%s'", err))
		return nil, err
	}
	return resp, nil
}

// DumpRequestToStdout dumps the specified request to stdout
func DumpRequestToStdout(request *http.Request) error {

	data, err := httputil.DumpRequest(request, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump request: %s", err))
		return err
	}

	fmt.Println(string(data))
	return nil
}
