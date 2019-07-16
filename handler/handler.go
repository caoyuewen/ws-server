package handler

import (
	"encoding/binary"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"ws-server/mconst/msgid"
)

// onMessage
func OnMessage(m *melody.Melody, session *melody.Session, bytes []byte) {

}

// onClose
func OnClose(session *melody.Session, i int, s string) error {
	logger.Info("Handler OnClose %d,%s :", i, s)
	return nil
}

// onConnection
func OnConnect(m *melody.Melody, session *melody.Session) {
	defer func() {
		logger.Info("当前连接数:", m.Len())
	}()
	logger.Info("a new connect :remote addr",session.Request.RemoteAddr)
	handlerAuthPush(session)
}

// onDisconnection
func OnDisconnect(m *melody.Melody, session *melody.Session) {
	value, _ := session.Get("auth")
	logger.Info("当前连接数", m.Len(), value)
	logger.Info("Handler OnDisconnect :")
}

// onErr
func OnErr(session *melody.Session, e error) {
	logger.Info("Handler OnErr :", e.Error())
}

// onPong
func OnPong(session *melody.Session) {

	// logger.Info("Handler OnPong :")
}

// onSentMsg
func OnSentMessage(session *melody.Session, bytes []byte) {
	logger.Info("Handler OnSentMessage :", string(bytes))

}

// OnMessageBinary
func OnMessageBinary(session *melody.Session, bytes []byte) {
	logger.Info("Handler OnMessageBinary original :", bytes)
	logger.Info("Handler OnMessageBinary :", string(bytes))
	logger.Info("Handler OnMessageBinary len::", len(bytes))

	dataLen := len(bytes)
	if dataLen <= 2 {
		_ = session.WriteBinary([]byte("data is nil !"))
		return
	}

	msgId := getMsgId(bytes[:2])
	data := bytes[2:]
	logger.Info("msgId", msgId)

	switch msgId {
	case msgid.Auth:
		handlerAuthReq(session, data)
	}

}

// OnSentMessageBinary
func OnSentMessageBinary(session *melody.Session, bytes []byte) {
	//logger.Info("Handler OnSentMessageBinary :", string(bytes))
}

func getMsgId(mid []byte) uint16 {
	u := binary.BigEndian.Uint16(mid)
	return u
}

func pkgMsg(msgId uint16, data []byte) []byte {

	bytes := make([]byte, 2+len(data))

	binary.BigEndian.PutUint16(bytes[:2], msgId)

	copy(bytes[2:],data)

	return bytes

}
