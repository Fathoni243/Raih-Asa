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
	Comment       []Comment       `gorm:"foreignKey:Lomba_ID"`
}

type CategoryLomba struct {
	ID            uint    `gorm : "primarykey"`
	Name_Category string  `json : "name_category"`
	Lomba         []Lomba `gorm:"many2many:lomba_category;"`
}

type Comment struct {
	ID       uint   `gorm : "primarykey"`
	UserID   uint   `json:"user_id"`
	Contents string `json:"contents"`
	Lomba_ID uint64
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

type PostCommentBody struct {
	Contents string `json:"contents"`
}
