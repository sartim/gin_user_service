package schemas

type PermissionSchema struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"required,max=500"`
}
