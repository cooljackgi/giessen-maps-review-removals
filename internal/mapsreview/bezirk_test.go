package mapsreview

import "testing"

func TestAssignBezirkDisabledByDefault(t *testing.T) {
	if got := AssignBezirk(50.5840, 8.6784); got != nil {
		t.Fatalf("expected no district assignment when districts are disabled, got %#v", got)
	}
}

func TestAssignBezirkForPostcodeDisabledByDefault(t *testing.T) {
	if got := AssignBezirkForPostcode(50.5840, 8.6784, "35390"); got != nil {
		t.Fatalf("expected no postcode-based district assignment when districts are disabled, got %#v", got)
	}
}

func TestBezirkBoundariesDisabledByDefault(t *testing.T) {
	if boundaries := BezirkBoundaries(); len(boundaries) != 0 {
		t.Fatalf("expected no district boundaries when districts are disabled, got %d", len(boundaries))
	}
}
