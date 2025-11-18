package parser

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	d "RowXml.com/internal/model"
)

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
		t := parts[0]
		switch t {
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
			addr := d.Address{Street: strings.TrimSpace(parts[1]), City: strings.TrimSpace(parts[2]), Postcode: strings.TrimSpace(parts[3])}
			if currentFamily != nil {
				currentFamily.Address = append(currentFamily.Address, addr)
			} else if currentPerson != nil {
				currentPerson.Address = append(currentPerson.Address, addr)
			} else {
				return d.People{}, fmt.Errorf("line %d: A record without a current person or family", lineNo)
			}

		case "T":
			if len(parts) < 3 {
				return d.People{}, fmt.Errorf("line %d: T record must be T|mobile|landline", lineNo)
			}
			ph := d.Phone{Mobile: strings.TrimSpace(parts[1]), Landline: strings.TrimSpace(parts[2])}
			if currentFamily != nil {
				currentFamily.Phone = append(currentFamily.Phone, ph)
			} else if currentPerson != nil {
				currentPerson.Phone = append(currentPerson.Phone, ph)
			} else {
				return d.People{}, fmt.Errorf("line %d: T record without a current person or family", lineNo)
			}

			validatePhone := func(s string) bool {
				for _, r := range s {
					if r >= '0' && r <= '9' {
						return true
					}
				}
				return false
			}

			mobile := strings.TrimSpace(parts[1])
			landline := strings.TrimSpace(parts[2])

			if !validatePhone(mobile) || !validatePhone(landline) {
				return d.People{}, fmt.Errorf("line %d: invalid phone number", lineNo)
			}

		default:
			return d.People{}, fmt.Errorf("line %d: unknown record type %q", lineNo, t)
		}
	}

	return out, nil
}

func RenderXML(p d.People) (string, error) {
	var buf strings.Builder
	// include xml header
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	// start <people>
	if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "people"}}); err != nil {
		return "", err
	}

	for _, person := range p.People {
		// <person>
		if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "person"}}); err != nil {
			return "", err
		}

		// firstname
		if err := enc.EncodeElement(person.FirstName, xml.StartElement{Name: xml.Name{Local: "firstname"}}); err != nil {
			return "", err
		}
		// lastname
		if err := enc.EncodeElement(person.LastName, xml.StartElement{Name: xml.Name{Local: "lastname"}}); err != nil {
			return "", err
		}

		// all addresses
		for _, a := range person.Address {
			if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(a.Street, xml.StartElement{Name: xml.Name{Local: "street"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(a.City, xml.StartElement{Name: xml.Name{Local: "city"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(a.Postcode, xml.StartElement{Name: xml.Name{Local: "postcode"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}}); err != nil {
				return "", err
			}
		}

		// all phones
		for _, ph := range person.Phone {
			if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "phone"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(ph.Mobile, xml.StartElement{Name: xml.Name{Local: "mobile"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(ph.Landline, xml.StartElement{Name: xml.Name{Local: "landline"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "phone"}}); err != nil {
				return "", err
			}
		}

		// families (arrival order)
		for _, f := range person.Family {
			if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "family"}}); err != nil {
				return "", err
			}
			if err := enc.EncodeElement(f.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
				return "", err
			}
			// born as text
			if err := enc.EncodeElement(fmt.Sprintf("%d", f.Born), xml.StartElement{Name: xml.Name{Local: "born"}}); err != nil {
				return "", err
			}

			// family addresses
			for _, a := range f.Address {
				if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeElement(a.Street, xml.StartElement{Name: xml.Name{Local: "street"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeElement(a.City, xml.StartElement{Name: xml.Name{Local: "city"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeElement(a.Postcode, xml.StartElement{Name: xml.Name{Local: "postcode"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}}); err != nil {
					return "", err
				}
			}

			// family phones
			for _, ph := range f.Phone {
				if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "phone"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeElement(ph.Mobile, xml.StartElement{Name: xml.Name{Local: "mobile"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeElement(ph.Landline, xml.StartElement{Name: xml.Name{Local: "landline"}}); err != nil {
					return "", err
				}
				if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "phone"}}); err != nil {
					return "", err
				}
			}

			// </family>
			if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "family"}}); err != nil {
				return "", err
			}
		}

		// </person>
		if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "person"}}); err != nil {
			return "", err
		}
	}

	// end </people>
	if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "people"}}); err != nil {
		return "", err
	}

	if err := enc.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ToXML(p d.People) (string, error) {
	return RenderXML(p)
}
