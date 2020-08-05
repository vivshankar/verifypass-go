package dpcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vivshankar/verifypass-go/internal/config"
)

const (
	duaPath string = "dpcm/v1.0/privacy/data-usage-approval"
)

var client *http.Client = http.DefaultClient

type DUAItem struct {
	PurposeID    string `json:"purposeId"`
	AccessTypeID string `json:"accessTypeId"`
	AttributeID  string `json:"attributeId"`
}

type DUARequest struct {
	Trace bool       `json:"trace"`
	Items []*DUAItem `json:"items"`
}

type DUAResponseResult struct {
	Approved    bool      `json:"approved"`
	Reason      *APIError `json:"reason,omitempty"`
	AttributeID string    `json:"attributeId,omitempty"`
}

type DUAResponseItem struct {
	PurposeID      string              `json:"purposeId"`
	AttributeID    string              `json:"attributeId,omitempty"`
	AttributeValue string              `json:"attributeValue,omitempty"`
	AccessTypeID   string              `json:"accessTypeId"`
	Result         []DUAResponseResult `json:"result"`
}

type APIError struct {
	MessageId          string `json:"messageId"`
	MessageDescription string `json:"messageDescription"`
}

func DUA(token string, r *DUARequest) ([]DUAResponseItem, error) {

	var buf io.ReadWriter
	if r != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(r)
		if err != nil {
			log.Errorf("Encoding of request failed; error %v", err)
			return nil, err
		}
	}

	log.Infof("DUA Request: %v", r)

	req, err := newRequest(token, "POST", duaPath, buf)
	if err != nil {
		log.Errorf("Could not create request; error %v", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Response failed; error %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Unable to read to end; error %v", err)
		return nil, err
	}

	//var duaRespItem DUAResponseItem
	var duaResp []DUAResponseItem

	err = json.Unmarshal(body, &duaResp)
	if err != nil {
		log.Errorf("Unable to decode response; body = %v, error %v", string(body), err)
	}

	return duaResp, err
}

func newRequest(token string, method string, path string, body io.ReadWriter) (*http.Request, error) {

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", config.Tenant, path), body)
	if err != nil {
		log.Errorf("Could not create request; error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	return req, nil
}
