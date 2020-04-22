package proto

//Pagination 分页
type Pagination struct {
	Page  int `form:"page" default:"1" json:"page"`
	Limit int `form:"limit" default:"10" json:"limit"`
}
