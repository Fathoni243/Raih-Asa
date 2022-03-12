package user

import (
	"raih-asa/lomba"
)

type User struct {
	ID         uint            `gorm : "primarykey"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"password"`
	Foto       string          `json:"foto"`
	Pengalaman string          `json:"pengalaman"`
	Skill      string          `json:"skill"`
	Deskripsi  string          `json:"deskripsi"`
	Comment    []lomba.Comment `gorm:"foreignKey:UserID"`
}

type CekPassword struct {
	Name_User    string `json:"name_user"`
	Cek_Password string `json:"cek_password"`
}

type PostRegisterBody struct {
	Name       string `json: "name"`
	Email      string `json: "email"`
	Password   string `json: "password"`
	Foto       string `json:"foto"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}

type PostLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PatchUserBody struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}

type DisplayUserBody struct {
	ID         uint
	Name       string `json:"name"`
	Email      string `json:"email"`
	Foto       string `json:"foto"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}
