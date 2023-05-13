package schemas

type UserPermissionSchema struct {
	UserId       string `json:"user_id" binding:"required"`
	PermissionId string `json:"permission_id" binding:"required"`
}
