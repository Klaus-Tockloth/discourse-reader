package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

/*
getDiscourseData gets data from Discourse forum site.
*/
func getDiscourseData(requestURL string) (int, []byte, error) {
	// build request
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return -1, nil, fmt.Errorf("error [%w] at http.NewRequest()", err)
	}
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Api-Key", *userapikey)

	// log.Printf("outgoing HTTP request\n%v", DumpOutgoingRequest(request, true))

	// send request
	response, err := httpClient.Do(request)
	if err != nil {
		return -1, nil, fmt.Errorf("error [%w] at httpClient.Do()", err)
	}
	defer response.Body.Close()

	// log.Printf("incoming HTTP response: header fields\n%v", DumpResponse(response, true))

	// process response
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, nil, fmt.Errorf("error [%v] at io.ReadAll()", err)
	}

	return response.StatusCode, bodyData, nil
}

/*
DumpOutgoingRequest dumps outgoing HTTP request.
*/
func DumpOutgoingRequest(req *http.Request, body bool) string {
	requestDump, err := httputil.DumpRequestOut(req, body)
	if err != nil {
		return fmt.Sprintf("error [%v] at httputil.DumpRequestOut()", err)
	}

	return string(requestDump)
}

/*
DumpResponse dumps incoming HTTP response
*/
func DumpResponse(resp *http.Response, body bool) string {
	responseDump, err := httputil.DumpResponse(resp, body)
	if err != nil {
		return fmt.Sprintf("error [%v] at httputil.DumpResponse()", err)
	}

	return string(responseDump)
}
