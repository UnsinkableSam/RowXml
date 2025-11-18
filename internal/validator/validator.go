package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"RowXml.com/internal/model"
)

// ValidatePeople performs declarative validation over the populated model.
// It returns a combined error of all problems (or nil).
func ValidatePeople(p *model.People) error {
	var problems []string
	// basic regexes
	nameRe := regexp.MustCompile(`^[\p{L} '\-]+$`)
	postRe := regexp.MustCompile(`^\d{2,5}$`)
	phoneRe := regexp.MustCompile(`^\d{2,4}-\d{2,8}$`)
	yearMax := time.Now().Year() + 1

	for i, person := range p.People {
		prefix := fmt.Sprintf("person[%d]", i)
		// names
		if strings.TrimSpace(person.FirstName) == "" {
			problems = append(problems, prefix+".firstname is empty")
		} else if !nameRe.MatchString(person.FirstName) {
			problems = append(problems, prefix+".firstname contains invalid characters")
		}
		if strings.TrimSpace(person.LastName) == "" {
			problems = append(problems, prefix+".lastname is empty")
		} else if !nameRe.MatchString(person.LastName) {
			problems = append(problems, prefix+".lastname contains invalid characters")
		}
		// addresses
		for j, a := range person.Address {
			if strings.TrimSpace(a.Street) == "" {
				problems = append(problems, fmt.Sprintf("%s.address[%d].street empty", prefix, j))
			}
			if strings.TrimSpace(a.City) == "" {
				problems = append(problems, fmt.Sprintf("%s.address[%d].city empty", prefix, j))
			}
			if !postRe.MatchString(a.Postcode) {
				problems = append(problems, fmt.Sprintf("%s.address[%d].postcode invalid", prefix, j))
			}
		}
		// phones
		for j, ph := range person.Phone {
			if ph.Mobile != "" && !phoneRe.MatchString(ph.Mobile) {
				problems = append(problems, fmt.Sprintf("%s.phone[%d].mobile invalid", prefix, j))
			}
			if ph.Landline != "" && !phoneRe.MatchString(ph.Landline) {
				problems = append(problems, fmt.Sprintf("%s.phone[%d].landline invalid", prefix, j))
			}
		}
		// families
		for j, f := range person.Family {
			fpre := fmt.Sprintf("%s.family[%d]", prefix, j)
			if strings.TrimSpace(f.Name) == "" {
				problems = append(problems, fpre+".name empty")
			} else if !nameRe.MatchString(f.Name) {
				problems = append(problems, fpre+".name invalid")
			}
			if f.Born < 1800 || f.Born > yearMax {
				problems = append(problems, fmt.Sprintf("%s.born out of range: %d", fpre, f.Born))
			}
			for k, a := range f.Address {
				if strings.TrimSpace(a.Street) == "" {
					problems = append(problems, fmt.Sprintf("%s.address[%d].street empty", fpre, k))
				}
				if strings.TrimSpace(a.City) == "" {
					problems = append(problems, fmt.Sprintf("%s.address[%d].city empty", fpre, k))
				}
				if !postRe.MatchString(a.Postcode) {
					problems = append(problems, fmt.Sprintf("%s.address[%d].postcode invalid", fpre, k))
				}
			}
			for k, ph := range f.Phone {
				if ph.Mobile != "" && !phoneRe.MatchString(ph.Mobile) {
					problems = append(problems, fmt.Sprintf("%s.phone[%d].mobile invalid", fpre, k))
				}
				if ph.Landline != "" && !phoneRe.MatchString(ph.Landline) {
					problems = append(problems, fmt.Sprintf("%s.phone[%d].landline invalid", fpre, k))
				}
			}
		}
	}

	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "\n"))
	}
	return nil
}

