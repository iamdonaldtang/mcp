package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Data: data, Message: "success"})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: 0, Data: data, Message: "created"})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Code: status, Data: nil, Message: message})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context) {
	Error(c, http.StatusUnauthorized, "unauthorized")
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "internal server error")
}

func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: gin.H{
			"items":    items,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
		Message: "success",
	})
}
