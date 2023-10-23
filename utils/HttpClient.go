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

type Request struct {
	NodeId int
	Url    string
}

func (nm *NodeManager) ForwardGetRequests(requests []Request) []Response {
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

	timeout := time.Duration(nm.ConfigManager.Data.Timeout)
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
