package message

import (
	"encoding/binary"
	"fmt"
	"github.com/iyouport-org/relaybaton/util"
	"github.com/iyouport-org/socks5"
	"net"
)

type ConnectMessage struct {
	Atyp    byte   //1 {1,3,4}
	Session uint16 //2
	DstPort uint16 //2
	DstAddr []byte //x
	Length  int
}

func (cm ConnectMessage) Pack() []byte {
	msgBytes := []byte{cm.Atyp}
	msgBytes = append(msgBytes, util.Uint16ToBytes(cm.Session)...)
	msgBytes = append(msgBytes, util.Uint16ToBytes(cm.DstPort)...)
	msgBytes = append(msgBytes, cm.DstAddr...)
	cm.Length = len(msgBytes)
	return msgBytes
}

func UnpackConnect(b []byte) ConnectMessage {
	return ConnectMessage{
		Atyp:    b[0],
		Session: binary.BigEndian.Uint16(b[1:3]),
		DstPort: binary.BigEndian.Uint16(b[3:5]),
		DstAddr: b[5:],
		Length:  len(b),
	}
}

func (cm ConnectMessage) GetRequest() *socks5.Request {
	return socks5.NewRequest(socks5.CmdConnect, cm.Atyp, cm.DstAddr, util.Uint16ToBytes(cm.DstPort))
}

func (cm ConnectMessage) GetDstMsg() string {
	switch cm.Atyp {
	case socks5.ATYPIPv4:
		return fmt.Sprintf("%s:%d", net.IP(cm.DstAddr).To4().String(), cm.DstPort)
	case socks5.ATYPIPv6:
		return fmt.Sprintf("[%s]:%d", net.IP(cm.DstAddr).To16().String(), cm.DstPort)
	default: //Domain
		return fmt.Sprintf("%s:%d", string(cm.DstAddr[1:]), cm.DstPort)
	}
}
