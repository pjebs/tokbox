// +build !appengine

package tokbox

import (
	"net/http"
)

func client() *http.Client {
	return &http.Client{}
}
