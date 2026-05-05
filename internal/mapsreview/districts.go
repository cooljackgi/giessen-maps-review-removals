package mapsreview

type Bezirk struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BezirkBoundary struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Label    string        `json:"label"`
	Polygons [][][]float64 `json:"polygons"`
}

func AssignBezirk(lat, lng float64) *Bezirk {
	return nil
}

func AssignBezirkForPostcode(lat, lng float64, postcode string) *Bezirk {
	if !ActiveCity.DistrictsEnabled || !DefaultPostcodeSet[postcode] {
		return nil
	}
	return AssignBezirk(lat, lng)
}

func AllBezirke() []Bezirk {
	return nil
}

func BezirkBoundaries() []BezirkBoundary {
	return nil
}
