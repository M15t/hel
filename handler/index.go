package handler

import (
	"hel/templ/index"

	"github.com/labstack/echo/v4"
)

type Index struct{}

func (h *Index) HandlerIndex(c echo.Context) error {
	return render(c, index.Show())
}
