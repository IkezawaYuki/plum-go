package presentation

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func (h *Handler) SupportForm(c *gin.Context) {
	form := domain.NewForm(
		c.PostForm("company"),
		c.PostForm("phone"),
		c.PostForm("lastname"),
		c.PostForm("firstname"),
		c.PostForm("email"),
		c.PostForm("content"),
	)
	fmt.Println(form)
	if err := form.Validation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	go func() {
		if err := h.contactService.RespondContact(&form); err != nil {
			logger.Logger.Error("RespondContact is failed", err)
		}
	}()
	c.Redirect(http.StatusSeeOther, "/plum/support/thank_you")
}

func (h *Handler) SupportMail(c *gin.Context) {
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
		if err := h.contactService.RespondContact(&mail); err != nil {
			logger.Logger.Error("RespondContact is failed", err)
		}
	}()
	c.JSON(http.StatusOK, "success")
}

func (h *Handler) SupportFormPage(c *gin.Context) {
	c.HTML(http.StatusOK, "form.tmpl", gin.H{})
}

func (h *Handler) ThankYouPage(c *gin.Context) {
	c.HTML(http.StatusOK, "thank_you.tmpl", gin.H{})
}

func (h *Handler) DashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{})
}
