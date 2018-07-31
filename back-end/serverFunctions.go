package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func assureAuthentication(handler func(context *gin.Context, apiUser APIAccount)) func(context *gin.Context) {
	return func(context *gin.Context) {
		var bodyBytes []byte

		if context.Request.Body == nil {
			errorResponse(
				context,
				"Invalid token.",
				"Empty request body.",
			)
			return
		}

		bodyBytes, _ = ioutil.ReadAll(context.Request.Body)

		context.Request.Body.Close()

		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var apiAuth ApiAuther

		if err := json.Unmarshal(bodyBytes, &apiAuth); err != nil {
			errorResponse(
				context,
				"Invalid token.",
				err.Error(),
			)
			return
		}

		var session Session

		if err := db.Preload("Account").Where(Session{UUID: apiAuth.Token}).First(&session).Error; err != nil {
			errorResponse(
				context,
				"Invalid token.",
				err.Error(),
			)
			return
		}

		handler(context, session.Account.APIAccount())
	}
}

func errorResponse(context *gin.Context, msg, debug string) {
	context.JSON(
		400,
		gin.H{
			"error": true,
			"msg":   msg,
			"debug": debug,
		},
	)
}

func hashPassword(password string) string {
	var err error
	var hashedPassword []byte

	for hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14); err != nil; {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	}

	return string(hashedPassword)
}
