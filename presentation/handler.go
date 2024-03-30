package presentation

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"plum/domain"
	"plum/usecase"
)

type Handler struct {
	contactService usecase.ContactService
}

func NewHandler(contactService usecase.ContactService) Handler {
	return Handler{contactService: contactService}
}

func (h *Handler) SupportContact(c *gin.Context) {
	var contact domain.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := contact.Validation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// Hubspotから顧客情報を取得

	go func() {
		if err := h.contactService.RespondContact(contact); err != nil {
			log.Printf("%v", err)
		}
	}()

	c.JSON(200, "hello!!")
}