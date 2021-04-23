package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-http-leaflet-geotag"
	"github.com/sfomuseum/go-http-leaflet-geotag/templates/html"
	"html/template"
	"log"
	"net/http"
)

func PageHandler(templates *template.Template, t_name string) (http.Handler, error) {

	t := templates.Lookup(t_name)

	if t == nil {
		return nil, errors.New("Missing 'map' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		err := t.Execute(rsp, nil)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn), nil
}

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080, "A valid aaronland/go-http-server URI")

	flag.Parse()

	t := template.New("example")

	var err error

	if *path_templates != "" {

		t, err = t.ParseGlob(*path_templates)

		if err != nil {
			log.Fatalf("Failed to parse templates (%s), %v", *path_templates, err)
		}

	} else {

		for _, name := range templates.AssetNames() {

			body, err := templates.Asset(name)

			if err != nil {
				log.Fatal(err)
			}

			t, err = t.Parse(string(body))

			if err != nil {
				log.Fatalf("Failed to parse template (%s), %v", name, err)
			}
		}
	}

	geotag_opts := geotag.DefaultLeafletGeotagOptions()

	mux := http.NewServeMux()

	err = geotag.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatalf("Failed to append leaflet-geotag asset handler, %v", err)
	}

	camera_handler, err := PageHandler(t, "camera")

	if err != nil {
		log.Fatalf("Failed to create camera handler, %v", err)
	}

	camera_handler = geotag.AppendResourcesHandler(camera_handler, geotag_opts)

	mux.Handle("/camera/", camera_handler)

	crosshair_handler, err := PageHandler(t, "crosshair")

	if err != nil {
		log.Fatalf("Failed to create crosshair handler, %v", err)
	}

	crosshair_handler = geotag.AppendResourcesHandler(crosshair_handler, geotag_opts)

	mux.Handle("/crosshair/", crosshair_handler)

	index_handler, err := PageHandler(t, "index")

	if err != nil {
		log.Fatalf("Failed to create index handler, %v", err)
	}

	mux.Handle("/", index_handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening for requests on %s\n", endpoint)

	err = http.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}

}
