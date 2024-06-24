package utils

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderTempl(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
