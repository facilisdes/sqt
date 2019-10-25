package message

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

const (
	STATUS_ENTRY_NOT_FOUND    = 0
	STATUS_OK_DB              = 10
	STATUS_OK_REDIS           = 11
	STATUS_MAX_QUEUE_EXCEEDED = 20
	STATUS_WRONG_CONFIG       = 21
	STATUS_NO_ACTIVE_STORAGE  = 22
	STATUS_OTHER_ERROR        = 30
	STATUS_WRONG_COMMAND_TYPE = 40
	STATUS_WRONG_COMMAND_KEY  = 41
)

var (
	STATUSES_TEXTS = map[int]string{
		STATUS_ENTRY_NOT_FOUND:    "Data isn't found",
		STATUS_OK_DB:              "Data is found using database server",
		STATUS_OK_REDIS:           "Data is found using cache server",
		STATUS_MAX_QUEUE_EXCEEDED: "Data wasn't queried because max query size was exceeded by this request",
		STATUS_WRONG_CONFIG:       "Data wasn't queried - something is wrong with client's config",
		STATUS_NO_ACTIVE_STORAGE:  "Data wasn't queried - neither database server or cache server is available",
		STATUS_OTHER_ERROR:        "Data wasn't queried due to some unknown error",
		STATUS_WRONG_COMMAND_TYPE: "Wrong command - unsupported command type",
		STATUS_WRONG_COMMAND_KEY:  "Wrong command - wrong key requested",
	}
)

type Message struct {
	IsExecuted       bool
	Status           int
	Data             string
	TimeElapsed      int
	TimeQueuedMin    int
	TimeElapsedTotal int
	QueueSize        int
}

func Serialize(m Message) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	return base64.StdEncoding.EncodeToString(b.Bytes()), err
}

func Deserialize(str string) (Message, error) {
	m := Message{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		b := bytes.Buffer{}
		b.Write(by)
		d := gob.NewDecoder(&b)
		err = d.Decode(&m)
	}
	return m, err
}
