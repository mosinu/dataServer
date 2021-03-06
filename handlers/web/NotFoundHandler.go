package web

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

//NotFoundHandler 404 not found handler
func NotFoundHandler(handlerData HandlerData, w http.ResponseWriter, r *http.Request) {
	log.Info("Not found: ", r.URL.Path)

	err := serveStaticFile(handlerData.Config, NotFoundFile, w)
	if err != nil {
		if os.IsNotExist(err) {
			log.Error("Can't find 404.html!")
			return
		}
	}
}
