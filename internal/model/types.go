package model

import "encoding/xml"

type People struct {
	XMLName xml.Name `xml:"people"`
	People  []Person `xml:"person"`
}

type Person struct {
	FirstName string    `xml:"firstname"`
	LastName  string    `xml:"lastname"`
	Address   []Address `xml:"address,omitempty"`
	Phone     []Phone   `xml:"phone,omitempty"`
	Family    []Family  `xml:"family,omitempty"`
}

type Family struct {
	Name    string    `xml:"name"`
	Born    int       `xml:"born"`
	Address []Address `xml:"address,omitempty"`
	Phone   []Phone   `xml:"phone,omitempty"`
}

type Address struct {
	Street   string `xml:"street"`
	City     string `xml:"city"`
	Postcode string `xml:"postcode"`
}

type Phone struct {
	Mobile   string `xml:"mobile"`
	Landline string `xml:"landline"`
}
