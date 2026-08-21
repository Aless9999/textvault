package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// возвращает да если нет ошибок
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)

}

// добавляет ошибку в струтуру
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exites := v.FieldErrors[key]; !exites {

		v.FieldErrors[key] = message
	}
}

// добавляет данные в структуру если валидация не пройдена
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// если строка не пустая вернет true
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

//возвращает true если символов не больше n

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) < n
}

// возвращает true если значение находится в списке разрешенных
func PermittedValue[T comparable](value T, permittedValue ...T) bool {
	return slices.Contains(permittedValue, value)
}
