package routes

import (
	"context"
	"net/http"

	"github.com/KejawenLab/bima/v4/loggers"
	"github.com/KejawenLab/bima/v4/middlewares"
	"github.com/KejawenLab/bima/v4/utils"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
)

type (
	JwtLogin struct {
		PathUrl       string
		UserField     string
		PasswordField string
		Secret        string
		SigningMethod string
		Expire        int
		RefreshToken  bool
		FindUser      FindUserByUsernameAndPassword
	}

	FindUserByUsernameAndPassword func(username string, password string) jwt.MapClaims
)

func NewJwtLogin(
	path string,
	userField string,
	passwordField string,
	secret string,
	signingMethod string,
	expire int,
	refreshToken bool,
	findUser FindUserByUsernameAndPassword,
) *JwtLogin {
	return &JwtLogin{
		PathUrl:       path,
		Secret:        secret,
		SigningMethod: signingMethod,
		FindUser:      findUser,
		UserField:     userField,
		PasswordField: passwordField,
		Expire:        expire,
		RefreshToken:  refreshToken,
	}
}

func DefaultJwtLogin(path string, secret string, signingMethod string, refreshToken bool, findUser FindUserByUsernameAndPassword) *JwtLogin {
	return NewJwtLogin(path, "username", "password", secret, signingMethod, 2, refreshToken, findUser)
}

func (j *JwtLogin) Path() string {
	return j.PathUrl
}

func (j *JwtLogin) Method() string {
	return http.MethodPost
}

func (j *JwtLogin) SetClient(client *grpc.ClientConn) {}

func (j *JwtLogin) Middlewares() []middlewares.Middleware {
	return nil
}

func (j *JwtLogin) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := context.WithValue(context.Background(), "scope", "jwt_login")
	user := map[string]string{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		loggers.Logger.Error(ctx, err.Error())
		http.Error(w, "username and/or password not match", http.StatusBadRequest)

		return
	}

	var username string
	var password string
	var ok bool
	if username, ok = user[j.UserField]; !ok {
		loggers.Logger.Error(ctx, "username or password not provided")
		http.Error(w, "username and/or password not match", http.StatusBadRequest)

		return
	}

	if password, ok = user[j.PasswordField]; !ok {
		loggers.Logger.Error(ctx, "username or password not provided")
		http.Error(w, "username and/or password not match", http.StatusBadRequest)

		return
	}

	claims := j.FindUser(username, password)
	if _, ok := claims["user"]; !ok {
		loggers.Logger.Error(ctx, "user not found")
		http.Error(w, "username and/or password not match", http.StatusBadRequest)

		return

	}

	payload := map[string]string{}
	token, _ := utils.CreateToken(j.Secret, j.SigningMethod, claims, j.Expire)
	payload["token"] = token
	if j.RefreshToken {
		refreshToken, _ := utils.CreateRefreshToken(j.Secret, j.SigningMethod, token)
		payload["refresh_token"] = refreshToken
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
