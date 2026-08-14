package types

type Student struct{
	Id int
	name string `validate:"required"`
	Email string `validate:"required"`
	Age int `validate:"required"`
}