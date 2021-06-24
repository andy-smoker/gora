package handlers

import (
	"encoding/json"
	"gora/pkg/service"
	"log"

	"github.com/gin-gonic/gin"
)

type ErrResponce struct {
	Err string `json:"err"`
}

type Handler struct {
	*service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{srv}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api") // midleware user verification
	{
		v1 := api.Group("/v1")
		{
			images := v1.Group("/image")
			{
				images.GET("/", h.GetImages)
				images.POST("/", h.AddImage)

				image := images.Group("/:id")
				{
					image.GET("/", h.GetImage)
					image.DELETE("/", h.DeleteImage)
				}
			}
		}
	}

	return router
}

// GetImage
func (h *Handler) GetImage(c *gin.Context) {

	id := c.Param("id")
	data, err := h.ImageServiceInterfase.GetImage(id)
	if err != nil {
		log.Println(err)
		c.JSON(501, &ErrResponce{Err: "Something wrong"})
		return
	}
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		c.JSON(501, &ErrResponce{Err: "Something wrong"})
		return
	}
	c.Data(200, "file/image", resp)
}

func (h *Handler) GetImages(c *gin.Context) {
	resp, err := h.ImageServiceInterfase.GetImageList()
	if err != nil {
		log.Println(err)
		c.JSON(501, &ErrResponce{Err: "Something wrong"})
		return
	}
	c.Data(200, "application/json", resp)
}

func (h *Handler) AddImage(c *gin.Context) {
	resp := struct {
		Result bool `json:"result"`
	}{
		Result: false,
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("file error", err)
		c.JSON(501, resp)
		return
	}
	if file == nil {
		log.Println("no surch file", err)
		c.JSON(501, resp)
	} else {
		defer file.Close()
	}
	resp.Result = h.Service.AddImage(file)
	if resp.Result {
		c.JSON(200, resp)
	}
}
func (h *Handler) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	resp := struct {
		Result bool `json:"result"`
	}{
		Result: false,
	}
	resp.Result = h.ImageServiceInterfase.DeleteImage(id)
	c.JSON(200, resp)
}
