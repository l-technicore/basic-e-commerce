package order

type output struct {
	code int
	message string
	data interface
}

type order struct {
	id int
	details string
	value int
}