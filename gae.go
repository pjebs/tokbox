// +build appengine

package tokbox

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func client(ctx *context.Context) *http.Client {
	if ctx == nil {
		return &http.Client{}
	} else {
		return urlfetch.Client(*ctx)
	}
}
