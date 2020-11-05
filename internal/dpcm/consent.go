package dpcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	consentPath              string = "dpcm/v1.0/privacy/consents"
	getConsentsPath          string = "dpcm-mgmt/config/v1.0/privacy/consents"
	ConsentStateAllow        int    = 1
	ConsentStateDeny         int    = 2
	ConsentStateOptIn        int    = 3
	ConsentStateOptOut       int    = 4
	ConsentStateTransparency int    = 5
)

/*
{
    "subjectId": "",
    "isExternalSubject": false,
    "isGlobal": false,
    "purposeId": "",
    "attributeId": 1,
    "attributeValue": "",
    "accessTypeId": "",
    "geoIP": "",
    "state": 1,
    "startTime": 1591370011,
    "endTime": 1591370011,
    "customAttributes": [
        {
            "name": "",
            "value": ""
        }
    ]
}
*/
type Consent struct {
	ID             string `json:"id,omitempty"`
	Version        int    `json:"version,omitempty"`
	PurposeID      string `json:"purposeId"`
	PurposeName    string `json:"purposeName,omitempty"`
	AttributeID    string `json:"attributeId,omitempty"`
	AttributeName  string `json:"attributeName,omitempty"`
	AccessTypeID   string `json:"accessTypeId"`
	AccessTypeName string `json:"accessTypeName,omitempty"`
	State          int    `json:"state"`
	StartTime      int64  `json:"startTime,omitempty"`
	EndTime        int64  `json:"endTime,omitempty"`
}

type ConsentOp struct {
	Op    string   `json:"op"`
	Value *Consent `json:"value"`
}

type ConsentsResponse struct {
	Consents []Consent `json:"consents"`
}

func BulkCreateConsents(token string, r []*Consent) error {

	if len(r) < 10 {
		return createConsents(token, r)
	}

	chunkSize := 10
	for i := 0; i < len(r); i += chunkSize {
		end := i + chunkSize

		if end > len(r) {
			end = len(r)
		}

		divided := r[i:end]
		log.Infof("Attempting create consents for %v", divided)
		createConsents(token, divided)
	}

	return nil
}

func createConsents(token string, r []*Consent) error {

	var ops []*ConsentOp = make([]*ConsentOp, len(r))
	for i, val := range r {
		// Hack for bug with null EndTime
		// TODO: Remember to remove this
		val.EndTime = time.Now().AddDate(0, 0, 1).Unix()

		ops[i] = &ConsentOp{
			Op:    "add",
			Value: val,
		}
	}

	var buf io.ReadWriter
	if r != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(ops)
		if err != nil {
			log.Errorf("Encoding of request failed; error %v", err)
			return err
		}
	}

	log.Infof("Consent Request: %+v", r)

	req, err := newRequest(token, "PATCH", consentPath, buf)
	if err != nil {
		log.Errorf("Could not create request; error %v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Response failed; error %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMultiStatus || resp.StatusCode == http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Unable to read to end; error %v", err)
			return err
		}

		return fmt.Errorf("Error occurred; body = %s", string(body))
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error occurred; resp.Status = %v", resp.Status)
	}

	return nil
}

func CreateConsent(token string, r *Consent) error {

	var buf io.ReadWriter
	if r != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(r)
		if err != nil {
			log.Errorf("Encoding of request failed; error %v", err)
			return err
		}
	}

	log.Infof("Consent Request: %+v", r)

	req, err := newRequest(token, "POST", consentPath, buf)
	if err != nil {
		log.Errorf("Could not create request; error %v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Response failed; error %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error occurred; resp.Status = %v", resp.Status)
	}

	return nil
}

func GetConsents(token string) ([]Consent, error) {

	req, err := newRequest(token, "GET", getConsentsPath, nil)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error occurred; resp.Status = %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Unable to read to end; error %v", err)
		return nil, err
	}

	log.Infof("Raw Response: %v", string(body))
	var cresp ConsentsResponse
	err = json.Unmarshal(body, &cresp)
	if err != nil {
		log.Errorf("Unable to decode response; body = %v, error %v", string(body), err)
		return nil, err
	}

	log.Infof("Consents: %+v", cresp)

	return cresp.Consents, nil
}
