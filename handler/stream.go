package handler

import (
	"github.com/benpate/derp"
	"github.com/benpate/ghost/service"
	"github.com/benpate/presto"
	"github.com/labstack/echo/v4"
)

// GetStream generates the base HTML for a stream
func GetStream(maker service.FactoryMaker, roles ...presto.RoleFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		// Get the service factory
		factory := maker.Factory(ctx.Request().Context())
		defer factory.Close()

		// Get the stream service
		streamService := factory.Stream()

		stream, err := streamService.LoadByToken(ctx.Param("token"))

		if err != nil {
			err = derp.Wrap(err, "ghost.handler.GetStream", "Error loading stream from service")
			derp.Report(err)
			return err
		}

		pageService := factory.PageService()

		// Generate the result
		result, err := streamService.Render(stream, ctx.Param("view"))

		if err != nil {
			derp.Report(err)
			return ctx.String(err.Code, "")
		}

		header, footer := pageService.Render(ctx, stream, ctx.Param("view"))

		// Success!
		response := ctx.Response()
		response.WriteHeader(200)
		response.Write([]byte(header))
		response.Write([]byte(result))
		response.Write([]byte(footer))

		return nil
	}
}
