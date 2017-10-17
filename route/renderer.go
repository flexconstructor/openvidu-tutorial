package route

import "github.com/gin-gonic/gin"

// renderHTML is a function that renders HTTP pages.
func renderHTML(ctx *gin.Context) {
	ctx.Next()
	ctx.HTML(ctx.Writer.Status(),
		ctx.MustGet("template").(string),
		ctx.MustGet("parameters").(gin.H))
}
