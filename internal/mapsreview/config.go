package mapsreview

import "strings"

type CityConfig struct {
	Name             string
	Slug             string
	SearchName       string
	DefaultPostcodes []string
	DistrictLabel    string
	DistrictsEnabled bool
}

var ActiveCity = CityConfig{
	Name:             "Gießen",
	Slug:             "giessen",
	SearchName:       "Gießen",
	DefaultPostcodes: []string{"35390", "35392", "35394", "35396", "35398"},
	DistrictLabel:    "Stadtteil",
	DistrictsEnabled: false,
}

var DefaultPostcodeSet = func() map[string]bool {
	set := make(map[string]bool, len(ActiveCity.DefaultPostcodes))
	for _, postcode := range ActiveCity.DefaultPostcodes {
		set[postcode] = true
	}
	return set
}()

func OutputPrefix() string {
	return ActiveCity.Slug
}

func ProjectTitle() string {
	return ActiveCity.Name + " Google-Bewertungen: Diffamierungs-Löschbanner"
}

func SearchQuery(query, postcode string) string {
	parts := []string{strings.TrimSpace(query), strings.TrimSpace(postcode), ActiveCity.SearchName}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, " ")
}
