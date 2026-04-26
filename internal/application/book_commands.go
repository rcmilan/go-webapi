package application

type RegisterBookCommand struct {
	Title       string
	ISBN        string
	Price       float64
	ReleaseYear int
}

type BookFilter struct {
	ID   uint32
	ISBN string
}

type BookResult struct {
	ID          uint32
	Title       string
	ISBN        string
	Price       float64
	ReleaseYear int
}
