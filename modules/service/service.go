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
	Name      string
	Endpoint  string
	Protocol  string
	Functions map[string]string
}

// Call ...
func (s *Service) Call(method string, payload interface{}, output interface{}) error {
	// var result map[string]interface{}

	path, ok := s.Functions[method]
	if !ok {
		return errors.New("No function registered")
	}

	requestBody, _ := json.Marshal(payload)
	httpReponse, err := http.Post(s.Endpoint+"/"+path, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New("No function registered")
	}

	body, _ := ioutil.ReadAll(httpReponse.Body)
	json.Unmarshal(body, &output)

	defer httpReponse.Body.Close()
	return nil
}

// GetService ...
func GetService(name string) (Service, error) {
	switch name {
	case "morbius":
		return Service{
			Name:     "Morbius",
			Endpoint: os.Getenv("SERVICE_MORBIUS_URL"),
			Protocol: "TCP",
			Functions: map[string]string{
				"sentiment": "",
			},
		}, nil

	case "storm":
		return Service{
			Name:     "Storm",
			Endpoint: os.Getenv("SERVICE_STORM_URL"),
			Protocol: "TCP",
			Functions: map[string]string{
				"summarizeText": "",
				"summarizeLink": "/link",
			},
		}, nil

	default:
		return Service{}, errors.New("Unknown Service")
	}
}
