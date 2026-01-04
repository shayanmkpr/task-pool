package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/shayanmkpr/task-pool/internal/models"
)

type CallBack struct {
	EndPoint string
	Client   *http.Client
}

func NewCallBack(url string) CallBack {
	return CallBack{
		EndPoint: url,
	}
}

type CbackResponse struct {
	Status  string
	Message string
	Error   string
}

type CbackRequestBody struct {
	WorkerID int
	Task     *models.Task
}

func (c CallBack) CallBack(workerID int, taskInfo *models.Task) (*CbackResponse, error) {
	// construct the json body.
	requestBody := CbackRequestBody{
		WorkerID: workerID,
		Task:     taskInfo,
	}

	buffer := new(bytes.Buffer)

	err := json.NewEncoder(buffer).Encode(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.EndPoint, buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	}

	decoder := json.NewDecoder(resp.Body)
	var cbackResponse CbackResponse
	err = decoder.Decode(&cbackResponse)
	if err != nil {
		return nil, err
	}

	return &cbackResponse, nil
}
