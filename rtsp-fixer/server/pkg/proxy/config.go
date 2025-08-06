package proxy

import (
	"encoding/json"
	"net/url"
)

type ParsableURL url.URL

func (me *ParsableURL) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*me = ParsableURL(*u)
	return nil
}

func (me *ParsableURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.String())
}

func (me *ParsableURL) String() string {
	return (*url.URL)(me).String()
}

type StreamConfig struct {
	Name                   string      `json:"name"`
	URL                    ParsableURL `json:"url"`
	FixForceTCPInTransport bool        `json:"fix_force_tcp_in_transport"`
}
