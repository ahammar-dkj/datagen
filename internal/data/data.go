package data

import (
	"github.com/bxcodec/faker"
)

type Login struct {
	Id string `faker:"uuid_hyphenated"`
	LoginTime int64 `faker:"unix_time"`
	UserName string `faker:"username"`
	Email string `faker:"email"`
	Position Position
}

type Position struct {
	Latitude float32 `faker:"lat"`
	Longitude float32 `faker:"long"`
}

func NewLogin() (Login, error) {
	l := Login{}
	err := faker.FakeData(&l)
	if err != nil {
		return Login{}, err
	}
	return l, nil
}
