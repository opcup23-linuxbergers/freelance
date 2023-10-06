package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Error          error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

const (
	ErrorUsernameTaken  = "username or email is already taken."
	ErrorBadJSON        = "json is poorly formatted."
	ErrorInternal       = "internal server error, try again later."
	ErrorBadPassword    = "password is too weak."
	ErrorBadCredentials = "email doesn't match with a password."
	ErrorLogic          = "business logic error"
	ErrorRender         = "failed to render a response."
	ErrorInvalidToken   = "token is invalid."
	ErrorInvalidInfo    = "user information mustn't contain whitespaces."
	ErrorNotFound       = "resource not found."
	ErrorReadOnly       = "you can't edit this resource."
	ErrorExists         = "resource already exists."
	ErrorUserNotFound   = "user not found."
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (s *Server) newErrorResponse(err error, code int, status string) render.Renderer {
	response := &ErrorResponse{
		Error:          err,
		HTTPStatusCode: code,
		StatusText:     status,
	}

	if s.DevMode && err != nil {
		response.ErrorText = err.Error()
	}

	return response
}

func (s *Server) SendError(w *http.ResponseWriter, r *http.Request, e error, c int, m string) {
	render.Render(*w, r, s.newErrorResponse(e, c, m))
}
