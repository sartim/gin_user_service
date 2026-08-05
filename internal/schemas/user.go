package schemas

type UserSchema struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
	Email     string `json:"email" binding:"required,email,max=254"`
	Password  string `json:"password" binding:"required,min=12,max=72"`
	IsActive  bool   `json:"is_active"`
	IsAdmin   bool   `json:"is_admin"`
}
