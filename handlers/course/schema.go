package course

type CreateCourseSchema struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Teacher     *int   `json:"teacher"`
}
type UpdateCourseSchema struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Teacher     *int    `json:"teacher"`
}

type ManageCourseMemberSchema struct {
	Students []int `json:"students"`
}

type ManageCourseMemberResult struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type CreateAssignmentSchema struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=empty files"`
}
