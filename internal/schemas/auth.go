package schemas

type AuthSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type HeaderSchema struct {
	Authorization *string `json:"Authorization" header:"Authorization" binding:"required"`
	ContentType   *string `json:"Content-Type" header:"Content-type" binding:"required"`
}
