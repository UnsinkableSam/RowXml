package xmlrender

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"RowXml.com/internal/model"
)

// RenderXML emits XML with deterministic ordering:
// firstname, lastname, all address, all phone, then family (each: name,born,addresses,phones)
func RenderXML(p model.People) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	// <people>
	if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "people"}}); err != nil {
		return "", err
	}

	for _, person := range p.People {
		if err := encodePerson(enc, &person); err != nil {
			return "", err
		}
	}

	// </people>
	if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "people"}}); err != nil {
		return "", err
	}
	if err := enc.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func encodePerson(enc *xml.Encoder, p *model.Person) error {
	if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "person"}}); err != nil {
		return err
	}
	if err := enc.EncodeElement(p.FirstName, xml.StartElement{Name: xml.Name{Local: "firstname"}}); err != nil {
		return err
	}
	if err := enc.EncodeElement(p.LastName, xml.StartElement{Name: xml.Name{Local: "lastname"}}); err != nil {
		return err
	}
	// addresses
	for _, a := range p.Address {
		if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(a.Street, xml.StartElement{Name: xml.Name{Local: "street"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(a.City, xml.StartElement{Name: xml.Name{Local: "city"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(a.Postcode, xml.StartElement{Name: xml.Name{Local: "postcode"}}); err != nil {
			return err
		}
		if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}}); err != nil {
			return err
		}
	}
	// phones
	for _, ph := range p.Phone {
		if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "phone"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(ph.Mobile, xml.StartElement{Name: xml.Name{Local: "mobile"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(ph.Landline, xml.StartElement{Name: xml.Name{Local: "landline"}}); err != nil {
			return err
		}
		if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "phone"}}); err != nil {
			return err
		}
	}
	// families
	for _, f := range p.Family {
		if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "family"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(f.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
			return err
		}
		if err := enc.EncodeElement(fmt.Sprintf("%d", f.Born), xml.StartElement{Name: xml.Name{Local: "born"}}); err != nil {
			return err
		}
		for _, a := range f.Address {
			if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
				return err
			}
			if err := enc.EncodeElement(a.Street, xml.StartElement{Name: xml.Name{Local: "street"}}); err != nil {
				return err
			}
			if err := enc.EncodeElement(a.City, xml.StartElement{Name: xml.Name{Local: "city"}}); err != nil {
				return err
			}
			if err := enc.EncodeElement(a.Postcode, xml.StartElement{Name: xml.Name{Local: "postcode"}}); err != nil {
				return err
			}
			if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "address"}}); err != nil {
				return err
			}
		}
		for _, ph := range f.Phone {
			if err := enc.EncodeToken(xml.StartElement{Name: xml.Name{Local: "phone"}}); err != nil {
				return err
			}
			if err := enc.EncodeElement(ph.Mobile, xml.StartElement{Name: xml.Name{Local: "mobile"}}); err != nil {
				return err
			}
			if err := enc.EncodeElement(ph.Landline, xml.StartElement{Name: xml.Name{Local: "landline"}}); err != nil {
				return err
			}
			if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "phone"}}); err != nil {
				return err
			}
		}
		if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "family"}}); err != nil {
			return err
		}
	}
	if err := enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "person"}}); err != nil {
		return err
	}
	return nil
}

