package errorHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(c echo.Context, st *status.Status) error {
	switch st.Code() {
	case codes.InvalidArgument:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": st.Message(),
		})
	case codes.AlreadyExists:
		return c.JSON(http.StatusConflict, map[string]string{
			"error": st.Message(),
		})
	case codes.Unauthenticated:
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": st.Message(),
		})
	case codes.NotFound:
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": st.Message(),
		})
	case codes.Internal:
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": st.Message(),
		})
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": st.Message(),
		})
	}
}
