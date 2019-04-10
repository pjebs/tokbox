package tokbox

//Adapted from https://github.com/cioc/tokbox

import (
	"log"
	"testing"
)

const key = "<your api key here>"

const secret = "<your partner secret here>"

var tb = New(key, secret)

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
func TestArchiveStart(t *testing.T) {
	sessionId := ""
	resp, err := tb.Archive.Start(ArchiveReq{
		SessionID: sessionId,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp)
}
func TestArchiveList(t *testing.T) {
	resp, err := tb.Archive.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp)
}

func TestArchiveGet(t *testing.T) {
	archiveId := ""
	resp, err := tb.Archive.Get(archiveId)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resp)

}
func TestArchiveDelete(t *testing.T) {
	archiveId := ""
	err := tb.Archive.Delete(archiveId)
	if err != nil {
		t.Fatal(err)
	}
}