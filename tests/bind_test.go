package tests

import (
	"net/http"
	"testing"

	"gowins/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

func TestBind(t *testing.T) {
	router := gin.New()
	router.GET("/getc", func(c *gin.Context) {
		var b StructC
		err := c.Bind(&b)

		assert.Nil(t, err)
		c.JSON(200, gin.H{
			"a": b.NestedStructPointer,
			"c": b.FieldC,
		})
	})

	rr := utils.PerformRequest(router, http.MethodGet, "/getc?field_a=hello&field_c=world")
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"a":{"FieldA":"hello"},"c":"world"}`, rr.Body.String())
}
