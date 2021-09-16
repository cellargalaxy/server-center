package controller

import (
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server-center/config"
	"github.com/cellargalaxy/server-center/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const tokenKey = "Authorization"
const claimsKey = "claims"

//token检查
func validate(c *gin.Context) {
	token := c.Request.Header.Get(tokenKey)
	logrus.WithContext(c).WithFields(logrus.Fields{"token": token}).Info("解析token")
	tokens := strings.SplitN(token, " ", 2)
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		c.Abort()
		c.JSON(http.StatusOK, util.CreateErrResponse("Authorization非法"))
		return
	}
	jwtToken, err := util.ParseJWT(c, tokens[1], config.Config.Secret, &model.LoginClaims{})
	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, createErrResponse(err.Error()))
		return
	}
	if jwtToken == nil {
		c.Abort()
		c.JSON(http.StatusOK, createErrResponse("JWT token为空"))
		return
	}
	if !jwtToken.Valid {
		c.Abort()
		c.JSON(http.StatusOK, createErrResponse("JWT token非法"))
		return
	}
	claims := jwtToken.Claims
	c.Set(claimsKey, claims)
}
