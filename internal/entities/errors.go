package entities

import "errors"

var (
	ErrEmptyBurgers           = errors.New("burgers are empty")
	ErrImageNotExists         = errors.New("image is not exist")
	ErrPhoneNumberNotValid    = errors.New("phone-number is not valid")
	ErrEmptyOrder             = errors.New("order is empty")
	ErrMealNotExists          = errors.New("meal is not exist")
	ErrValidateAddressOrFloor = errors.New("error validate address or floor")
)
