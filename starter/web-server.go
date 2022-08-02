package starter

import (
	"github.com/XM-GO/PandaKit/config"
	"github.com/XM-GO/PandaKit/logger"
	"github.com/gin-gonic/gin"
)

func RunWebServer(web *gin.Engine) {
	server := config.Conf.Server
	port := server.GetPort()
	if app := config.Conf.App; app != nil {
		logger.Log.Infof("%s- Listening and serving HTTP on %s", app.GetAppInfo(), port)
	} else {
		logger.Log.Infof("Listening and serving HTTP on %s", port)
	}

	if server.Tls != nil && server.Tls.Enable {
		web.RunTLS(port, server.Tls.CertFile, server.Tls.KeyFile)
	} else {
		web.Run(port)
	}
}
