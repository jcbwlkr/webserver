package context

import (
	"net/http"

	"git.wreckerlabs.com/in/webserver/render"

	"github.com/julienschmidt/httprouter"
	"github.com/sdming/gosnow"
)

// Snowflake ID generator
var snowflake, _ = gosnow.Default()

// Event is created on every request and models the event, or request.
// The webserver will populate a RequestContext with any data provided by the
// client from a form, URL, or recognized data type sent in the request body.
type Event struct {
	// ID is generated for each new RequestContext.
	id uint64 `json:"requestID"`
	// RequestContentLength contains a count of incoming bytes.
	RequestContentLength int `json:"requestContentLength"`
	// ResponseContentLength contains a count of outgoing bytes.
	ResponseContentLength int `json:"requestContentLength"`
	// StatusCode
	StatusCode int `json:"statusCode"`

	renderer render.Renderer

	Input          *Input
	Output         *Output
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// New produces a new request context event.
func New(w http.ResponseWriter, req *http.Request, params httprouter.Params) *Event {
	var e = new(Event)

	id, err := snowflake.Next()
	if err != nil {
		panic("Snowflake failed?")
	}
	e.id = id

	e.Input = NewInput()
	e.Output = NewOutput(e)
	e.Request = req
	e.ResponseWriter = w

	return e
}

// HTML renders the HTML view specified by it's filename omitting the file extension.
func (e *Event) HTML(name string, args interface{}) error {
	content, err := render.HTML.Render(name, args)
	if err != nil {
		return err
	}

	e.Output.Body(content)

	return nil
}
