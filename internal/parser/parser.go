package parser

import (
	"fmt"
	"strconv"

	"RowXml.com/internal/model"
)

// Parse implements a clear finite-state parser over line tokens.
// Returns a fully populated model.People or ParseError.
func Parse(input string) (model.People, error) {
	toks := lexLines(input)
	var out model.People

	var currentPerson *model.Person
	var currentFamily *model.Family

	var prevType string

	for i, tk := range toks {

		lineNo := i + 1

		switch prevType {
		case "F":
			if tk.typ != "T" && tk.typ != "A" {
				return model.People{}, newParseError(lineNo, fmt.Sprintf("F cannot be followed by %s", tk.typ))
			}
		case "P":
			if tk.typ != "T" && tk.typ != "A" && tk.typ != "F" {
				return model.People{}, newParseError(lineNo, fmt.Sprintf("P cannot be followed by %s", tk.typ))
			}
		}

		switch tk.typ {
		case "P":
			if len(tk.parts) < 3 {
				return model.People{}, newParseError(lineNo, "P record must be P|First|Last")
			}
			p := model.Person{FirstName: tk.parts[1], LastName: tk.parts[2]}
			out.People = append(out.People, p)
			currentPerson = &out.People[len(out.People)-1]
			currentFamily = nil

		case "F":
			if currentPerson == nil {
				return model.People{}, newParseError(lineNo, "F record without a current person")
			}
			if len(tk.parts) < 3 {
				return model.People{}, newParseError(lineNo, "F record must be F|Name|Born")
			}
			born, err := strconv.Atoi(tk.parts[2])
			if err != nil {
				return model.People{}, newParseError(lineNo, fmt.Sprintf("invalid birth year %q", tk.parts[2]))
			}
			f := model.Family{Name: tk.parts[1], Born: born}
			currentPerson.Family = append(currentPerson.Family, f)
			currentFamily = &currentPerson.Family[len(currentPerson.Family)-1]

		case "A":
			if len(tk.parts) < 4 {
				return model.People{}, newParseError(lineNo, "A record must be A|Street|City|Postcode")
			}
			addr := model.Address{Street: tk.parts[1], City: tk.parts[2], Postcode: tk.parts[3]}
			if currentFamily != nil {
				currentFamily.Address = append(currentFamily.Address, addr)
			} else if currentPerson != nil {
				currentPerson.Address = append(currentPerson.Address, addr)
			} else {
				return model.People{}, newParseError(lineNo, "A record without a current person or family")
			}

		case "T":
			if len(tk.parts) < 3 {
				return model.People{}, newParseError(lineNo, "T record must be T|mobile|landline")
			}
			ph := model.Phone{Mobile: tk.parts[1], Landline: tk.parts[2]}
			if currentFamily != nil {
				currentFamily.Phone = append(currentFamily.Phone, ph)
			} else if currentPerson != nil {
				currentPerson.Phone = append(currentPerson.Phone, ph)
			} else {
				return model.People{}, newParseError(lineNo, "T record without a current person or family")
			}
		default:
			return model.People{}, newParseError(lineNo, fmt.Sprintf("unknown record type %q", tk.typ))
		}

		prevType = tk.typ

	}

	return out, nil
}
