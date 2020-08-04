package dpcm

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

const (
	dspPath string = "dpcm/v1.0/privacy/data-subject-presentation"
)

/*
{
		"purposeId": ["string"],
		"subjectId": "string",
		"isExternalSubject": false,
		"geoIP": "string"
	}
*/
type DSPRequest struct {
	PurposeID []string `json:"purposeId"`
}

/*
{
      "name": "p7",
      "version": 0,
      "description": "marketing",
      "state": 0,
      "previousConsentApply": false,
      "similarToVersion": 0,
      "tags": [
        "marketing",
        "test"
      ],
      "category": "eula",
      "lastModifiedTime": "1593521285",
      "termsOfUse": {
        "ref": "https://www.google.com"
      },
      "disableDeleteConsent": true,
      "dataCount": 0
    },
    "purpose-id-2": {}
  }
*/
type Purpose struct {
	Name                 string               `json:"name"`
	Version              int                  `json:"version"`
	Description          string               `json:"description"`
	State                int                  `json:"state"`
	PreviousConsentApply bool                 `json:"previousConsentApply"`
	SimilarToVersion     int                  `json:"similarToVersion"`
	Tags                 []string             `json:"tags"`
	Category             string               `json:"category"`
	LastModifiedTime     string               `json:"lastModifiedTime"`
	Terms                *TermsOfUse          `json:"termsOfUse"`
	DisableDeleteConsent bool                 `json:"disableDeleteConsent"`
	DataCount            int                  `json:"dataCount"`
	AccessTypes          []*PurposeAccessType `json:"accessTypes,omitempty"`
	Attributes           []*PurposeAttribute  `json:"attributes,omitempty"`
}

/*
{
                    "id": "1",
                    "mandatory": true,
                    "accessTypes": [
                        {
                            "id": "37a77dd6-d937-4489-92c1-cbc6657d5cc3",
                            "legalCategory": "3",
                            "assentUIDefault": false
                        }
                    ]
                }
*/
type PurposeAttribute struct {
	AttributeID string               `json:"id"`
	Mandatory   bool                 `json:"mandatory"`
	AccessTypes []*PurposeAccessType `json:"accessTypes"`
}

type PurposeAccessType struct {
	AccessTypeID string `json:"id"`
}

type TermsOfUse struct {
	Ref string `json:"ref"`
}

/*
{
      "name": "uid",
      "description": "The unique identifier for the user.",
      "scope": "global",
      "sourceType": "schema",
      "datatype": "string",
      "tags": [
        "sso",
        "prov"
      ],
      "credName": "uid",
      "schemaAttribute": {
        "name": "uid",
        "attributeName": "id",
        "scimName": "id",
        "customAttribute": false
      }
    },
*/
type Attribute struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Scope       string   `json:"scope"`
	SourceType  string   `json:"sourceType"`
	DataType    string   `json:"dataType"`
	Tags        []string `json:"tags"`
}

/*
{
      "name": "RW"
    }
*/
type AccessType struct {
	Name string `json:"name"`
}

type DSPConsent struct {
	PurposeID    string `json:"purposeId"`
	AttributeID  string `json:"attributeId,omitempty"`
	AccessTypeID string `json:"accessTypeId"`
	State        int    `json:"state"`
	StartTime    int64  `json:"startTime,omitempty"`
	EndTime      int64  `json:"endTime,omitempty"`
	Version      int    `json:"version"`
}

/*
{
  "purposes": {
	"purpose-id-1": {},
  },
  "consents": {
    "consent-id-1": {
      "purposeId": "8837584f-e642-4d31-9b68-6c1769404022",
      "purposeVersion": 1,
      "isGlobal": false,
      "applicationId": "881213923715506226",
      "attributeId": "12",
      "attributeValue": "attr value",
      "accessTypeId": "770c9e99-ac0b-4b76-9011-1b87c07b82aa",
      "geoIP": "64.64.64.64",
      "geoId": "4691930",
      "state": 2,
      "createdTime": 1595834929,
      "lastModifiedTime": 1595834929,
      "startTime": 1689545816,
      "endTime": 1690445826,
      "version": 1,
      "status": 3,
      "customAttributes": [
        {
          "name": "ca1",
          "value": "cv1"
        },
        {
          "name": "ca2",
          "value": "cv2"
        }
      ]
    },
    "consent-id-2": {}
  },
  "attributes": {
    "5": {
      "name": "uid",
      "description": "The unique identifier for the user.",
      "scope": "global",
      "sourceType": "schema",
      "datatype": "string",
      "tags": [
        "sso",
        "prov"
      ],
      "credName": "uid",
      "schemaAttribute": {
        "name": "uid",
        "attributeName": "id",
        "scimName": "id",
        "customAttribute": false
      }
    },
    "attributes-id-2": {}
  },
  "accessTypes": {
    "accessTypes-id-1": {
      "name": "RW"
    },
    "accessTypes-id-2": {}
  },
  "geographies": {
    "6252001": {
      "continent_code": "NA",
      "continent_name": "North America",
      "country_code": "US",
      "country_name": "United States",
      "subdivision_one_code": "",
      "subdivision_one_name": "",
      "subdivision_two_code": "",
      "subdivision_two_name": "",
      "city_name": "",
      "is_in_european_union": "0"
    },
    "geo-id-2": {}
  }
}
*/
type DSPResponse struct {
	Purposes    map[string]*Purpose    `json:"purposes"`
	Consents    map[string]*DSPConsent `json:"consents"`
	Attributes  map[string]*Attribute  `json:"attributes"`
	AccessTypes map[string]*AccessType `json:"accessTypes"`
}

func DSP(token string, r *DSPRequest) (*DSPResponse, error) {

	var buf io.ReadWriter
	if r != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(r)
		if err != nil {
			log.Errorf("Encoding of request failed; error %v", err)
			return nil, err
		}

		log.Infof("DSP Raw Request; %s", buf)
	}

	log.Infof("DSP Request: %+v", r)

	req, err := newRequest(token, "POST", dspPath, buf)
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

	log.Infof("Raw DSP Response: %v", string(body))
	var dspResp DSPResponse
	err = json.Unmarshal(body, &dspResp)
	if err != nil {
		log.Errorf("Unable to decode response; body = %v, error %v", string(body), err)
		return nil, err
	}

	log.Infof("DSP Response: %+v", dspResp)

	return &dspResp, err
}
