package geotag

import (
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-leaflet"
	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-http-leaflet-geotag/static"
)

var INCLUDE_LEAFLET = true

type LeafletGeotagOptions struct {
	JS  []string
	CSS []string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
}

func DefaultLeafletGeotagOptions() *LeafletGeotagOptions {

	opts := &LeafletGeotagOptions{
		CSS: []string{
			"/css/Leaflet.GeotagPhoto.css",
			"/css/highlight.js.default.min.css",
		},
		JS: []string{
			"/javascript/Leaflet.GeotagPhoto.js",
			"/javascript/highlight.min.js",
		},
	}

	return opts
}

func AppendResourcesHandler(next http.Handler, opts *LeafletGeotagOptions) http.Handler {

	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

func AppendResourcesHandlerWithPrefix(next http.Handler, opts *LeafletGeotagOptions, prefix string) http.Handler {

	if INCLUDE_LEAFLET {
		leaflet_opts := leaflet.DefaultLeafletOptions()
		leaflet_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF
		next = leaflet.AppendResourcesHandlerWithPrefix(next, leaflet_opts, prefix)
	}

	static_opts := aa_static.DefaultResourcesOptions()
	static_opts.CSS = opts.CSS
	static_opts.JS = opts.JS
	static_opts.AppendJavaScriptAtEOF = opts.AppendJavaScriptAtEOF

	return aa_static.AppendResourcesHandlerWithPrefix(next, static_opts, prefix)
}

func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	if INCLUDE_LEAFLET {

		err := leaflet.AppendAssetHandlersWithPrefix(mux, prefix)

		if err != nil {
			return fmt.Errorf("Failed to append Leaflet assets, %w", err)
		}
	}

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, prefix)
}
