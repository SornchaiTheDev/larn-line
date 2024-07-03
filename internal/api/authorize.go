package api

import (
	"larn-line/internal/constants"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func GetAuthorizationCode(c *gin.Context) {

	params := url.Values{}

	params.Add("response_type", constants.RESPONSE_TYPE)
	params.Add("client_id", os.Getenv("VERIFY_CLIENT_ID"))
	params.Add("redirect_uri", os.Getenv("VERIFY_REDIRECT_URI"))
	params.Add("scope", constants.SCOPE)
	params.Add("state", os.Getenv("VERIFY_REDIRECT_URI"))
	params.Add("nonce", os.Getenv("VERIFY_REDIRECT_URI"))

	authUrl := "https://iot-auththan.ais.co.th/auth/v3.2/oauth/authorize?" + params.Encode()

	c.Redirect(http.StatusFound, authUrl)
}
