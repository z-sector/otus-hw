package response

import "github.com/gin-gonic/gin"

func WriteError(c *gin.Context, code int, err error) {
	body := gin.H{}
	if err != nil {
		body["error"] = err.Error()
	}
	c.JSON(code, body)
}
