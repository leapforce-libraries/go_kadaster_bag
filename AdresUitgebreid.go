package go_kadaster_bag

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	url2 "net/url"
)

type AdressenUitgebreidConfig struct {
	Postcode                         *string
	Huisnummer                       *uint
	Huisnummertoevoeging             *string
	Huisletter                       *string
	ExacteMatch                      *bool
	AdresseerbaarObjectIdentificatie *string
	WoonplaatsNaam                   *string
	OpenbareRuimteNaam               *string
	Query                            *string
	InclusiefEindStatus              *bool
}

type AdressenUitgebreidResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Embedded struct {
		Adressen []Adres `json:"adressen"`
	} `json:"_embedded"`
}

func (service *Service) AdressenUitgebreid(config *AdressenUitgebreidConfig) (*[]Adres, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be a nil pointer")
	}

	var values = url2.Values{}
	if config.Postcode != nil {
		values.Set("postcode", *config.Postcode)
	}
	if config.OpenbareRuimteNaam != nil {
		values.Set("openbareRuimteNaam", service.cleanQuery(*config.OpenbareRuimteNaam))
	}
	if config.Huisnummer != nil {
		values.Set("huisnummer", fmt.Sprintf("%v", *config.Huisnummer))
	}
	if config.Huisnummertoevoeging != nil {
		values.Set("huisnummertoevoeging", *config.Huisnummertoevoeging)
	}
	if config.Huisletter != nil {
		values.Set("huisletter", *config.Huisletter)
	}
	if config.ExacteMatch != nil {
		values.Set("exacteMatch", fmt.Sprintf("%v", *config.ExacteMatch))
	}
	if config.AdresseerbaarObjectIdentificatie != nil {
		values.Set("adresseerbaarObjectIdentificatie", service.cleanQuery(*config.AdresseerbaarObjectIdentificatie))
	}
	if config.WoonplaatsNaam != nil {
		values.Set("woonplaatsNaam", service.cleanQuery(*config.WoonplaatsNaam))
	}
	if config.OpenbareRuimteNaam != nil {
		values.Set("openbareRuimteNaam", service.cleanQuery(*config.OpenbareRuimteNaam))
	}
	if config.Query != nil {
		values.Set("q", service.cleanQuery(*config.Query))
	}
	if config.InclusiefEindStatus != nil {
		values.Set("inclusiefEindStatus", fmt.Sprintf("%v", *config.InclusiefEindStatus))
	}

	var header = http.Header{}
	header.Set("Accept-Crs", "epsg:28992")

	var response AdressenUitgebreidResponse

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodGet,
		Url:               service.url(fmt.Sprintf("adressenuitgebreid?%s", values.Encode())),
		BodyModel:         config,
		ResponseModel:     &response,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Embedded.Adressen, nil
}
