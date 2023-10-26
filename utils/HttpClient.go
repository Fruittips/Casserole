package utils

import (
	"bytes"
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

type RequestType int

const (
	Read  RequestType = iota // 0 for read request
	Write                    // 1 for write request
)

var RequestTypeStr = map[RequestType]string{
	Read:  "read",
	Write: "write",
}

// Intra-system request
type Request struct {
	NodeId    NodeId
	Url       string
	Payload   *Row
	CourseId  string
	StudentId string
}

// Intra-system response
type Response struct {
	NodeId     NodeId
	StatusCode int
	Data       *Row
	Error      error
}

func (nm *NodeManager) IntraSystemRequests(requests []Request) []Response {
	var wg sync.WaitGroup
	respC := make(chan Response, len(requests))

	for _, req := range requests {
		wg.Add(1)
		go nm.nonBlockingRequest(req, &wg, &respC)
	}

	wg.Wait()
	close(respC)

	var results []Response
	for resp := range respC {
		results = append(results, resp)
	}

	return results
}

func (nm *NodeManager) nonBlockingRequest(req Request, wg *sync.WaitGroup, ch *chan Response) {
	defer wg.Done()

	timeout := time.Duration(nm.GetConfig().Timeout)
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	var err error
	var resp *http.Response
	if req.Payload != nil {
		newStudentJson, err := json.Marshal(req)
		if err == nil {
			resp, err = client.Post(req.Url, "application/json", bytes.NewBuffer(newStudentJson))
		}
	} else {
		resp, err = client.Get(req.Url)
	}

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

	*ch <- Response{NodeId: req.NodeId, StatusCode: resp.StatusCode, Data: httpResponse.Data}
	return
}
