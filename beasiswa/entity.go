package beasiswa

type Beasiswa struct {
	ID            uint               `gorm : "primarykey"`
	Judul         string             `json:"judul"`
	Penyelenggara string             `json:"penyelenggara"`
	Deskripsi     string             `json:"deskripsi"`
	Poster        string             `json:"poster"`
	TanggalDaftar string             `json:"tanggaldaftar"`
	TanggalAkhir  string             `json:"tanggalakhir"`
	Syarat        string             `json:"syarat"`
	CP            string             `json:"cp"`
	Link          string             `json:"link"`
	Category      []CategoryBeasiswa `gorm:"many2many:beasiswa_category;"`
}

type CategoryBeasiswa struct {
	ID            uint       `gorm : "primarykey"`
	Name_Category string     `json : "name_category"`
	Beasiswa      []Beasiswa `gorm:"many2many:beasiswa_category;"`
}

type PostBeasiswaBody struct {
	Judul         string `json:"judul"`
	Penyelenggara string `json:"penyelenggara"`
	Deskripsi     string `json:"deskripsi"`
	Poster        string `json:"poster"`
	TanggalDaftar string `json:"tanggaldaftar"`
	TanggalAkhir  string `json:"tanggalakhir"`
	Syarat        string `json:"syarat"`
	CP            string `json:"cp"`
	Link          string `json:"link"`
	Category      []CategoryBeasiswa
}

type PostCategoryBody struct {
	ID            uint   `gorm : "primarykey"`
	Name_Category string `json : "name_category"`
}