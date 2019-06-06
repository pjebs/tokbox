// Author: majianyu
// Create Date: 2019-06-06
// Description:
// version V1.0
package tokbox

import (
	"encoding/json"
	"fmt"
)

const (
	streamsBaseUrl = "https://api.opentok.com/v2/project/%s/broadcast" // api_key
)

type StreamService interface {
	List(interface{}) (resp *StreamResponse, error error)
}
type StreamResponse struct {
	Count int          `json:"count"`
	Items []StreamItem `json:"items"`
}
type StreamItem struct {
	ID            string `json:"id"`
	SessionID     string `json:"sessionId"`
	ProjectID     int    `json:"projectId"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
	Resolution    string `json:"resolution"`
	BroadcastUrls struct {
		Hls  string `json:"hls"`
		Rtmp struct {
			Foo struct {
				ServerURL  string `json:"serverUrl"`
				StreamName string `json:"streamName"`
				Status     string `json:"status"`
			} `json:"foo"`
			Bar struct {
				ServerURL  string `json:"serverUrl"`
				StreamName string `json:"streamName"`
				Status     string `json:"status"`
			} `json:"bar"`
		} `json:"rtmp"`
	} `json:"broadcastUrls"`
	Status string `json:"status"`
}
type StreamOption struct {
	Offset    int    `url:"offset,omitempty"`
	Count     int    `url:"count,omitempty"`
	SessionID string `url:"session_id,omitempty"`
}

func (this StreamResponse) Json() string {
	bytes, err := json.MarshalIndent(this, "", "\t")
	if err != nil {
		return ""
	}
	return string(bytes)

}

type StreamServiceOp struct {
	client *Tokbox
}

func (this *StreamServiceOp) List(options interface{}) (resp *StreamResponse, err error) {
	path := fmt.Sprintf(streamsBaseUrl, this.client.apiKey)
	resp = new(StreamResponse)
	err = this.client.Get(path, resp, options)
	return resp, err
}
