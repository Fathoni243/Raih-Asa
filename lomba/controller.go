package lomba

import (
	"net/http"
	"raih-asa/auth"
	"raih-asa/beasiswa"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

		if err := db.Where("id = ?", idLomba).Take(&Lomba{}); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal mencari id lomba.",
				"error":   err.Error.Error(),
			})
			return
		}
		conv, _ := strconv.ParseUint(idLomba, 10, 64)

		comment := Comment{
			UserID:   uint(id.(float64)),
			Contents: body.Contents,
			Lomba_ID: conv,
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
}
