package auth

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	identityKey = "username"
	AuthUserKey = "auth_user"
)

var (
	AuthMiddleware *jwt.GinJWTMiddleware

	gjm = &jwt.GinJWTMiddleware{
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     generatePayLoad,
		IdentityHandler: identityHandler,
		Authenticator:   authenticate,
		Authorizator:    authorize,
		LoginResponse:   loginResponse,
		RefreshResponse: loginResponse,
		LogoutResponse:  logoutResponse,
		Unauthorized:    unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

type UnauthorizedResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func init() {
	AuthMiddleware, _ = jwt.New(gjm)
}

func generatePayLoad(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	u := &User{
		Username: claims[identityKey].(string),
	}

	c.Set(AuthUserKey, u)
	return u
}

func authenticate(c *gin.Context) (interface{}, error) {
	user := &User{}
	if err := c.ShouldBind(user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	err := validateUser(user)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, nil
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusCreated,
		&LoginResp{
			Token:  token,
			Expire: expire.Format(time.RFC3339),
		})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusCreated, struct{}{})
}

func authorize(data interface{}, c *gin.Context) bool {
	v, ok := data.(*User)
	if !ok {
		return false
	}

	err := validateUser(v)
	if err != nil {
		return false
	}

	return true
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code,
		&UnauthorizedResp{
			Code:    code,
			Message: message,
		},
	)
}

func validateUser(u *User) error {
	return nil
}
