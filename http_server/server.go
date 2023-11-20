package http_server

import (
	"awesomeProject/cache"
	"awesomeProject/repository"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"net/http"
)

type Handler struct {
	Storage *cache.Cache
	Sc      stan.Conn
}

func (h *Handler) homePageGET(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) homePagePOST(c *gin.Context) {
	t := c.PostForm("message")
	c.Redirect(http.StatusFound, "/orders?order_id="+t)
}

func (h *Handler) getOrder(c *gin.Context) {
	orderId := c.Query("order_id")

	dt, found := h.Storage.Get(orderId)
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, dt)
}

func (h *Handler) publishOrder(c *gin.Context) {
	var dt repository.Order

	if err := c.BindJSON(&dt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	validate := validator.New()

	if err := validate.Struct(dt); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	jsonData, _ := json.Marshal(dt)
	if err := h.Sc.Publish("foo", jsonData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusCreated, dt)
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", h.homePageGET)
	router.POST("/", h.homePagePOST)
	router.GET("/orders", h.getOrder)
	router.POST("/publish", h.publishOrder)
	return router
}
