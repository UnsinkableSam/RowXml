package parser

import (
	"strings"
	"testing"
)

// helpers
func mustFail(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", contains)
	}
	if !strings.Contains(err.Error(), contains) {
		t.Fatalf("expected error containing %q, got %q", contains, err.Error())
	}
}

func mustOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ------------------------------------------------------------------
// Basic valid parsing tests
// ------------------------------------------------------------------

func TestSimplePerson(t *testing.T) {
	input := `P|Joe|Biden
A|White House|Washington|00000`

	people, err := Parse(input)
	mustOK(t, err)

	if len(people.People) != 1 {
		t.Fatalf("expected 1 person, got %d", len(people.People))
	}

	p := people.People[0]
	if p.FirstName != "Joe" || p.LastName != "Biden" {
		t.Fatalf("invalid person parsed: %+v", p)
	}

	if len(p.Address) != 1 {
		t.Fatalf("expected 1 address, got %d", len(p.Address))
	}

	adr := p.Address[0]
	if adr.Street != "White House" || adr.City != "Washington" || adr.Postcode != "00000" {
		t.Fatalf("invalid address parsed: %+v", adr)
	}
}

// ------------------------------------------------------------------
// Grouping tests (T and A stay grouped and multiple allowed)
// ------------------------------------------------------------------

func TestGroupingOrder(t *testing.T) {
	input := `P|Carl|Bernadotte
T|111|222
T|333|444
A|Street 1|CityA|10000
A|Street 2|CityB|20000`

	people, err := Parse(input)
	mustOK(t, err)

	p := people.People[0]

	if len(p.Phone) != 2 {
		t.Fatalf("expected 2 phones, got %d", len(p.Phone))
	}

	if p.Phone[0].Mobile != "111" || p.Phone[1].Mobile != "333" {
		t.Fatalf("phone grouping incorrect: %+v", p.Phone)
	}

	if len(p.Address) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(p.Address))
	}

	if p.Address[0].Street != "Street 1" || p.Address[1].Street != "Street 2" {
		t.Fatalf("address grouping incorrect: %+v", p.Address)
	}
}

// ------------------------------------------------------------------
// Families tests
// ------------------------------------------------------------------

func TestFamilies(t *testing.T) {
	input := `P|Alice|Doe
F|Bob|2010
A|FamilyStreet|FamilyCity|77777
F|Charlie|2015
T|999|888`

	people, err := Parse(input)
	mustOK(t, err)

	p := people.People[0]

	if len(p.Family) != 2 {
		t.Fatalf("expected 2 family members, got %d", len(p.Family))
	}

	f1 := p.Family[0]
	if f1.Name != "Bob" || f1.Born != 2010 {
		t.Fatalf("invalid family parsed: %+v", f1)
	}
	if len(f1.Address) != 1 {
		t.Fatalf("expected 1 address for family Bob")
	}

	f2 := p.Family[1]
	if f2.Name != "Charlie" || f2.Born != 2015 {
		t.Fatalf("invalid family parsed: %+v", f2)
	}
	if len(f2.Phone) != 1 {
		t.Fatalf("expected 1 phone for family Charlie")
	}
}

// ------------------------------------------------------------------
// Strict validation tests
// ------------------------------------------------------------------

func TestInvalidAddressMissingPostcode(t *testing.T) {
	input := `P|Joe|Biden
A|White House|Washington`

	_, err := Parse(input)
	mustFail(t, err, "A record must be A|Street|City|Postcode")
}

func TestInvalidPhone(t *testing.T) {
	input := `P|John|Smith
T|not-a-phone|222`

	_, err := Parse(input)
	mustFail(t, err, "invalid")
	// we assert generic "invalid" since validator messages may vary; adjust if you want a stricter substring
}

// ------------------------------------------------------------------
// Swedish Spec Example (EXPECTED FAILURE)
// ------------------------------------------------------------------
//
// The spec includes a faulty address line without postcode:
//   A|White House|Washington, D.C
//
// This breaks the A = Street|City|Postcode rule; we expect Parse to fail.
// ------------------------------------------------------------------

func TestFullExampleFromSpec_ShouldFail(t *testing.T) {
	input := `P|Victoria|Bernadotte
T|070-0101010|0459-123456
A|Haga Slott|Stockholm|101
F|Estelle|2012
A|Solliden|Öland|10002
F|Oscar|2016
T|0702-020202|02-202020
P|Joe|Biden
A|White House|Washington, D.C`

	_, err := Parse(input)
	mustFail(t, err, "A record must be A|Street|City|Postcode")
}

// ------------------------------------------------------------------
// Corrected full example that SHOULD succeed
// ------------------------------------------------------------------

func TestFullExampleCorrected(t *testing.T) {
	input := `P|Victoria|Bernadotte
T|070-0101010|0459-123456
A|Haga Slott|Stockholm|101
F|Estelle|2012
A|Solliden|Öland|10002
F|Oscar|2016
T|0702-020202|02-202020
P|Joe|Biden
A|White House|Washington, D.C|00000`

	people, err := Parse(input)
	mustOK(t, err)

	if len(people.People) != 2 {
		t.Fatalf("expected 2 persons, got %d", len(people.People))
	}
}
