package lomba

type Lomba struct {
	ID            uint            `gorm : "primarykey"`
	Judul         string          `json:"judul"`
	Penyelenggara string          `json:"penyelenggara"`
	Deskripsi     string          `json:"deskripsi"`
	Poster        string          `json:"poster"`
	TanggalDaftar string          `json:"tanggaldaftar"`
	TanggalAkhir  string          `json:"tanggalakhir"`
	Syarat        string          `json:"syarat"`
	CP            string          `json:"cp"`
	Link          string          `json:"link"`
	Category      []CategoryLomba `gorm:"many2many:lomba_category;"`
}

type CategoryLomba struct {
	ID            uint    `gorm : "primarykey"`
	Name_Category string  `json : "name_category"`
	Lomba         []Lomba `gorm:"many2many:lomba_category;"`
}

type PostLombaBody struct {
	Judul         string `json:"judul"`
	Penyelenggara string `json:"penyelenggara"`
	Deskripsi     string `json:"deskripsi"`
	Poster        string `json:"poster"`
	TanggalDaftar string `json:"tanggaldaftar"`
	TanggalAkhir  string `json:"tanggalakhir"`
	Syarat        string `json:"syarat"`
	CP            string `json:"cp"`
	Link          string `json:"link"`
	Category      []CategoryLomba
}

// type postCategoryBody struct {
// 	ID            uint   `gorm : "primarykey"`
// 	Name_Category string `json : "name_category"`
// }