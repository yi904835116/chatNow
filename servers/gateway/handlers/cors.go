package handlers

import "net/http"

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/
// CORSHandler creates a CORS middleware handler that wraps these other handler functions
type CORSHandler struct {
	wrappedHandler http.Handler
}

// NewCORShandler returns CORSHandler
func NewCORShandler(handlerToWrap http.Handler) http.Handler {
	return &CORSHandler{handlerToWrap}
}

func (c *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(headerAccessControlAllowOrigin, "*")
	w.Header().Add(headerAccessControlAllowMethods, "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add(headerAccessControlAllowHeaders, "Content-Type, Authorization")
	w.Header().Add(headerAccessControlExposeHeaders, "Authorization")
	w.Header().Add(headerAccessControlMaxAge, "600")

	if r.Method != "OPTIONS" {
		c.wrappedHandler.ServeHTTP(w, r)
	}

}
