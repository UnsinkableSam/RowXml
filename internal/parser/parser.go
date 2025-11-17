package parser

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	d "RowXml.com/internal/model"
)

// Parse consumes the pipe-delimited input and builds model.People.
// It performs syntactic validation (record shapes) and returns helpful errors.
func Parse(input string) (d.People, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	var out d.People

	var currentPerson *d.Person
	var currentFamily *d.Family

	lineNo := 0
	for _, raw := range lines {
		lineNo++
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		typeToken := parts[0]
		switch typeToken {
		case "P":
			if len(parts) < 3 {
				return d.People{}, fmt.Errorf("line %d: P record must be P|First|Last", lineNo)
			}
			p := d.Person{FirstName: strings.TrimSpace(parts[1]), LastName: strings.TrimSpace(parts[2])}
			out.People = append(out.People, p)
			currentPerson = &out.People[len(out.People)-1]
			currentFamily = nil

		case "F":
			if currentPerson == nil {
				return d.People{}, fmt.Errorf("line %d: F record without a current person", lineNo)
			}
			if len(parts) < 3 {
				return d.People{}, fmt.Errorf("line %d: F record must be F|Name|Born", lineNo)
			}
			born, err := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err != nil {
				return d.People{}, fmt.Errorf("line %d: invalid birth year %q", lineNo, parts[2])
			}
			f := d.Family{Name: strings.TrimSpace(parts[1]), Born: born}
			currentPerson.Family = append(currentPerson.Family, f)
			currentFamily = &currentPerson.Family[len(currentPerson.Family)-1]

		case "A":
			if len(parts) < 4 {
				return d.People{}, fmt.Errorf("line %d: A record must be A|Street|City|Postcode", lineNo)
			}
			addr := &d.Address{Street: strings.TrimSpace(parts[1]), City: strings.TrimSpace(parts[2]), Postcode: strings.TrimSpace(parts[3])}
			if currentFamily != nil {
				currentFamily.Address = addr
			} else if currentPerson != nil {
				currentPerson.Address = addr
			} else {
				return d.People{}, fmt.Errorf("line %d: A record without a current person", lineNo)
			}

		case "T":
			if len(parts) < 3 {
				return d.People{}, fmt.Errorf("line %d: T record must be T|mobile|landline", lineNo)
			}
			ph := &d.Phone{Mobile: strings.TrimSpace(parts[1]), Landline: strings.TrimSpace(parts[2])}
			if currentFamily != nil {
				currentFamily.Phone = ph
			} else if currentPerson != nil {
				currentPerson.Phone = ph
			} else {
				return d.People{}, fmt.Errorf("line %d: T record without a current person", lineNo)
			}

		default:
			return d.People{}, fmt.Errorf("line %d: unknown record type %q", lineNo, typeToken)
		}
	}

	return out, nil
}

// ToXML marshals domain.People to pretty XML. We keep this here to avoid an
// unnecessary dependency on the CLI layer; XML-specific formatting belongs in parser package.
func ToXML(p d.People) (string, error) {
	b, err := xml.MarshalIndent(p, " ", " ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(b) + "\n", nil
}
