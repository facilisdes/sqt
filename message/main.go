package message

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

const (
	STATUS_ENTRY_NOT_FOUND = 0
	STATUS_OK_DB = 10
	STATUS_OK_REDIS = 11
	STATUS_MAX_QUEUE_EXCEEDED = 20
	STATUS_WRONG_CONFIG = 21
	STATUS_OTHER_EROR = 30
)

type Message struct {
	IsExecuted bool
	Status int
	Data string
	TimeElapsed int
	TimeQueuedMin int
	QueueSize int
}


// go binary encoder
func ToGOB64(m Message) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil { fmt.Println(`failed gob Encode`, err) }
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) Message {
	m := Message{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil { fmt.Println(`failed base64 Decode`, err); }
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil { fmt.Println(`failed gob Decode`, err); }
	return m
}
