package httpserver

import (
	"net/http"
	"os"

	"github.com/agntcy/pyramid/pkg/httpserver/assets"
)

func FileServer(dir string) http.Handler {
	staticFS := assets.NewHttpStaticFS(os.DirFS(dir))
	handler := http.FileServer(http.FS(staticFS))
	return handler
}
