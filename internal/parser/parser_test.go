package parser

import (
	"strings"
	"testing"
)

func TestParse_ValidSimple(t *testing.T) {
	input := strings.Join([]string{
		"P|Joe|Biden",
		"A|White House|Washington|00000",
	}, "\n")
	people, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(people.People) != 1 {
		t.Fatalf("expected 1 person, got %d", len(people.People))
	}
	p := people.People[0]
	if p.FirstName != "Joe" || p.LastName != "Biden" {
		t.Fatalf("wrong person parsed: %+v", p)
	}
	if len(p.Address) != 1 {
		t.Fatalf("expected 1 address, got %d", len(p.Address))
	}
}

func TestParse_FamilyAndPhones(t *testing.T) {
	input := strings.Join([]string{
		"P|Victoria|Bernadotte",
		"T|070-0101010|0459-123456",
		"A|Haga Slott|Stockholm|101",
		"F|Estelle|2012",
		"A|Solliden|Öland|10002",
		"F|Oscar|2016",
		"T|0702-020202|02-202020",
	}, "\n")
	people, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(people.People) != 1 {
		t.Fatalf("expected 1 person, got %d", len(people.People))
	}
	p := people.People[0]
	if len(p.Family) != 2 {
		t.Fatalf("expected 2 families, got %d", len(p.Family))
	}
	// addresses and phones grouped
	if len(p.Address) != 1 {
		t.Fatalf("expected 2 addresses for person, got %d", len(p.Address))
	}
	if len(p.Phone) != 1 {
		t.Fatalf("expected 2 phones for person, got %d", len(p.Phone))
	}
}

func TestParse_InvalidA(t *testing.T) {
	input := "P|Joe|Biden\nA|White House|Washington"
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for missing postcode")
	}
	if !strings.Contains(err.Error(), "A record must be A|Street|City|Postcode") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FBeforeP(t *testing.T) {
	input := "F|Child|2010"
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for F before P")
	}
	if !strings.Contains(err.Error(), "F record without a current person") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_TWithoutContext(t *testing.T) {
	input := "T|070-1111111|08-111111"
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for T without a current person or family")
	}
	if !strings.Contains(err.Error(), "T record without a current person or family") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_AWithoutContext(t *testing.T) {
	input := "A|Storgatan 1|Stockholm|11111"
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for A without a current person or family")
	}
	if !strings.Contains(err.Error(), "A record without a current person or family") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FFollowedByF(t *testing.T) {
	input := strings.Join([]string{
		"P|Anna|Svensson",
		"F|Elsa|2010",
		"F|Olivia|2012",
	}, "\n")

	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for F followed by F")
	}
	if !strings.Contains(err.Error(), "F cannot be followed by F") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InvalidP(t *testing.T) {
	input := "P|Johan"

	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for malformed P record")
	}
	if !strings.Contains(err.Error(), "P record must be P|First|Last") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InvalidT(t *testing.T) {
	input := strings.Join([]string{
		"P|Greta|Holm",
		"T|070-1234567",
	}, "\n")

	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for malformed T record")
	}
	if !strings.Contains(err.Error(), "T record must be T|mobile|landline") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_InvalidF(t *testing.T) {
	input := strings.Join([]string{
		"P|Olle|Andersson",
		"F|Saknas",
	}, "\n")

	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for malformed F record")
	}
	if !strings.Contains(err.Error(), "F record must be F|Name|Born") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParse_FFollowedByP(t *testing.T) {
	input := strings.Join([]string{
		"P|Karl|Nilsson",
		"F|Erik|2012",
		"P|Lisa|Karlsson",
	}, "\n")

	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error for P after F")
	}
	if !strings.Contains(err.Error(), "F cannot be followed by P") {
		t.Fatalf("unexpected error: %v", err)
	}
}
