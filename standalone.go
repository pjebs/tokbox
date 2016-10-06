// +build !appengine

package tokbox

import (
	"golang.org/x/net/context"
	"net/http"
)

func client(ctx *context.Context) *http.Client {
	return &http.Client{}
}
