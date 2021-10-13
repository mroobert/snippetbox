package forms

// errors type is use to hold the validation error
// messages for forms. The name of the form field will be used as the key in
// this map.
type errors map[string][]string

// Add error messages for a given field to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Retrieve the first error message for a given
// field from the map.
func (e errors) Get(field string) string {
	errorMessages := e[field]
	if len(errorMessages) == 0 {
		return ""
	}
	return errorMessages[0]
}
