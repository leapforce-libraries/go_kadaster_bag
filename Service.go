package go_kadaster_bag

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	apiName          string = "BAG"
	apiURL           string = "https://api.bag.kadaster.nl/lvbag/individuelebevragingen/v2"
	apiURLAcceptatie string = "https://api.bag.acceptatie.kadaster.nl/lvbag/individuelebevragingen/v2"
	// this regex appears in error response if postcode has invalid format
	// we use it in the function ValidatePostcode which provides a way to check the postcode format before calling the API itself
	regexPostcode string = `^[1-9]{1}[0-9]{3}[ ]{0,1}[a-zA-Z]{2}$`
	regexQuery    string = `[^\w\s,-]`
)

type Service struct {
	apiKey        string
	useAcceptatie bool
	httpService   *go_http.Service
	rPostcode     *regexp.Regexp
	rQuery        *regexp.Regexp
}

type ServiceConfig struct {
	ApiKey        string
	UseAcceptatie *bool
}

func (service *Service) ValidatePostcode(postcode string) bool {
	return service.rPostcode.Match([]byte(postcode))
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
		rPostcode:     regexp.MustCompile(regexPostcode),
		rQuery:        regexp.MustCompile(regexQuery),
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	if requestConfig.NonDefaultHeaders != nil {
		header = *requestConfig.NonDefaultHeaders
	}
	header.Set("X-Api-Key", service.apiKey)
	header.Set("Accept", "application/hal+json")
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	attempts := 10
	for {
		request, response, e := service.httpService.HttpRequest(requestConfig)
		if response != nil && attempts > 0 {
			if response.StatusCode == http.StatusTooManyRequests {
				kadasterRateLimitDayLimitRemaining := response.Header.Get("Kadaster-RateLimit-DayLimit-Remaining")
				if kadasterRateLimitDayLimitRemaining != "" {
					kadasterRateLimitDayLimitRemainingInt, err := strconv.ParseInt(kadasterRateLimitDayLimitRemaining, 10, 64)
					if err == nil {
						if kadasterRateLimitDayLimitRemainingInt <= 0 {
							return nil, nil, errortools.ErrorMessage("Daily quotum exceeded")
						}
					}
				}
				rateLimitReset := response.Header.Get("RateLimit-Reset")
				if rateLimitReset != "" {
					rateLimitResetInt, err := strconv.ParseInt(rateLimitReset, 10, 64)
					if err == nil {
						if rateLimitResetInt > 0 {
							time.Sleep(time.Duration(rateLimitResetInt+1) * time.Second)
							attempts--
							continue
						}
					}
				}
			}
		}
		if errorResponse.Title != "" {
			e.SetMessage(errorResponse.Title)
		}

		return request, response, e
	}
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
