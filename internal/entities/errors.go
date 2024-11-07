package entities

import "errors"

var (
	ErrCacheEmpty             = errors.New("data is not exists")
	ErrEmptyCategories        = errors.New("categories are empty")
	ErrImageNotExists         = errors.New("image is not exist")
	ErrPhoneNumberNotValid    = errors.New("phone-number is not valid")
	ErrEmptyOrder             = errors.New("order is empty")
	ErrMealNotExists          = errors.New("meal is not exist")
	ErrValidateAddressOrFloor = errors.New("error validate address or floor")
	ErrCategoryIdOrPageValid  = errors.New("error when to validate category-id or page")
	ErrPageIsZero             = errors.New("page is equal to zero")
)
