package oidc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vivshankar/verifypass-go/internal/config"
)

const (
	introspectPath string = "/v1.0/endpoint/default/introspect"
)

var client *http.Client = http.DefaultClient

func Introspect(token string) (map[string]interface{}, error) {

	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("token", token)

	req, err := newRequest(token, "POST", introspectPath, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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
	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Errorf("Unable to decode response; body = %v, error %v", string(body), err)
		return nil, err
	}

	return payload, nil
}

func newRequest(token string, method string, path string, contentType string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", config.Tenant, path), body)
	if err != nil {
		log.Errorf("Could not create request; error %v", err)
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	return req, nil
}
