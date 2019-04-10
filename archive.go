package tokbox

import (
	"encoding/json"
	"fmt"
)

/*
TODO safari
To include video from streams published from Safari clients, you must use a Safari OpenTok project. Otherwise, streams published from Safari show up as audio-only.
https://tokbox.com/developer/guides/archiving/#individual-stream-and-composed-archives

TODO
The recommended maximum number of streams in an archive is five. You can record up to nine streams; however, in a composed archive quality may degrade if you record more than five streams. If more than nine streams are published to the session, they are not recorded. Also, archive recordings are limited to 120 minutes in length.
https://tokbox.com/developer/guides/archiving/
*/

const (
	archiveBaseUrl = "https://api.opentok.com/v2/project/%s/archive" // api_key
)

type ArchiveService interface {
	Start(ArchiveReq) (resp *Archive, err error)
	Stop(archiveId string) (resp *Archive, err error)
	List(interface{}) (resp *ArchiveResponse, error error)
	Get(archiveId string) (resp *Archive, err error)
	Delete(archiveId string) error
}

type ArchiveServiceOp struct {
	client *Tokbox
}

func (this *ArchiveServiceOp) Start(data ArchiveReq) (resp *Archive, err error) {
	path := fmt.Sprintf(archiveBaseUrl, this.client.apiKey)
	resource := new(Archive)
	err = this.client.Post(path, data, resource)
	return resource, err
}

func (this *ArchiveServiceOp) Stop(archiveId string) (resp *Archive, err error) {
	path := fmt.Sprintf(archiveBaseUrl+"/%s/stop", this.client.apiKey, archiveId)
	resource := new(Archive)
	err = this.client.Post(path, nil, resource)
	return resource, err
}

func (this *ArchiveServiceOp) List(options interface{}) (resp *ArchiveResponse, err error) {
	path := fmt.Sprintf(archiveBaseUrl, this.client.apiKey)
	resource := new(ArchiveResponse)
	err = this.client.Get(path, resource, options)
	return resource, err
}

func (this *ArchiveServiceOp) Get(archiveId string) (resp *Archive, err error) {
	path := fmt.Sprintf(archiveBaseUrl+"/%s", this.client.apiKey, archiveId)
	resource := new(Archive)
	err = this.client.Get(path, resource, nil)
	return resource, err
}

func (this *ArchiveServiceOp) Delete(archiveId string) (err error) {
	path := fmt.Sprintf(archiveBaseUrl+"/%s", this.client.apiKey, archiveId)
	err = this.client.Delete(path, nil, nil)
	return err
}

// url query options
type ArchiveOptions struct {
	Offset    int    `url:"offset,omitempty"`     // default 0
	Count     int    `url:"count,omitempty"`      // default 50
	SessionId string `url:"session_id,omitempty"` // default null, get all archive list
}

// request body
type ArchiveReq struct {
	SessionID string `json:"sessionId"`
	HasAudio  bool   `json:"hasAudio"`
	HasVideo  bool   `json:"hasVideo"`
	Layout    struct {
		Type       string `json:"type"`       // bestFit | custom | horizontalPresentation | pip | verticalPresentation
		Stylesheet string `json:"stylesheet"` // only used with type == custom
	} `json:"layout"`
	Name       string `json:"name"`       // (Optional) The name of the archive (for your own identification)
	OutputMode string `json:"outputMode"` // composed (default) | individual
	Resolution string `json:"resolution"` // 640x480 (default) | 1280x720
}

// response
type ArchiveResponse struct {
	Count int       `json:"count"`
	Items []Archive `json:"items"`
}

func (this *ArchiveResponse) Json() string {
	b, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

type Archive struct {
	CreatedAt  int64       `json:"createdAt"`
	Duration   int         `json:"duration"`
	HasAudio   bool        `json:"hasAudio"`
	HasVideo   bool        `json:"hasVideo"`
	ID         string      `json:"id"` // archive id
	Name       string      `json:"name"`
	OutputMode string      `json:"outputMode"`
	ProjectID  int         `json:"projectId"`
	Reason     string      `json:"reason"`
	Resolution string      `json:"resolution"`
	SessionID  string      `json:"sessionId"`
	Size       int         `json:"size"`
	Status     string      `json:"status"`
	URL        interface{} `json:"url"`
}

func (this *Archive) GetUrl() string {
	if this.URL == nil {
		return ""
	}
	return fmt.Sprintf("%v", this.URL)
}
func (this *Archive) Json() string {
	b, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}
