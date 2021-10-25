package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form struct anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors errors
}

// Initialize a custom Form struct. Notice that
// this takes the form data as the parameter?
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Checks that specific fields in the form data are present and not blank.
// If any fields fail this check, add the appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Checks that a specific field in the form contains a maximum number of characters.
// If the check fails then add the appropriate message to the form errors.
func (f *Form) MaxLength(field string, max int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > max {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", max))
	}
}

// Checks that a specific field in the form contains a minimum number of characters.
// If the check fails then add the appropriate message to the form errors.
func (f *Form) MinLength(field string, min int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < min {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", min))
	}
}

// Checks that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Returns true if the email address is valid conform RFC 5322 and extended by RFC 6532
func (f *Form) ParseEmail(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		f.Errors.Add(field, "Email address is not valid")
	}
}
