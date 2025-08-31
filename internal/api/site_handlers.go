package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
)

type createSiteRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func (server *Server) createSite(ctx *gin.Context) {
	var req createSiteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Ambil userID dari context yang sudah di-set oleh middleware
	authPayload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization payload does not exist")))
		return
	}
	userID := authPayload.(int64)

	arg := db.CreateSiteParams{
		UserID: userID,
		URL:    req.URL,
	}

	site, err := server.store.CreateSite(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, site)
}

func (server *Server) listSites(ctx *gin.Context) {
	authPayload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization payload does not exist")))
		return
	}
	userID := authPayload.(int64)

	sites, err := server.store.GetSitesByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sites)
}

func (server *Server) deleteSite(ctx *gin.Context) {
	siteID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid site ID")))
		return
	}

	authPayload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization payload does not exist")))
		return
	}
	userID := authPayload.(int64)

	err = server.store.DeleteSite(ctx, siteID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "site deleted successfully"})
}