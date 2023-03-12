package go_kadaster_bag

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

const (
	apiName          string = "BAG"
	apiURL           string = "https://api.bag.kadaster.nl/lvbag/individuelebevragingen/v2"
	apiURLAcceptatie string = "https://api.bag.acceptatie.kadaster.nl/lvbag/individuelebevragingen/v2"
)

type Service struct {
	apiKey        string
	useAcceptatie bool
	httpService   *go_http.Service
}

type ServiceConfig struct {
	ApiKey        string
	UseAcceptatie *bool
}

func NewService(config *ServiceConfig) (*Service, *errortools.Error) {
	if config.ApiKey == "" {
		return nil, errortools.ErrorMessage("ApiKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	var useAcceptatie = false
	if config.UseAcceptatie != nil {
		useAcceptatie = *config.UseAcceptatie
	}

	return &Service{
		apiKey:        config.ApiKey,
		useAcceptatie: useAcceptatie,
		httpService:   httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("X-Api-Key", service.apiKey)
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if len(errorResponse.Error) > 0 {
		e.SetMessage(errorResponse.ErrorDescription)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	if service.useAcceptatie {
		return fmt.Sprintf("%s/%s", apiURLAcceptatie, path)
	}
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
