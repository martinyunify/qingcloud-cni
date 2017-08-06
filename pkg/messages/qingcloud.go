package messages

type QingcloudInitializeMessage struct {
	QingAccessFile string
	Zone           string
}

type CreateNewNicMessage struct {
	Vxnet []string
}

type GetNicsUnderVxnetMessage struct {
	Hostid  string
	Vxnetid []string
}

type DeleteNicMessage struct {
	nicid string
}
