package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var r *gin.Engine

type User struct {
	ID         uint   `gorm : "primarykey"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Foto       string `json:"foto"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}

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

type CekPassword struct {
	Name_User    string `json:"name_user"`
	Cek_Password string `json:"cek_password"`
}

type CategoryBeasiswa struct {
	ID            uint       `gorm : "primarykey"`
	Name_Category string     `json : "name_category"`
	Beasiswa      []Beasiswa `gorm:"many2many:beasiswa_category;"`
}

type CategoryLomba struct {
	ID            uint    `gorm : "primarykey"`
	Name_Category string  `json : "name_category"`
	Lomba         []Lomba `gorm:"many2many:lomba_category;"`
}

type postRegisterBody struct {
	Name       string `json: "name"`
	Email      string `json: "email"`
	Password   string `json: "password"`
	Foto       string `json:"foto"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}

type postLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postBeasiswaBody struct {
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

type postLombaBody struct {
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

type postCategoryBody struct {
	ID            uint   `gorm : "primarykey"`
	Name_Category string `json : "name_category"`
}

type patchUserBody struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Foto       string `json:"foto"`
	Pengalaman string `json:"pengalaman"`
	Skill      string `json:"skill"`
	Deskripsi  string `json:"deskripsi"`
}

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/raih_asa?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	if err = db.AutoMigrate(&User{}, &Beasiswa{}, &CategoryBeasiswa{},
		&Lomba{}, &CategoryLomba{}, &CekPassword{}); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte("passwordBuatSigning"), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
	}
}

func InitRouter() {
	//CREATE DATA
	r.POST("/user/register", func(c *gin.Context) {
		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}

		hash, _ := HashPassword(body.Password)

		user := User{
			Name:       body.Name,
			Email:      body.Email,
			Password:   hash,
			Foto:       body.Foto,
			Pengalaman: body.Pengalaman,
			Skill:      body.Skill,
			Deskripsi:  body.Deskripsi,
		}
		//pengecekan isi password
		cekPass := CekPassword{
			Name_User:    body.Name,
			Cek_Password: body.Password,
		}

		var cek = []byte(body.Password)
		var angka = false

		for i := 0; i < len(body.Password); i++ {
			if cek[i] >= 48 && cek[i] <= 57 {
				angka = true
			}
		}

		if (len(body.Password) >= 8) && (cek[0] >= 65 && cek[0] <= 90) && angka == true {
			result := db.Create(&user)
			db.Create(&cekPass) //create Cek Password
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when inserting into the database.",
					"error":   result.Error.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Password harus lebih dari 8 karakter, Huruf pertama harus kapital, Password harus terdapat angka",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Akun berhasil dibuat.",
			"data": gin.H{
				"nama":  user.Name,
				"email": user.Email,
			},
		})
	})

	r.POST("/beasiswa", func(c *gin.Context) {
		var body postBeasiswaBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		beasiswa := Beasiswa{
			Judul:         body.Judul,
			Penyelenggara: body.Penyelenggara,
			Deskripsi:     body.Deskripsi,
			Poster:        body.Poster,
			TanggalDaftar: body.TanggalDaftar,
			TanggalAkhir:  body.TanggalAkhir,
			Syarat:        body.Syarat,
			CP:            body.CP,
			Link:          body.Link,
			Category:      body.Category,
		}

		if result := db.Create(&beasiswa); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Beasiswa berhasil dibuat.",
			"data": gin.H{
				"judul":          beasiswa.Judul,
				"penyelenggara":  beasiswa.Penyelenggara,
				"deskripsi":      beasiswa.Deskripsi,
				"poster":         beasiswa.Poster,
				"tanggal daftar": beasiswa.TanggalDaftar,
				"tanggal akhir":  beasiswa.TanggalAkhir,
				"syarat":         beasiswa.Syarat,
				"cp":             beasiswa.CP,
				"link":           beasiswa.Link,
				"Category":       beasiswa.Category,
			},
		})
	})

	r.POST("/category/beasiswa", func(c *gin.Context) {
		var body postCategoryBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		category := CategoryBeasiswa{
			Name_Category: body.Name_Category,
		}

		if result := db.Create(&category); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Kategori berhasil dibuat.",
			"data": gin.H{
				"nama kategori": category.Name_Category,
			},
		})
	})

	r.POST("/lomba", func(c *gin.Context) {
		var body postLombaBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		lomba := Lomba{
			Judul:         body.Judul,
			Penyelenggara: body.Penyelenggara,
			Deskripsi:     body.Deskripsi,
			Poster:        body.Poster,
			TanggalDaftar: body.TanggalDaftar,
			TanggalAkhir:  body.TanggalAkhir,
			Syarat:        body.Syarat,
			CP:            body.CP,
			Link:          body.Link,
			Category:      body.Category,
		}

		if result := db.Create(&lomba); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Lomba berhasil dibuat.",
			"data": gin.H{
				"judul":          lomba.Judul,
				"penyelenggara":  lomba.Penyelenggara,
				"deskripsi":      lomba.Deskripsi,
				"poster":         lomba.Poster,
				"tanggal daftar": lomba.TanggalDaftar,
				"tanggal akhir":  lomba.TanggalAkhir,
				"syarat":         lomba.Syarat,
				"cp":             lomba.CP,
				"link":           lomba.Link,
				"Category":       lomba.Category,
			},
		})
	})

	r.POST("/category/lomba", func(c *gin.Context) {
		var body postCategoryBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		category := CategoryLomba{
			Name_Category: body.Name_Category,
		}

		if result := db.Create(&category); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Kategori berhasil dibuat.",
			"data": gin.H{
				"nama kategori": category.Name_Category,
			},
		})
	})

	//READ
	r.GET("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		user := User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query success.",
			"data": gin.H{
				"ID":         user.ID,
				"nama":       user.Name,
				"email":      user.Email,
				"password":   user.Password,
				"foto":       user.Foto,
				"pengalaman": user.Pengalaman,
				"skill":      user.Skill,
				"deskripsi":  user.Deskripsi,
			},
		})
	})

	r.GET("/user", func(c *gin.Context) {
		var allUsersFromDB []User

		if result := db.Find(&allUsersFromDB); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    allUsersFromDB,
		})
	})

	r.POST("/user/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		if result := db.Where("email = ?", body.Email).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Email tidak ditemukan.",
				"error":   result.Error.Error(),
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
			})
			tokenString, err := token.SignedString([]byte("passwordBuatSigning"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when generating the token.",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login Berhasil.",
				"data": gin.H{
					"name":  user.Name,
					"email": user.Email,
					"token": tokenString,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Password salah.",
			})
			return
		}
	})

	r.GET("/user/token", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query success",
			"data":    user,
		})
	})

	r.GET("/beasiswa/search", func(c *gin.Context) {
		var queryResults []Beasiswa
		trx := db

		judul, isJudulExists := c.GetQuery("judul")
		penyelenggara, isPenyelenggaraExists := c.GetQuery("penyelenggara")
		deskripsi, isDeskripsiExists := c.GetQuery("deskripsi")

		if !isJudulExists && !isPenyelenggaraExists && !isDeskripsiExists {
			if result := trx.Find(&queryResults); result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Query is not supplied.",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"data": queryResults,
			})
			return
		}

		if isJudulExists {
			trx = trx.Where("judul LIKE ?", "%"+judul+"%")
		}
		if isPenyelenggaraExists {
			trx = trx.Where("penyelenggara LIKE ?", "%"+penyelenggara+"%")
		}
		if isDeskripsiExists {
			trx = trx.Where("deskripsi LIKE ?", "%"+deskripsi+"%")
		}

		result := trx.Find(&queryResults)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search success",
			"data": gin.H{
				"query": gin.H{
					"judul":         judul,
					"penyelenggara": penyelenggara,
					"deskripsi":     deskripsi,
				},
				"result": queryResults,
			},
		})
	})

	r.GET("/beasiswa/category/:category_beasiswa_id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("category_beasiswa_id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}

		queryCategory := CategoryBeasiswa{
			ID: uint(parsedId),
		}
		if result := db.Preload("Beasiswa").Take(&queryCategory); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query success.",
			"data":    queryCategory,
		})

	})

	r.GET("/lomba/search", func(c *gin.Context) {
		var queryResults []Lomba
		trx := db

		judul, isJudulExists := c.GetQuery("judul")
		penyelenggara, isPenyelenggaraExists := c.GetQuery("penyelenggara")
		deskripsi, isDeskripsiExists := c.GetQuery("deskripsi")

		if !isJudulExists && !isPenyelenggaraExists && !isDeskripsiExists {
			if result := trx.Find(&queryResults); result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Query is not supplied.",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"data": queryResults,
			})
			return
		}

		if isJudulExists {
			trx = trx.Where("judul LIKE ?", "%"+judul+"%")
		}
		if isPenyelenggaraExists {
			trx = trx.Where("penyelenggara LIKE ?", "%"+penyelenggara+"%")
		}
		if isDeskripsiExists {
			trx = trx.Where("deskripsi LIKE ?", "%"+deskripsi+"%")
		}

		result := trx.Find(&queryResults)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search success",
			"data": gin.H{
				"query": gin.H{
					"judul":         judul,
					"penyelenggara": penyelenggara,
					"deskripsi":     deskripsi,
				},
				"result": queryResults,
			},
		})
	})

	r.GET("/lomba/category/:category_lomba_id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("category_lomba_id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}

		queryCategory := CategoryLomba{
			ID: uint(parsedId),
		}
		if result := db.Preload("Lomba").Take(&queryCategory); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query success.",
			"data":    queryCategory,
		})

	})

	//UPDATE
	r.PATCH("/user/update/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
				"error":   err.Error(),
			})
			return
		}

		user := User{
			ID:         uint(parsedId),
			Name:       body.Name,
			Email:      body.Email,
<<<<<<< HEAD
=======
			Password:   hash,
>>>>>>> 2b74c910d60efe2b933562790c9ace13de7e872b
			Foto:       body.Foto,
			Pengalaman: body.Pengalaman,
			Skill:      body.Skill,
			Deskripsi:  body.Deskripsi,
		}

		result := db.Model(&user).Updates(user)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Success.",
			"data":    user,
		})
	})

	r.PATCH("/user/update/token", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}

		// hash, _ := HashPassword(body.Password)

		user := User{
			ID:         uint(id.(float64)),
			Name:       body.Name,
			Email:      body.Email,
<<<<<<< HEAD
			// Password:   hash,
=======
			Password:   hash,
>>>>>>> 2b74c910d60efe2b933562790c9ace13de7e872b
			Foto:       body.Foto,
			Pengalaman: body.Pengalaman,
			Skill:      body.Skill,
			Deskripsi:  body.Deskripsi,
		}
		result := db.Model(&user).Updates(user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Success.",
			"data":    user,
		})
	})

	//DELETE
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		user := User{
			ID: uint(parsedId),
		}
		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete success.",
		})
	})

}

func InitGin() {
	r = gin.Default()
}

func StartServer() error {
	return r.Run()
}

func main() {
	if err := InitDB(); err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}
	InitGin()
	InitRouter()

	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
