package apis

import (
	"reflect"
)

type User struct {
	User_id     int    `json:"user_id"`
	User_name   string `json:"user_name"`
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
	Email       string `json:"email"`
	User_status string `json:"user_status"`
	Department  string `json:"department"`
}

func CheckNullFields(user User) map[string]interface{} {
	nonNullFields := make(map[string]interface{})

	val := reflect.ValueOf(user)
	typ := reflect.TypeOf(user)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Tag.Get("json")

		if !isZeroValue(field) {
			nonNullFields[fieldName] = field.Interface()
		}
	}
	return nonNullFields
}

func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
