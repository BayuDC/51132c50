package teacher

type CreateTeacherSchema struct {
	Username string `json:"username" binding:"required,alphanum"`
	Fullname string `json:"fullname" binding:"required"`
}
type UpdateTeacherSchema struct {
	Fullname *string `json:"fullname"`
}
