package course

type CreateCourseSchema struct {
	Name    string `json:"name" binding:"required"`
	Teacher *int   `json:"teacher"`
}
type UpdateCourseSchema struct {
	Name *string `json:"name"`
}
