package user

type output struct {
	code int
	message string
	data interface
}

type user struct {
	id int
	name string
	mobile int
	address string
}

type order struct {
	id int
	product string
	quantity int
	value int
}