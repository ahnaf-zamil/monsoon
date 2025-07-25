package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"

	"monsoon/api"

	"github.com/gin-gonic/gin"
)

/*
This is the place where I dump all my helper functions

May refactor/reorganize later
*/

/* General helper functions */
func GenerateSocketID() string {
	// Generates 8 byte IDs for sockets
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func PrettyPrintSyncMap(data *sync.Map) {
	// Handy function for pretty printing sync.Map objects
	m := map[string]any{}
	data.Range(func(key, value any) bool {
		m[fmt.Sprint(key)] = value
		return true
	})

	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

/* API-sided */
func WriteAPIResponse(c *gin.Context, data any, rs *api.APIResponse, r_code int) {
	rs.Err = false
	rs.Message = ""
	rs.Status = r_code
	rs.Data = data
	c.JSON(r_code, rs)
}

func WriteAPIError(c *gin.Context, message string, rs *api.APIResponse, r_code int) {
	rs.Err = true
	rs.Message = message
	rs.Status = r_code
	c.JSON(r_code, rs)
}

func HandleServerError(c *gin.Context, rs *api.APIResponse, err error) {
	// Instead of rewriting the same two lines every single time, just run it for error cases
	log.Println(err)
	WriteAPIError(c, "Internal server error", rs, http.StatusInternalServerError)
}

func RandomBase16String(l int) (string, error) {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	_, err := rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l], err // strip 1 extra character we get from odd length results
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	/* For dev environments, setting SameSite = lax since no HTTPS
	Chrome will not send this cookie on API requests, so use Firefox for development

	On prod, assume HTTPS and set SameSite = none with Secure
	*/

	var sameSite http.SameSite
	var secure bool
	domain := c.Request.Host

	if IsDevEnv() {
		sameSite = http.SameSiteLaxMode
		secure = false
		domain = "localhost"
	} else {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     api.COOKIE_REFRESH_TOKEN,
		Value:    refreshToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func GetCurrentUser(c *gin.Context) (*api.UserModel, bool) {
	userAny, exists := c.Get("current_user")
	if !exists {
		return nil, false
	}

	user, ok := userAny.(*api.UserModel)
	if !ok {
		return nil, false
	}

	return user, true
}
