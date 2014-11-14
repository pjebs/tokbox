package tokbox

//Adapted from https://github.com/cioc/tokbox

import (
	"log"
	"testing"
)

const key = "<your api key here>"
const secret = "<your partner secret here>"

func TestToken(t *testing.T) {
	tokbox := New(key, secret)
	session, err := tokbox.NewSession("", P2P)
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}
	log.Println(session)
	token, err := session.Token(Publisher, "", Hours24)
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}
	log.Println(token)
}
