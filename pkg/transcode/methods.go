package transcode

import (
	"bytes"
	"encoding/json"

	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
)

func Transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func CopyMap(in map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {
	out := make(map[string]rtspsimpleserver.SConf)
	for k, v := range in {
		out[k] = v
	}
	return out
}
