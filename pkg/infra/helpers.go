package infra

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-create-dns-record
type DNSRecord struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
	TTL     int64    `json:"ttl"`
}

// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-create-dns-record
type CfResponse struct {
	Errors []struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"messages"`
	Success bool     `json:"success"`
	Result  struct{} `json:"result"`
}

// prepareJSON is a helper function to marshal input data into JSON string, and to convert them into a pointer to bytes.Reader struct
func prepareJSON(data interface{}) (*bytes.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader([]byte(jsonData))

	return bodyReader, nil
}

func parseResponse(response *http.Response) (*CfResponse, error) {
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var cfResponse CfResponse

	err = json.Unmarshal(respData, &cfResponse)
	if err != nil {
		return nil, err
	}

	return &cfResponse, nil
}

// callCfAPI function is a helper function to call remote Cloudflare API. It is used for creating new DNS records at the moment, so it serves its one purpose only at the moment.
// func callCfAPI(authParams []string, zoneID string, payload interface{}) error {
func callCfAPI(zoneID string, payload interface{}) error {
	/*if authParams == nil || len(authParams) != 2 {
		return errors.New("given auth params are empty, or of invalid length")
	}*/

	if zoneID == "" {
		return errors.New("given zoneID is empty")
	}

	if payload == nil {
		return errors.New("payload not provided (is nil)")
	}

	bodyReader, err := prepareJSON(payload)
	if err != nil {
		return err
	}

	// TODO: make this more generic
	endpoint := "https://api.cloudflare.com/client/v4/zones/" + zoneID + "/dns_records"

	//var reqest *http.Request

	reqest, err := http.NewRequest("POST", endpoint, bodyReader)
	if err != nil {
		return err
	}

	// set request headers according to the Cloudflare API docs
	reqest.Header.Set("X-Auth-Email", os.Getenv("CF_API_EMAIL"))
	reqest.Header.Set("X-Auth-Key", os.Getenv("CF_API_TOKEN"))
	reqest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	response, err := client.Do(reqest)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	cfResponse, err := parseResponse(response)
	if err != nil {
		return err
	}

	// return a very nasty error directly from Cloudflare
	if !cfResponse.Success {
		errs := []string{
			cfResponse.Errors[0].Message,
			//cfResponse.Errors[1].Message,
		}
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
