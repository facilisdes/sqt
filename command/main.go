package command

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

const (
	COMMAND_RUN_QUEUE   = 0
	COMMAND_HEALTHCHECK = 1
)

type Command struct {
	Type       int
	KeyToCheck string
}

func Serialize(c Command) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(c)
	return base64.StdEncoding.EncodeToString(b.Bytes()), err
}

func Deserialize(str string) (Command, error) {
	c := Command{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		b := bytes.Buffer{}
		b.Write(by)
		d := gob.NewDecoder(&b)
		err = d.Decode(&c)
	}
	return c, err
}
