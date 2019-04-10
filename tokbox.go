package tokbox

import (
	"bytes"
	"github.com/google/go-querystring/query"
	"io/ioutil"

	"net/http"
	"net/url"

	"encoding/json"

	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/dgrijalva/jwt-go"
	"github.com/myesui/uuid"
)

const (
	apiHost    = "https://api.opentok.com"
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

type Tokbox struct {
	Client        *http.Client
	apiKey        string
	partnerSecret string
	BetaUrl       string         //Endpoint for Beta Programs
	Archive       ArchiveService // Archive sdk
}

func New(apikey, partnerSecret string) *Tokbox {
	//return &Tokbox{apikey, partnerSecret, "",nil}
	tb := &Tokbox{
		Client:        http.DefaultClient,
		apiKey:        apikey,
		partnerSecret: partnerSecret,
		BetaUrl:       "",
	}
	tb.Archive = &ArchiveServiceOp{tb}
	return tb
}

func (t *Tokbox) jwtToken() (string, error) {

	type TokboxClaims struct {
		Ist string `json:"ist,omitempty"`
		jwt.StandardClaims
	}

	claims := TokboxClaims{
		"project",
		jwt.StandardClaims{
			Issuer:    t.apiKey,
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().UTC().Unix() + (2 * 24 * 60 * 60), // 2 hours; //NB: The maximum allowed expiration time range is 5 minutes.
			Id:        uuid.NewV4().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.partnerSecret))
}

// Creates a new tokbox session or returns an error.
// See README file for full documentation: https://github.com/pjebs/tokbox
// NOTE: ctx must be nil if *not* using Google App Engine
func (t *Tokbox) NewSession(location string, mm MediaMode, ctx ...context.Context) (*Session, error) {
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

	//Create jwt token
	jwt, err := t.jwtToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-OPENTOK-AUTH", jwt)

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

	var s []Session
	if err = json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}

	if len(s) < 1 {
		return nil, fmt.Errorf("Tokbox did not return a session")
	}

	o := s[0]
	o.T = t
	return &o, nil
}

func (this *Tokbox) NewRequest(method, urlStr string, body, options interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		for k, values := range rel.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
	}
	var js []byte = nil
	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, rel.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	token, err := this.jwtToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-OPENTOK-AUTH", token)
	return req, nil

}

func (this *Tokbox) CreateAndDo(method, path string, data, options, resource interface{}) error {
	req, err := this.NewRequest(method, path, data, options)
	if err != nil {
		return err
	}
	err = this.Do(req, resource)
	if err != nil {
		return err
	}
	return nil
}
func (this *Tokbox) Do(req *http.Request, v interface{}) error {
	resp, err := this.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = CheckResponseError(resp)
	if err != nil {
		return err
	}
	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckResponseError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	responseError := &ResponseError{}
	err = json.Unmarshal(bodyBytes, responseError)
	if err != nil {
		responseError.Code = resp.StatusCode
		responseError.Message = string(bodyBytes)
		responseError.Description = codeErrDict[resp.StatusCode]
	}
	return responseError
}

func (this *Tokbox) Get(path string, resource, options interface{}) error {
	return this.CreateAndDo("GET", path, nil, options, resource)
}
func (this *Tokbox) Post(path string, data, resource interface{}) error {
	return this.CreateAndDo("POST", path, data, nil, resource)
}
func (this *Tokbox) Put(path string, data, resource interface{}) error {
	return this.CreateAndDo("PUT", path, data, nil, resource)
}
func (this *Tokbox) Delete(path string, resource, options interface{}) error {
	return this.CreateAndDo("DELETE", path, nil, options, resource)
}
