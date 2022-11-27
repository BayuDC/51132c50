package auth

type LoginSchema struct {
	Username string `json:"username" bind:"required"`
	Password string `json:"password" bind:"required"`
}
