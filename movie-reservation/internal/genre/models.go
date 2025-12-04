package genre

type Genre struct {
	ID   int64
	Name string
}

type InputDTO struct {
	Name string `json:"name"`
}

type OutputDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
