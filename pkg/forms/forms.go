package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Embeds a url.Values object (to hold the form data)
// and an Errors field to hold any validation errors
type Form struct {
	url.Values
	Errors errors
}

// Initialize a custom Form struct
func New(data url.Values) *Form {
	return &Form {
		data,
		errors(map[string][]string{}),
	}
}

// Check that specific fields in the form data are present and not blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Check that a specific field in the form contains a maximum number of characters
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d)", d))
	}
}

// Check that a specific field in the form matches one of the set of specific permitted values
func (f *Form) PermittedValues(field string, opts...string) {
	value := f.Get(field)
	if field == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}