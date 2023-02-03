package handler

import (
	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/server"
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/labstack/echo/v4"
)

func GetThemeBundle(serverFactory *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		themeID := ctx.Param("themeId")
		bundleID := ctx.Param("bundleId")

		themeService := serverFactory.Theme()
		theme := themeService.GetTheme(themeID)

		return getBundle(theme.Bundles, bundleID, ctx.Response())
	}
}

func GetTemplateBundle(serverFactory *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		templateID := ctx.Param("templateID")
		bundleID := ctx.Param("bundleID")

		templateService := serverFactory.Template()
		template, err := templateService.Load(templateID)

		if err != nil {
			return derp.NewNotFoundError("handler.GetTemplateBundle", "Template not found", templateID)
		}

		return getBundle(template.Bundles, bundleID, ctx.Response())
	}
}

func getBundle(bundles mapof.Object[model.Bundle], bundleID string, response *echo.Response) error {

	bundle, ok := bundles[bundleID]

	if !ok {
		return derp.NewNotFoundError("handler.getBundle", "Bundle not found", bundleID)
	}

	header := response.Header()
	header.Set("Content-Type", bundle.ContentType)
	header.Set("Cache-Control", bundle.GetCacheControl())
	response.WriteHeader(200)

	if _, err := response.Write(bundle.Content); err != nil {
		return derp.Wrap(err, "handler.getBundle", "Error writing bundle content", bundleID)
	}

	return nil
}
