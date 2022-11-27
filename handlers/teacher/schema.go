package teacher

type CreateTeacherSchema struct {
	Name string `json:"name" binding:"required"`
}
type UpdateTeacherSchema struct {
	Name *string `json:"name"`
}
