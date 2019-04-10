package tokbox

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

type Session struct {
	SessionId      string  `json:"session_id"`
	ProjectId      string  `json:"project_id"`
	PartnerId      string  `json:"partner_id"`
	CreateDt       string  `json:"create_dt"`
	SessionStatus  string  `json:"session_status"`
	MediaServerURL string  `json:"media_server_url"`
	T              *Tokbox `json:"-"`
}

func (s *Session) Token(role Role, connectionData string, expiration int64) (string, error) {
	now := time.Now().UTC().Unix()

	dataStr := ""
	dataStr += "session_id=" + url.QueryEscape(s.SessionId)
	dataStr += "&create_time=" + url.QueryEscape(fmt.Sprintf("%d", now))
	if expiration > 0 {
		dataStr += "&expire_time=" + url.QueryEscape(fmt.Sprintf("%d", now+expiration))
	}
	if len(role) > 0 {
		dataStr += "&role=" + url.QueryEscape(string(role))
	}
	if len(connectionData) > 0 {
		dataStr += "&connection_data=" + url.QueryEscape(connectionData)
	}
	dataStr += "&nonce=" + url.QueryEscape(fmt.Sprintf("%d", rand.Intn(999999)))

	h := hmac.New(sha1.New, []byte(s.T.partnerSecret))
	n, err := h.Write([]byte(dataStr))
	if err != nil {
		return "", err
	}
	if n != len(dataStr) {
		return "", fmt.Errorf("hmac not enough bytes written %d != %d", n, len(dataStr))
	}

	preCoded := ""
	preCoded += "partner_id=" + s.T.apiKey
	preCoded += "&sig=" + fmt.Sprintf("%x:%s", h.Sum(nil), dataStr)

	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(preCoded))
	encoder.Close()
	return fmt.Sprintf("T1==%s", buf.String()), nil
}

func (s *Session) Tokens(n int, multithread bool, role Role, connectionData string, expiration int64) []string {
	ret := []string{}

	if multithread {
		var w sync.WaitGroup
		var lock sync.Mutex
		w.Add(n)

		for i := 0; i < n; i++ {
			go func(role Role, connectionData string, expiration int64) {
				a, e := s.Token(role, connectionData, expiration)
				if e == nil {
					lock.Lock()
					ret = append(ret, a)
					lock.Unlock()
				}
				w.Done()
			}(role, connectionData, expiration)

		}

		w.Wait()
		return ret
	} else {
		for i := 0; i < n; i++ {

			a, e := s.Token(role, connectionData, expiration)
			if e == nil {
				ret = append(ret, a)
			}
		}
		return ret
	}
}
