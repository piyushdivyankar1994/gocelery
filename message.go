package gocelery

import (
	"fmt"
	"os"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type CeleryMessageProperties struct {
	CorrelationID   string
	ContentType     string
	ContentEncoding string
	ReplyTo         string
}

type CeleryMessageHeaders map[string]interface{}

type CeleryMessageBody struct {
	Args   []interface{}          `json:"args"`
	Kwargs map[string]interface{} `json:"kwargs"`
	Embed  BodyEmbed              `json:"embed"`
}

type BodyEmbed struct {
	Callbacks []string `json:"callbacks"`
	Errbacks  []string `json:"errbacks"`
	Chain     []string `json:"chain"`
	Chord     string   `json:"chord"`
}

type TaskMessageV2 struct {
	headers    CeleryMessageHeaders
	body       CeleryMessageBody
	properties CeleryMessageProperties
}

func (tm2 *TaskMessageV2) reset() {
	tm2.properties.CorrelationID = uuid.Must(uuid.NewV4()).String()
	tm2.properties.ReplyTo = uuid.Must(uuid.NewV4()).String()
	tm2.headers = CeleryMessageHeaders{
		"task":       "",
		"argsrepr":   "[]",
		"kwargsrepr": "{}",
	}
	tm2.body.Args = nil
	tm2.body.Kwargs = nil
	tm2.body.Embed = BodyEmbed{}
}

var taskMessageV2Pool = sync.Pool{
	New: func() interface{} {
		eta := time.Now().Format(time.RFC3339)
		hostname, _ := os.Hostname()
		return &TaskMessageV2{
			headers: CeleryMessageHeaders{
				"lang":       "go",
				"origin":     fmt.Sprintf("@%v%v", os.Getpid(), hostname),
				"eta":        eta,
				"task":       "",
				"argsrepr":   "[]",
				"kwargsrepr": "{}",
			},
			body: CeleryMessageBody{},
			properties: CeleryMessageProperties{
				CorrelationID:   uuid.Must(uuid.NewV4()).String(),
				ReplyTo:         uuid.Must(uuid.NewV4()).String(),
				ContentType:     "application/json",
				ContentEncoding: "base64",
			},
		}
	},
}

func getTaskMessageV2(task string) *TaskMessageV2 {
	msg := taskMessageV2Pool.Get().(*TaskMessageV2)
	msg.headers["task"] = task
	msg.headers["argsrepr"] = "[]"
	msg.headers["kwargsrepr"] = "{}"
	return msg
}

func releaseTaskMessageV2(v *TaskMessageV2) {
	v.reset()
	taskMessageV2Pool.Put(v)
}
