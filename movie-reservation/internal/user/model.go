package user

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type InputDTO struct{
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OutputDTO struct{
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginDTO struct{
	Email    string `json:"email"`
	Password string `json:"password"`
}