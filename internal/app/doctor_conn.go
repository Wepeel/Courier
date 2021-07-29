package app

import (
	"github.com/Wepeel/Courier/internal/common"
)

type DoctorConn struct {
	rabbitConn common.RabbitConn
}

func (this DoctorConn) SetupDoctorConn() {
	this.rabbitConn = common.RabbitConn{}
}

func (this DoctorConn) SendMsgToDoctorConn(msg []byte) {
	this.rabbit
}
