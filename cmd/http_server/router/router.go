package router

import (
	"net/http"

	"github.com/aria3ppp/rag-server/internal/rag/usecase"

	"github.com/labstack/echo/v4"
)

type QueryServer struct {
	uc usecase.UseCase
}

func New(uc usecase.UseCase) *QueryServer {
	return &QueryServer{uc: uc}
}

func (server *QueryServer) HandleQuery(ctx echo.Context) error {
	query := ctx.FormValue("query")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "query form is not provided"})
	}

	queryResult, err := server.uc.QuerySync(ctx.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": queryResult})
}
