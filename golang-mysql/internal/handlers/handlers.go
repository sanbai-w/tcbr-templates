package handlers

import (
	"net/http"

	"mysql-ip-connect/golang/internal/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	DB *db.Database
}

type CountRequest struct {
	Action string `json:"action"`
}

func NewServer(database *db.Database) *Server {
	return &Server{DB: database}
}

func (s *Server) GetIndex(c *gin.Context) {
	c.File("web/index.html")
}

// GetStatus 新增：获取数据库状态接口
func (s *Server) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"connected": s.DB.IsConnected(),
		},
	})
}

func (s *Server) GetCount(c *gin.Context) {
	if !s.DB.IsConnected() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    -1,
			"message": "数据库未连接，请检查配置",
			"data":    nil,
		})
		return
	}

	cnt, err := s.DB.CountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "数据库操作失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cnt})
}

func (s *Server) PostCount(c *gin.Context) {
	if !s.DB.IsConnected() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    -1,
			"message": "数据库未连接，请检查配置",
			"data":    nil,
		})
		return
	}

	var req CountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的请求",
			"data":    nil,
		})
		return
	}

	switch req.Action {
	case "inc":
		if err := s.DB.InsertOne(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "数据库操作失败",
				"data":    nil,
			})
			return
		}
	case "clear":
		if err := s.DB.Truncate(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "数据库操作失败",
				"data":    nil,
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的操作",
			"data":    nil,
		})
		return
	}

	cnt, err := s.DB.CountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "数据库操作失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cnt})
}

func (s *Server) WxOpenID(c *gin.Context) {
	if c.GetHeader("x-wx-source") != "" {
		c.String(http.StatusOK, c.GetHeader("x-wx-openid"))
		return
	}
	c.Status(http.StatusOK)
}
