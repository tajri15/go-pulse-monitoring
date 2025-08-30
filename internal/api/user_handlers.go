package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type registerUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// Nanti di sini kita bisa cek error spesifik, misal email sudah ada
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Token string   `json:"token"`
	User  db.User `json:"user"`
}


func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Jika user tidak ditemukan, beri pesan unauthorized
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Jika password salah
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	
	// Jika berhasil, generate token JWT
	// (Kita akan buat fungsi generateToken di langkah selanjutnya)
	token, err := generateToken(user.ID, 3 * time.Hour) // Token valid selama 3 jam
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
	
	rsp := loginUserResponse{
		Token: token,
		User: user,
	}

	ctx.JSON(http.StatusOK, rsp)
}