package main

import (
	"events/services/events"
	"events/services/users"
	"github.com/iris-contrib/middleware/pg"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/logger"
)

// Read the example and its comments carefully.
func makeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	ac := accesslog.File("access.log")
	// Remove this line to disable logging to console:
	//ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = true
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = true
	ac.BytesSent = true
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	return ac
}

func main() {
	ac := makeAccessLog()
	defer ac.Close() // Close the underline file.
	// Customize the logger configuration
	customLogger := logger.New(logger.Config{
		Status:             true,
		IP:                 true,
		Method:             true,
		Path:               true,
		Query:              true,
		TraceRoute:         false,
		PathAfterHandler:   false,
		MessageContextKeys: []string{"logger_message"}, // Customize the log message key
		MessageHeaderKeys:  []string{"X-Request-ID"},   // Include additional header keys in the log
	})

	app := iris.New()
	// Register the middleware (UseRouter to catch http errors too).
	app.UseRouter(ac.Handler)
	// Add the custom logger middleware
	app.Use(customLogger)
	//app.Use(logger.New())
	app.Logger().SetLevel("debug")
	app.Use(iris.Compression)

	db := newPostgresMiddleware()
	events.InitRoutes(app, db)
	users.InitRoutes(app, db)
	app.Listen(":8080")
}

func newPostgresMiddleware() iris.Handler {
	schema := pg.NewSchema()
	schema.MustRegister("users", users.User{})
	schema.MustRegister("events", events.Event{})

	opts := pg.Options{
		Host:          "localhost",
		Port:          5432,
		User:          "godevops",
		Password:      "godevops123",
		DBName:        "godevops",
		Schema:        "public",
		SSLMode:       "disable",
		Transactional: true,  // or false to disable the transactional feature.
		Trace:         false, // or false to production to disable query logging.
		CreateSchema:  true,  // true to create the schema if it doesn't exist.
		CheckSchema:   false, // true to check the schema for missing tables and columns.
		ErrorHandler:  func(ctx iris.Context, err error) {},
	}
	p := pg.New(schema, opts)
	return p.Handler()
}
