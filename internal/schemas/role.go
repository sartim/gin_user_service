package schemas

type RoleSchema struct {
	Name        string `json:"first_name" binding:"required"`
	Description string `json:"last_name" binding:"required"`
}
