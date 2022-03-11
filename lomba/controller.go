package lomba

import (
	"database/sql"
	"net/http"
	"raih-asa/auth"
	"raih-asa/beasiswa"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
func sendMail(to []string, subject, message string) error {
	return nil
}*/

func InitRouter(db *gorm.DB, r *gin.Engine) {

	r.POST("/lomba", func(c *gin.Context) {
		var body PostLombaBody
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

	r.POST("/lomba/category", func(c *gin.Context) {
		var body beasiswa.PostCategoryBody
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

	r.POST("/lomba/comment/:id_lomba", auth.AuthMiddleware(), func(c *gin.Context) {
		idLomba, _ := c.Params.Get("id_lomba")
		id, _ := c.Get("id")

		var body PostCommentBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}

		replied_to1, _ := c.GetQuery("replied_to")

		if err := db.Where("id = ?", idLomba).Take(&Lomba{}); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal mencari id lomba.",
				"error":   err.Error.Error(),
			})
			return
		}

		convReply, _ := strconv.ParseUint(replied_to1, 10, 64)
		repliedTo := sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
		if convReply != 0 {
			repliedTo.Int64 = int64(convReply)
			repliedTo.Valid = true
		}

		conv, _ := strconv.ParseUint(idLomba, 10, 64)

		comment := Comment{
			UserID:     uint(id.(float64)),
			Contents:   body.Contents,
			Lomba_ID:   conv,
			Replied_To: repliedTo,
		}

		if result := db.Create(&comment); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Komentar berhasil dibuat.",
			"data":    comment,
		})
	})

	/*
		r.GET("/email/:user_id", auth.AuthMiddleware(), func(c *gin.Context) {
					id, isIdExists := c.Params.Get("id")
					idToken, _ := c.Get("id")
					if !isIdExists {
						c.JSON(http.StatusBadRequest, gin.H{
							"success": false,
							"message": "ID is not supplied.",
						})
						return
					}

					user1 := user.User{
					}

					db.Where("id = ?", uint(idToken.(float64))).Take(&user1)

					nameToken := user1.Name
					emailToken := user1.Email
					passToken := user1.Password

					// fmt.Println(emailToken)
					// fmt.Println(passToken)

					parsedId, _ := strconv.ParseUint(id, 10, 64)

					user2 := user.User{
						ID: uint(parsedId),
					}

					if err := db.Find(&user2); err.Error != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"success": false,
							"message": "Id tidak ditemukan.",
							"error":   err.Error.Error(),
						})
						return
					}
					email := user2.Email
					// fmt.Println(email)

					var body PostEmailBody

					const CONFIG_SMTP_HOST = "smtp.gmail.com"
					const CONFIG_SMTP_PORT = 587
					const CONFIG_SENDER_NAME = nameToken
					const CONFIG_AUTH_EMAIL = ""+emailToken
					const CONFIG_AUTH_PASSWORD = ""+passToken

					to := []string{email}
					subject := ""+body.Subject
					message := ""+body.Message

					bodyEmail := "From: " + CONFIG_SENDER_NAME + "\n" +
			        	"To: " + strings.Join(to, ",") + "\n" +
			        	"Subject: " + subject + "\n\n" +
			        	message

			    	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)
			    	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

			   		err := smtp.SendMail(smtpAddr, auth, CONFIG_AUTH_EMAIL, append(to), []byte(bodyEmail))
			    	if err != nil {
			        	return
			    	}

					log.Println("Mail sent!")
		})
	*/

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

		if result := db.Model(&Lomba{}).Preload("Lomba").Preload("Category").Preload("Comment").Take(&queryCategory); result.Error != nil {
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

	r.GET("/lomba/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		lomba := Lomba{}

		if result := db.Preload("Category").Preload("Comment").Where("id = ?", id).Take(&lomba); result.Error != nil {
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
			"data":    lomba,
		})

	})

}
