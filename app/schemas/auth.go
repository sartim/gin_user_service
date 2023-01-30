package schemas

type AuthSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type HeaderSchema struct {
	Authorization *string `header:"Authorization" binding:"required"`
}
