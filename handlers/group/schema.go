package group

type CreateGroupSchema struct {
	Name string `json:"name" binding:"required"`
}
type UpdateGroupSchema struct {
	Name *string `json:"name"`
}
