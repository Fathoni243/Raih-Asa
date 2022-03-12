package beasiswa

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, r *gin.Engine) {

	r.POST("/beasiswa", func(c *gin.Context) {
		var body PostBeasiswaBody
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

		if result := db.Create(&beasiswa).Preload("category_beasiswa"); result.Error != nil {
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
			"data":    beasiswa,
		})
	})

	r.POST("/beasiswa/category", func(c *gin.Context) {
		var body PostCategoryBody
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
			trx = trx.Model(&Beasiswa{}).Preload("Category").Where("judul LIKE ?", "%"+judul+"%")
		}
		if isPenyelenggaraExists {
			trx = trx.Model(&Beasiswa{}).Preload("Category").Where("penyelenggara LIKE ?", "%"+penyelenggara+"%")
		}
		if isDeskripsiExists {
			trx = trx.Model(&Beasiswa{}).Preload("Category").Where("deskripsi LIKE ?", "%"+deskripsi+"%")
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
		trx := db
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
		queryBeasiswa := []Beasiswa{}

		if result := db.Preload("Beasiswa").Where("id = ?", queryCategory.ID).Take(&queryCategory); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		var listArrId = []uint{}
		for i := 0; i < len(queryCategory.Beasiswa); i++ {
			add := append(listArrId, queryCategory.Beasiswa[i].ID)
			listArrId = add
		}

		trx = trx.Model(&Beasiswa{}).Preload("Category").Where("id IN ?", listArrId).Find(&queryBeasiswa)

		trx.Model(&Beasiswa{}).Preload("Category").Find(&queryBeasiswa)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query success.",
			"data": gin.H{
				"query": gin.H{
					"ID":            queryCategory.ID,
					"Name_Category": queryCategory.Name_Category,
				},
				"result": queryBeasiswa,
			},
		})

	})

	r.GET("/beasiswa/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		beasiswa := Beasiswa{}

		if result := db.Preload("Category").Where("id = ?", id).Take(&beasiswa); result.Error != nil {
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
			"data":    beasiswa,
		})

	})
}
