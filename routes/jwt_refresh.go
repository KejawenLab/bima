package routes

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/middlewares"
	"github.com/KejawenLab/bima/v4/utils"
	"google.golang.org/grpc"
)

type JwtRefresh struct {
	PathUrl       string
	Secret        string
	SigningMethod string
	Expire        int
}

func (j *JwtRefresh) Path() string {
	return j.PathUrl
}

func (j *JwtRefresh) Method() string {
	return http.MethodPost
}

func (j *JwtRefresh) SetClient(client *grpc.ClientConn) {}

func (j *JwtRefresh) Middlewares() []middlewares.Middleware {
	return nil
}

func (j *JwtRefresh) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := context.WithValue(context.Background(), "scope", "jwt_refresh_token")
	body := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
		http.Error(w, "token is empty", http.StatusBadRequest)

		return
	}

	token, ok := body["token"]
	if !ok {
		loggers.Logger.Error(ctx, "token is empty")
		http.Error(w, "token is empty", http.StatusBadRequest)

		return
	}

	claims, err := utils.ValidateRefreshToken(j.Secret, j.SigningMethod, token)
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
		http.Error(w, "invalid token", http.StatusBadRequest)

		return
	}

	token, _ = utils.CreateToken(j.Secret, j.SigningMethod, claims, j.Expire)
	refreshToken, _ := utils.CreateRefreshToken(j.Secret, j.SigningMethod, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "refresh_token": refreshToken})
}
