package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

// Service ...
type Service struct {
	Name     string
	Endpoint string
	Protocol string
	Methods  map[string]string
}

type serviceMap map[string]Service

// SentimentResult ...
type SentimentResult struct {
	Data struct {
		Class       int    `json:"class"`
		Description string `json:"description"`
	} `json:"data"`
	Message string `json:"message"`
}

// SummarizationResult ...
type SummarizationResult struct {
	Data struct {
		Summary string `json:"summary"`
	} `json:"data"`
	Message string `json:"message"`
}

// CallSync ...
func (s *serviceMap) CallSync(serviceName string, methodName string, payload interface{}, output interface{}) error {

	service, ok := (*s)[serviceName]
	if !ok {
		return errors.New("No service registered")
	}

	path, ok := service.Methods[methodName]
	if !ok {
		return errors.New("No methos registered")
	}

	requestBody, _ := json.Marshal(payload)
	httpReponse, err := http.Post(service.Endpoint+"/"+path, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("Failed calling service method")
	}

	body, _ := ioutil.ReadAll(httpReponse.Body)
	json.Unmarshal(body, &output)

	defer httpReponse.Body.Close()
	return nil
}

// All ...
var All serviceMap

// InitializeServices ...
func InitializeServices() error {
	All = serviceMap{
		"morbius": Service{
			Name:     "Morbius",
			Endpoint: os.Getenv("SERVICE_MORBIUS_URL"),
			Protocol: "TCP",
			Methods: map[string]string{
				"sentiment": "",
			},
		},
		"storm": Service{
			Name:     "Storm",
			Endpoint: os.Getenv("SERVICE_STORM_URL"),
			Protocol: "TCP",
			Methods: map[string]string{
				"summarizeText": "",
				"summarizeLink": "/link",
			},
		},
	}
	return nil
}
