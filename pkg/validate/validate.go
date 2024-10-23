package validate

import "regexp"

type Validator interface{
	ValidatePhoneNumber(phoneNumber string) bool
}

type Validation struct{}

func NewValidation() *Validation {
	return &Validation{}
}

func (v *Validation) ValidatePhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`992\d{9}$`)
	return re.MatchString(phoneNumber)
}
