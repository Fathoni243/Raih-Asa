package user

import (
	"net/http"
	"raih-asa/auth"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func InitRouter(db *gorm.DB, r *gin.Engine) {

	r.POST("/user/register", func(c *gin.Context) {
		var body PostRegisterBody
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
				"foto":       user.Foto,
				"pengalaman": user.Pengalaman,
				"skill":      user.Skill,
				"deskripsi":  user.Deskripsi,
			},
		})
	})

	r.GET("/user", func(c *gin.Context) {
		var body []DisplayUserBody

		if result := db.Model(&User{}).Find(&body); result.Error != nil {
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
			"data":    body,
		})
	})

	r.POST("/user/login", func(c *gin.Context) {
		var body PostLoginBody
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

	r.GET("/user/token", auth.AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var body []DisplayUserBody

		if result := db.Model(&User{}).Where("id = ?", id).Take(&body); result.Error != nil {
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
			"data":    body,
		})
	})

	r.PATCH("/user/update/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body PatchUserBody
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
			"data": gin.H{
				"ID":         user.ID,
				"nama":       user.Name,
				"Email":      user.Email,
				"Foto":       user.Foto,
				"Pengalaman": user.Pengalaman,
				"Skill":      user.Skill,
				"Deskripsi":  user.Deskripsi,
			},
		})
	})

	r.PATCH("/user/update/token", auth.AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var body PatchUserBody
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
			ID:    uint(id.(float64)),
			Name:  body.Name,
			Email: body.Email,
			// Password:   hash,
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
			"data": gin.H{
				"ID":         user.ID,
				"nama":       user.Name,
				"Email":      user.Email,
				"Foto":       user.Foto,
				"Pengalaman": user.Pengalaman,
				"Skill":      user.Skill,
				"Deskripsi":  user.Deskripsi,
			},
		})
	})

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