Tokbox Golang [![GoDoc](http://godoc.org/github.com/pjebs/tokbox?status.svg)](http://godoc.org/github.com/pjebs/tokbox)
=============

This Library is for creating sessions and tokens for the Tokbox Video, Voice & Messaging Platform.
[See Tokbox website](https://tokbox.com/)

It is a hybrid library (supports **Google App Engine** and **Stand-Alone** binary).
It supports **multi-threading** for faster generation of tokens.

**WARNING:** This library uses the **deprecated** api which is only valid until July 2017. I will update this library to the [new API](https://www.tokbox.com/developer/rest/#authentication) before that date.

Install
-------

```shell
go get -u github.com/pjebs/tokbox
```

Usage
-----

```go
import "github.com/pjebs/tokbox"

//setup the api to use your credentials
tb := tokbox.New("<my api key>","<my secret key>")

//create a session
session, err := tb.NewSession("", tokbox.P2P) //no location, peer2peer enabled

//create a token
token, err := session.Token(tokbox.Publisher, "", tokbox.Hours24) //type publisher, no connection data, expire in 24 hours

//Or create multiple tokens
tokens := session.Tokens(5, true, tokbox.Publisher, "", tokbox.Hours24) //5 tokens, multi-thread token generation, type publisher, no connection data, expire in 24 hours. Returns a []string

```

See the unit test for a more detailed example.

Settings
----------

```go
type MediaMode string

const (
	/**
	 * The session will send streams using the OpenTok Media Router.
	 */
	MediaRouter MediaMode = "disabled"
	/**
	* The session will attempt to send streams directly between clients. If clients cannot connect
	* due to firewall restrictions, the session uses the OpenTok TURN server to relay streams.
	 */
	P2P = "enabled"
)

```

**MediaMode** is the second argument in `NewSession` method.


```go
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

```

**Role** is the first argument in `Token` method.

```go
const (
	Days30  = 2592000 //30 * 24 * 60 * 60
	Weeks1  = 604800 //7 * 24 * 60 * 60
	Hours24 = 86400  //24 * 60 * 60
	Hours2  = 7200   //60 * 60 * 2
	Hours1  = 3600   //60 * 60
)

```

**Expiration** value forms the third argument in `Token` method. It dictates how long a token is valid for. The unit is in (seconds) up to a maximum of 30 days.


Methods
----------

	func (t *Tokbox) NewSession(location string, mm MediaMode) (*Session, error)

Creates a new session or returns an error. A session represents a 'virtual chat room' where participants can 'sit in' and communicate with one another. A session can not be deregistered. If you no longer require the session, just discard it's details.

*location string*

The *location* setting is optional, and generally you should keep it as `"".` This setting is an IP address that TokBox will use to situate the session in its global network. If no location hint is passed in (which is recommended), the session uses a media server based on the location of the first client connecting to the session. Pass a location hint in only if you know the general geographic region (and a representative IP address) and you think the first client connecting may not be in that region. If you need to specify an IP address, replace *location* with an IP address that is representative of the geographical location for the session. ([Tokbox - REST API reference](https://tokbox.com/opentok/api/#session_id_production))

*mm MediaMode*

`P2P` will direct clients to transfer video-audio data between each other directly (if possible).

`MediaRouter` directs data to go through Tokbox's Media Router servers. Integrates **Intelligent Quality Control** technology to improve user-experience (albeit at higher pricing). ([Tokbox - REST API reference](https://tokbox.com/opentok/api/#session_id_production))


	func (s *Session) Token(role Role, connectionData string, expiration int64) (string, error)

Generates a token for a corresponding session. Returns a string representing the token value or returns an error. A token represents a 'ticket' allowing participants to 'sit in' a session. The permitted range of activities is determined by the `role` setting.

*role Role*

`Publisher` - allows participant to broadcast their own audio and video feed to other participants in the session. They can also listen to and watch broadcasts by other members of the session.

`Subscriber` - allows participants to **only** listen to and watch broadcasts by other participants in the session with **Publisher** rights.

*connectionData string*

`connectionData` - Extra arbitrary data that can be read by other clients. ([Tokbox - Generating Tokens](https://tokbox.com/opentok/libraries/server/php/))

*expiration int64*

`expiration` - How long the token is valid for. The unit is in (seconds) up to a maximum of 30 days. See above for built-in enum values, or use your own.

 	func (s *Session) Tokens(n int, multithread bool, role Role, connectionData string, expiration int64) []string

 Generates multiple (`n`) tokens in one go. Returns a `[]string.` All tokens generated have the same settings as dictated by *role*, *connectionData* and *expiration*. See above for more details. Since the function repeatedly calls the `Token` method, any error in token generation is ignored. Therefore it may be prudent to check if the length of the returned `[]string` matches `n.`

 *multithread bool*

 If `true,` the function strives to generate the tokens concurrently (if multiple CPU cores are available). Preferred if *many, many, many* tokens need to be generated in one go.


Credits: 
--------
(This library is based on the older tokbox library – no longer in active development)

https://github.com/cioc/tokbox by Charles Cary


