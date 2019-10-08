package reducer

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/SakiiR/ReduceRequest/internal/pkg/parser"
)

func reduceURIParameters(request *http.Request, parser *parser.Parser) http.Request {
	params := request.URL.Query()
	for key, values := range params {
		// Remove the current query key
		params.Del(key)
		// Construct the new parameters
		request.URL.RawQuery = params.Encode()

		status, _ := validateResponse(parser.InitialResponse, request, parser)
		if status == true {
			fmt.Println(fmt.Sprintf("Ok, parameter %s is useless", key))
		} else {
			for _, value := range values {
				params.Add(key, value)
			}
			request.URL.RawQuery = params.Encode()
		}
	}

	return *request
}

func reduceHeaders(request *http.Request, parser *parser.Parser) http.Request {

	headers := request.Header
	for key, values := range headers {

		request.Header.Del(key)

		status, _ := validateResponse(parser.InitialResponse, request, parser)
		if status == true {
			fmt.Println(fmt.Sprintf("Ok, header %s is useless", key))
		} else {
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}

	}

	return *request
}

func reduceCookies(request *http.Request, parser *parser.Parser) http.Request {

	return *request
}

// ReduceRequest reduces request
func ReduceRequest(parser *parser.Parser) (*http.Request, error) {
	r := *parser.Request
	// TODO: iterate over URI parameters
	r = reduceURIParameters(&r, parser)
	r = reduceHeaders(&r, parser)
	r = reduceCookies(&r, parser)
	// TODO: iterate over data parameters if form
	// TODO: iterate over json parameters if json
	// TODO: iterate over xml parameters if xml
	DumpRequestToStdout(&r)
	return nil, nil
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

// validateResponse checks the request lengths to identify valid/invalid request
func validateResponse(initialResponse *http.Response, request *http.Request, parser *parser.Parser) (bool, error) {

	response, err := parser.Do(request)
	if err != nil {
		return false, err
	}

	data1, err := httputil.DumpResponse(response, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump response: %s", err))
		return false, err
	}

	data2, err := httputil.DumpResponse(initialResponse, true)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to dump response: %s", err))
		return false, err
	}

	if len(data1) != len(data2) {
		return false, nil
	}

	return true, nil
}
