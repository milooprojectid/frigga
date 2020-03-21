package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"frigga/modules/service/storm"
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

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return errors.New("Failed parsing json body")
	}

	httpReponse, err := http.Post(service.Endpoint+"/"+path, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("Failed calling service method")
	}

	body, err := ioutil.ReadAll(httpReponse.Body)
	if err != nil {
		return errors.New("Failed parsing json body")
	}
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
			Endpoint: os.Getenv("SERVICE_REST_MORBIUS_URL"),
			Protocol: "TCP",
			Methods: map[string]string{
				"sentiment": "",
			},
		},
		"storm": Service{
			Name:     "Storm",
			Endpoint: os.Getenv("SERVICE_REST_STORM_URL"),
			Protocol: "TCP",
			Methods: map[string]string{
				"summarizeText": "",
				"summarizeLink": "/link",
			},
		},
	}
	return nil
}

// InitializeGrpcServices ...
func InitializeGrpcServices() error {
	storm.Init(os.Getenv("SERVICE_GRPC_STORM_URL"))
	return nil
}
