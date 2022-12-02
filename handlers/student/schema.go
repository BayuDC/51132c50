package student

type CreateStudentSchema struct {
	Username string `json:"username" binding:"required,alphanum"`
	Fullname string `json:"fullname" binding:"required"`
}
type UpdateStudentSchema struct {
	Fullname *string `json:"fullname"`
}
