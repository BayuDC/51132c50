package group

type CreateGroupSchema struct {
	Name string `json:"name" binding:"required"`
}
type UpdateGroupSchema struct {
	Name *string `json:"name"`
}
type ManageStudentGroupSchema struct {
	Students []int `json:"students" binding:"required"`
}

type ManageStudentGroupResult struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
