package apperror

import "errors"

var (
	ErrNotFound          = errors.New("Не найден")
	ErrAlreadyRegistered = errors.New("Пользователь с таким email уже зарегистрирован")
)
