package student

type CreateStudentSchema struct {
	Name string `json:"name" binding:"required"`
}
type UpdateStudentSchema struct {
	Name *string `json:"name"`
}
