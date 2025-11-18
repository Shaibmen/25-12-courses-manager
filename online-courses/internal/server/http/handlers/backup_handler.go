package handlers

import (
	"context"
	"net/http"
	"online-courses/internal/domain/service"
	"online-courses/internal/server/http/models"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	handler service.BackupService
}

func NewBackupHandler(handler service.BackupService) *BackupHandler {
	return &BackupHandler{handler: handler}
}

func (b *BackupHandler) BackupDB(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := b.handler.BackupDB(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "бекап базы данных успешно создан"})
}

func (b *BackupHandler) RestoreDB(c *gin.Context) {

	nameFile := c.Query("file")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := b.handler.RestoreDB(ctx, nameFile)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "база данных успешно восстановлена"})
}

func (b *BackupHandler) AllBackup(c *gin.Context) {
	data, err := b.handler.AllBackup()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models.HttpResponseWithData{Data: data})
}
