package go_kadaster_bag

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	url2 "net/url"
	"strings"
)

type AdressenConfig struct {
	WoonplaatsNaam       *string
	Postcode             *string
	OpenbareRuimteNaam   *string
	Huisnummer           *uint
	Huisnummertoevoeging *string
	Huisletter           *string
	Query                *string
	ExacteMatch          *bool
	InclusiefEindStatus  *bool
}

type AdressenResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Embedded struct {
		Adressen []Adres `json:"adressen"`
	} `json:"_embedded"`
}

type Adres struct {
	OpenbareRuimteNaam               string   `json:"openbareRuimteNaam"`
	KorteNaam                        string   `json:"korteNaam"`
	Huisnummer                       int      `json:"huisnummer"`
	Postcode                         string   `json:"postcode"`
	WoonplaatsNaam                   string   `json:"woonplaatsNaam"`
	NummeraanduidingIdentificatie    string   `json:"nummeraanduidingIdentificatie"`
	OpenbareRuimteIdentificatie      string   `json:"openbareRuimteIdentificatie"`
	WoonplaatsIdentificatie          string   `json:"woonplaatsIdentificatie"`
	AdresseerbaarObjectIdentificatie string   `json:"adresseerbaarObjectIdentificatie"`
	PandIdentificaties               []string `json:"pandIdentificaties"`
	Adresregel5                      string   `json:"adresregel5"`
	Adresregel6                      string   `json:"adresregel6"`
	Links                            struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		OpenbareRuimte struct {
			Href string `json:"href"`
		} `json:"openbareRuimte"`
		Nummeraanduiding struct {
			Href string `json:"href"`
		} `json:"nummeraanduiding"`
		Woonplaats struct {
			Href string `json:"href"`
		} `json:"woonplaats"`
		Adres struct {
			Href string `json:"href"`
		} `json:"adres"`
		Panden []struct {
			Href string `json:"href"`
		} `json:"panden"`
	} `json:"_links"`
}

func (service *Service) cleanQuery(rawQuery string) string {
	return service.rQuery.ReplaceAllString(strings.ReplaceAll(strings.ReplaceAll(rawQuery, "\t", ""), "\n", ""), " ")
}

func (service *Service) Adressen(config *AdressenConfig) (*[]Adres, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config must not be a nil pointer")
	}

	var values = url2.Values{}
	if config.WoonplaatsNaam != nil {
		values.Set("woonplaatsNaam", service.cleanQuery(*config.WoonplaatsNaam))
	}
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
	if config.Query != nil {
		values.Set("q", service.cleanQuery(*config.Query))
	}
	if config.ExacteMatch != nil {
		values.Set("exacteMatch", fmt.Sprintf("%v", *config.ExacteMatch))
	}
	if config.InclusiefEindStatus != nil {
		values.Set("inclusiefEindStatus", fmt.Sprintf("%v", *config.InclusiefEindStatus))
	}

	var response AdressenResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("adressen?%s", values.Encode())),
		BodyModel:     config,
		ResponseModel: &response,
	}
	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Embedded.Adressen, nil
}
