package googleapi

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"errors"
)

const (
	pagespeedApiTemplate = "https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=%s&key=%s&prettyprint=false"
)

type APIService interface {
	GetPagespeedResults(target string) (result *Result, err error)
}

type Service struct {
	apiKey string
}

func NewGoogleAPIService(apiKey string) APIService {
	return Service{
		apiKey: apiKey,
	}
}
func (service Service) GetPagespeedResults(target string) (result *Result, err error) {
	requestUrl := fmt.Sprintf(pagespeedApiTemplate, target, service.apiKey)
	fmt.Println(requestUrl)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("Invalid response status " + resp.Status + " and body " + string(responseBody))
	}

	return ParseResultFromReader(resp.Body)
}

func ParseResultFromReader(reader io.Reader) (result *Result, err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ParseResultFromData(data)
}

func ParseResultFromData(data []byte) (result *Result, err error) {
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return
}