package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

//returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0

}

//New initializes a new Form
func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}

}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this cannot be blank")

		}

	}
}

//Has checks if form field is in request and not empty
func (f *Form) Has(field string, r *http.Request) bool {

	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true

}

// MinLength checks the string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be atleast %d characters long", length))
		return false

	}

	return true

}

// IsEmail checks for email address
func (f *Form) IsEmail(field string) {

	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
	}

}
