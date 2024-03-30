package domain

import "errors"

var ContentsIsEmptyError = errors.New("contents is empty")
var EmailIsEmptyError = errors.New("email is empty")
var InvalidEmailError = errors.New("invalid email")
