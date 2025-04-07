package models

type TargetDataBase struct {
	Id       uint   `json:"id"`
	Host     string `json:"host" binding:"required"`
	Port     uint   `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Dialect  string `json:"dialect" default:"33"`
}
