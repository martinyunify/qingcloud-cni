package messages

import (
	"fmt"
	"os"
)

type QingcloudInitializeMessage struct {
	QingAccessFile string
	Zone           string
}

const (
	PEK3A = iota
	PEK3B
	GD1
	GD2A
	SHA1
)

var zoneMap = map[string]uint{
	"pek3a": PEK3A,
	"pek3b": PEK3B,
	"gd1":   GD1,
	"gd2a":  GD2A,
	"sha1":  SHA1,
}

func NewQingcloudInitializeMessage(filepath string, zone string) (*QingcloudInitializeMessage, error) {
	if filepath == "" {
		return nil, fmt.Errorf("Access File Path is emtpy")
	}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Access File does not exist %s", filepath)
	}
	if zone == "" {
		return nil, fmt.Errorf("Zone is empty")
	}
	if _, ok := zoneMap[zone]; !ok {
		return nil, fmt.Errorf("Zone is invalid %s", zone)
	}
	return &QingcloudInitializeMessage{
		Zone:           zone,
		QingAccessFile: filepath,
	}, nil
}

type GetNicsUnderVxnetMessage struct {
	Hostid  string
	Vxnetid []string
}

type GetNicsUnderVxnetReplyMessage struct {
}
