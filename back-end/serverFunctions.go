package main

import "github.com/gin-gonic/gin"

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
