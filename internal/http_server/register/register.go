package register

import "reflect"

type userRegister struct {
	email    string `validate:"mail"`
	password string `validate:"enoughLen"`
}

//todo create own validation for user

func Register() {
	user := reflect.TypeOf(userRegister{})
}
