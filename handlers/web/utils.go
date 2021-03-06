package web

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/JojiiOfficial/DataManagerServer/models"
	"github.com/JojiiOfficial/gaw"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//HandlerData handlerData for web
type HandlerData struct {
	Config *models.Config
	Db     *gorm.DB
	User   *models.User
}

//LogError returns true on error
func LogError(err error, context ...map[string]interface{}) bool {
	if err == nil {
		return false
	}

	if len(context) > 0 {
		log.WithFields(context[0]).Error(err.Error())
	} else {
		log.Error(err.Error())
	}
	return true
}

//Copy stream
func serveFileStream(config *models.Config, reader io.Reader, w http.ResponseWriter) {
	_ = gaw.BufferedCopy(config.Webserver.DownloadFileBuffer, w, reader)
}

//Detect and set Content-Type by extension
func autoSetContentType(w http.ResponseWriter, file string) {
	setContentType(w, mime.TypeByExtension(file))
}

//Set Content-Type
func setContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set(models.HeaderContentType, fmt.Sprintf("%s; charset=utf-8", contentType))
}

//Serve static file
func serveStaticFile(config *models.Config, file string, w http.ResponseWriter, contentType ...string) error {
	//Open file
	f, err := os.Open(config.GetHTMLFile(file))
	defer f.Close()

	if LogError(err) {
		return err
	}

	//Set contentType
	if len(contentType) == 0 || len(contentType[0]) == 0 {
		autoSetContentType(w, file)
	} else {
		w.Header().Set(models.HeaderContentType, contentType[0])
	}

	serveFileStream(config, f, w)
	return nil
}

//Handles errors and respond with 404 if this caused the error
func handleBrowserServeError(err error, handlerData HandlerData, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println(err)
		if os.IsNotExist(err) {

			NotFoundHandler(handlerData, w, r)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
