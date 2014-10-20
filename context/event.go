package context

import (
	"net/http"
	"time"

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
	// StartTime can be used for performance metrics.
	StartTime time.Time `json:"startTime"`
	// Duration is set when a request is concluded and is a measure of how
	// long a request has taken.
	Duration time.Time `json:"duration"`
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
	e.StartTime = time.Now()

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

func (e Event) getID() uint64 {
	return e.id
}

// ---------------------
// ---------------------
// ---------------------

// HTML renders the HTML view specified by it's filename omitting the file extension.
func (e *Event) HTML(name string, args interface{}) error {
	content, err := render.HTML.Render(name, nil)
	if err != nil {
		return err
	}

	e.Output.Body(content)

	return nil
}

// LayoutHTML renders the HTML layout and embedded view specified by their filenames omitting the file extension.
func (e *Event) LayoutHTML(layout string, template string, args interface{}) {

}
