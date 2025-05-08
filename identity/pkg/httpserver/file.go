// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package httpserver

import (
	"net/http"
	"os"

	"github.com/agntcy/identity/pkg/httpserver/assets"
)

func FileServer(dir string) http.Handler {
	staticFS := assets.NewHttpStaticFS(os.DirFS(dir))
	handler := http.FileServer(http.FS(staticFS))

	return handler
}
