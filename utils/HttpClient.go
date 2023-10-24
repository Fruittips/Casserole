package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

type HTTPResponse struct {
	Data         *Row   `json:"data"`
	ErrorMessage string `json:"errorMessage"`
}

type Response struct {
	NodeId     int
	StatusCode int
	Data       Row
	Error      error
}

type RequestType int

const (
	Read  RequestType = iota // 0 for read request
	Write                    // 1 for write request
)

var RequestTypeStr = map[RequestType]string{
	Read:  "read",
	Write: "write",
}

type Request struct {
	NodeId int
	Url    string
}

func (nm *NodeManager) ForwardGetRequests(requests []Request) []Response {
	// shouldn't this call an intra-system request? if this is an external request, then those other nodes will just ask all the replicas again causing inf loop
	var wg sync.WaitGroup
	respC := make(chan Response, len(requests))

	for _, req := range requests {
		wg.Add(1)
		go nm.nonBlockingGet(req, &wg, &respC)
	}

	wg.Wait()
	close(respC)

	var results []Response
	for resp := range respC {
		results = append(results, resp)
	}

	return results
}

func (nm *NodeManager) nonBlockingGet(req Request, wg *sync.WaitGroup, ch *chan Response) {
	defer wg.Done()

	timeout := time.Duration(nm.GetConfig().Timeout)
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	resp, err := client.Get(req.Url)
	if err != nil {
		*ch <- Response{NodeId: req.NodeId, Error: err}
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		*ch <- Response{NodeId: req.NodeId, StatusCode: resp.StatusCode, Error: err}
		return
	}

	var httpResponse HTTPResponse
	err = json.Unmarshal(body, &httpResponse)
	if err != nil {
		*ch <- Response{NodeId: req.NodeId, StatusCode: resp.StatusCode, Error: err}
		return
	}

	if resp.StatusCode != http.StatusOK {
		*ch <- Response{NodeId: req.NodeId, StatusCode: resp.StatusCode, Error: errors.New(httpResponse.ErrorMessage)}
		return
	}

	*ch <- Response{NodeId: req.NodeId, StatusCode: resp.StatusCode, Data: *httpResponse.Data}
	return
}
