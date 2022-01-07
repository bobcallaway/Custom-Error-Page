package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/sirupsen/logrus"
	"mime"
	"os"
	"strconv"
	"strings"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"
	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// RequestId is a unique ID that identifies the request - same as for backend service
	RequestId = "X-Request-ID"

	// ErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler.
	ErrFilesPathVar = "ERROR_FILES_PATH"

	DefaultErrorFilesPath = "templates"

	// DefaultFormatVar is the name of the environment variable indicating
	// the default error MIME type that should be returned if either the
	// client does not specify an Accept header, or the Accept header provided
	// cannot be mapped to a file extension.
	DefaultFormatVar = "DEFAULT_RESPONSE_FORMAT"

	ServerNameVar = "SERVER_NAME"

	DebugVar = "DEBUG"

	ServerListenVar = "LISTEN"

	DefaultServerListen = ":8080"
)

func setupRouter() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file")
	}
	r := gin.New()

	/*	// Optional custom metrics list
		customMetrics := []*ginprometheus.Metric{
			&ginprometheus.Metric{
				ID:	"1234",				// optional string
				Name:	"test_metric",			// required string
				Description:	"Counter test metric",	// required string
				Type:	"counter",			// required string
			},
			&ginprometheus.Metric{
				ID:	"1235",				// Identifier
				Name:	"test_metric_2",		// Metric Name
				Description:	"Summary test metric",	// Help Description
				Type:	"summary", // type associated with prometheus collector
			},
			// Type Options:
			//	counter, counter_vec, gauge, gauge_vec,
			//	histogram, histogram_vec, summary, summary_vec
		}
		p := ginprometheus.NewPrometheus("gin", customMetrics)
	*/

	p := ginprometheus.NewPrometheus("gin")

	p.Use(r)
	templateRoot := os.Getenv(ErrFilesPathVar)
	if templateRoot == "" {
		templateRoot = DefaultErrorFilesPath
	}
	defaultFormat := os.Getenv(DefaultFormatVar)
	if defaultFormat == "" {
		defaultFormat = "text/html"
	}
	isDebug := false
	if os.Getenv(DebugVar) != "" {
		isDebug = true
	}

	serverName := os.Getenv(ServerNameVar)
	r.LoadHTMLGlob(fmt.Sprintf("%s/*", templateRoot))
	r.GET("/", errorHandler(defaultFormat, serverName, isDebug))
	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})
	return r
}

func main() {
	r := setupRouter()
	listen := os.Getenv(ServerListenVar)
	if listen == "" {
		listen = DefaultServerListen
	}
	_ = r.Run(listen)
}
func errorHandler(defaultFormat string, serverName string, debug bool) func(*gin.Context) {
	defaultExts, err := mime.ExtensionsByType(defaultFormat)
	if err != nil || len(defaultExts) == 0 {
		panic("couldn't get file extension for default format")
	}
	defaultExt := defaultExts[0]
	return func(c *gin.Context) {
		ext := defaultExt
		if debug {
			c.Writer.Header().Set(FormatHeader, c.Request.Header.Get(FormatHeader))
			c.Writer.Header().Set(CodeHeader, c.Request.Header.Get(CodeHeader))
			c.Writer.Header().Set(ContentType, c.Request.Header.Get(ContentType))
			c.Writer.Header().Set(OriginalURI, c.Request.Header.Get(OriginalURI))
			c.Writer.Header().Set(Namespace, c.Request.Header.Get(Namespace))
			c.Writer.Header().Set(IngressName, c.Request.Header.Get(IngressName))
			c.Writer.Header().Set(ServiceName, c.Request.Header.Get(ServiceName))
			c.Writer.Header().Set(ServicePort, c.Request.Header.Get(ServicePort))
			c.Writer.Header().Set(RequestId, c.Request.Header.Get(RequestId))
		}
		errCode := c.Request.Header.Get(CodeHeader)
		code, err := strconv.Atoi(errCode)
		if err != nil {
			code = 404
			logrus.Printf("unexpected error reading return code: %v. Using %v", err, code)
		}
		format := c.Request.Header.Get(FormatHeader)
		if format == "" {
			format = defaultFormat
			logrus.Printf("format not specified. Using %v", format)
		}

		cext, err := mime.ExtensionsByType(format)
		if err != nil {
			logrus.Printf("unexpected error reading media type extension: %v. Using %v", err, ext)
			format = defaultFormat
		} else if len(cext) == 0 {
			logrus.Printf("couldn't get media type extension. Using %v", ext)
		} else {
			ext = cext[0]
		}
		c.Writer.Header().Set(ContentType, format)

		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		// special case for compatibility
		if ext == ".htm" {
			ext = ".html"
		}
		openVariables := gin.H{
			"code":        code,
			"originalUrl": c.Request.Header.Get(OriginalURI),
			"ingressName": c.Request.Header.Get(IngressName),
			"serverName":  serverName,
			"requestId":   c.Request.Header.Get(RequestId),
		}
		if ext == ".html" {
			c.HTML(code, "default.html", openVariables)
		} else {
			c.JSON(code, openVariables)
		}
	}

}
