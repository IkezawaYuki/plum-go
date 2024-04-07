package presentation

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"plum/domain"
	"plum/logger"
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

	// TODO Hubspotから顧客情報を取得

	go func() {
		if err := h.contactService.RespondContact(contact); err != nil {
			logger.Logger.Error("RespondContact is failed", err)
		}
	}()

	c.JSON(http.StatusOK, "success")
}

func (h *Handler) GmailToHubspot(c *gin.Context) {
	var mail domain.Gmail
	if err := c.BindJSON(&mail); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := mail.Validation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	go func() {
		if err := h.contactService.GmailToHubspot(mail); err != nil {
			log.Printf("%v", err)
		}
	}()
	c.JSON(http.StatusOK, "success")
}

func (h *Handler) GmailToAiSearch(c *gin.Context) {
	var mailList domain.GmailList
	if err := c.BindJSON(&mailList); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := mailList.Validation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.contactService.GmailToAiSearch(mailList); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}
