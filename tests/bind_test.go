package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"gowins/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

func TestBind(t *testing.T) {
	router := gin.New()
	router.GET("/getc", func(c *gin.Context) {
		var bc StructC
		err := c.ShouldBind(&bc)
		var bc2 StructC
		_ = c.ShouldBind(&bc2)
		fmt.Println(bc, bc2)
		assert.Nil(t, err)
		c.JSON(200, gin.H{
			"a": bc.NestedStructPointer,
			"c": bc.FieldC,
		})
	})

	rr := utils.PerformRequest(router, http.MethodGet, "/getc?field_a=hello&field_c=world", nil, nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"a":{"FieldA":"hello"},"c":"world"}`, rr.Body.String())
}

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Person struct {
	Name string `form:"user"`
}

func TestBindJSON(t *testing.T) {
	router := gin.New()
	router.POST("/post", func(c *gin.Context) {
		var login Login

		err := c.ShouldBind(&login)

		assert.Nil(t, err)
		assert.Equal(t, "kenny", login.User)
		assert.Equal(t, "1234", login.Password)

		var person Person
		errQ := c.ShouldBindQuery(&person)
		assert.Nil(t, errQ)
		assert.Equal(t, "appleboy", person.Name)
	})

	rr := utils.PerformRequest(router, http.MethodPost, "/post?user=appleboy",
		[]utils.Header{{Key: "Content-Type", Value: "application/json"}},
		strings.NewReader(`{"user":"kenny","password":"1234"}`))

	assert.Equal(t, http.StatusOK, rr.Code)
}
