package tokbox

import (
	"bytes"

	"net/http"
	"net/url"

	"encoding/base64"
	"encoding/xml"

	"crypto/hmac"
	"crypto/sha1"

	"fmt"
	"math/rand"
	"strings"
	"time"

	"sync"

	"golang.org/x/net/context"
)

const (
	apiHost    = "https://api.opentok.com/hl"
	apiSession = "/session/create"
)

const (
	Days30  = 2592000 //30 * 24 * 60 * 60
	Weeks1  = 604800  //7 * 24 * 60 * 60
	Hours24 = 86400   //24 * 60 * 60
	Hours2  = 7200    //60 * 60 * 2
	Hours1  = 3600    //60 * 60
)

type MediaMode string

const (
	/**
	 * The session will send streams using the OpenTok Media Router.
	 */
	MediaRouter MediaMode = "disabled"
	/**
	* The session will attempt send streams directly between clients. If clients cannot connect
	* due to firewall restrictions, the session uses the OpenTok TURN server to relay streams.
	 */
	P2P = "enabled"
)

type Role string

const (
	/**
	* A publisher can publish streams, subscribe to streams, and signal.
	 */
	Publisher Role = "publisher"
	/**
	* A subscriber can only subscribe to streams.
	 */
	Subscriber = "subscriber"
	/**
	* In addition to the privileges granted to a publisher, in clients using the OpenTok.js 2.2
	* library, a moderator can call the <code>forceUnpublish()</code> and
	* <code>forceDisconnect()</code> method of the Session object.
	 */
	Moderator = "moderator"
)

type sessions struct {
	Sessions []Session `xml:"Session"`
}

type Tokbox struct {
	apiKey        string
	partnerSecret string
	BetaUrl       string //Endpoint for Beta Programs
}

type Session struct {
	SessionId     string `xml:"session_id"`
	PartnerId     string `xml:"partner_id"`
	CreateDt      string `xml:"create_dt"`
	SessionStatus string `xml:"session_status"`
	T             *Tokbox
}

func New(apikey, partnerSecret string) *Tokbox {
	return &Tokbox{apikey, partnerSecret, ""}
}

func (t *Tokbox) NewSession(location string, mm MediaMode, ctx ...*context.Context) (*Session, error) {
	params := url.Values{}

	if len(location) > 0 {
		params.Add("location", location)
	}

	params.Add("p2p.preference", string(mm))

	var endpoint string
	if t.BetaUrl == "" {
		endpoint = apiHost
	} else {
		endpoint = t.BetaUrl
	}
	req, err := http.NewRequest("POST", endpoint+apiSession, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	authHeader := t.apiKey + ":" + t.partnerSecret
	req.Header.Add("X-TB-PARTNER-AUTH", authHeader)

	if len(ctx) == 0 {
		ctx = append(ctx, nil)
	}
	res, err := client(ctx[0]).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Tokbox returns error code: %v", res.StatusCode)
	}

	var s sessions
	if err = xml.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}

	if len(s.Sessions) < 1 {
		return nil, fmt.Errorf("Tokbox did not return a session")
	}

	o := s.Sessions[0]
	o.T = t
	return &o, nil
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
