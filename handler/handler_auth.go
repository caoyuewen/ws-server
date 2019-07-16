package handler

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/wonderivan/logger"
	"gopkg.in/olahol/melody.v1"
	"ws-server/cert"
	"ws-server/mconst/msgid"
	"ws-server/msg/mproto"

	"time"
)

// 100 Auth push
func handlerAuthPush(session *melody.Session) {
	authKey := uuid.New().String()

	ensNonce, err := cert.GetSNonce(authKey)
	if err != nil {
		_ = session.WriteBinary([]byte("server auth err"))
		return
	}

	var authResp mproto.AuthResp
	authResp.Nonce = ensNonce

	bytes, err := proto.Marshal(&authResp)
	if err != nil {
		_ = session.WriteBinary([]byte("server auth err"))
		return
	}

	session.Set("auth_key", authKey)

	logger.Debug("ensNonce:", ensNonce)
	logger.Debug("original", authKey)

	_ = session.WriteBinary(pkgMsg(msgid.Auth, bytes))
	go delayRemoveSession(session)
}

// 100 verify client
func handlerAuthReq(session *melody.Session, bytes []byte) {
	logger.Info("handler auth req")
	var authReq mproto.AuthReq
	err := proto.Unmarshal(bytes, &authReq)
	if err != nil {
		logger.Error("data is not correct:", err.Error())
		_ = session.WriteBinary([]byte("data is not correct"))
		return
	}

	saveAuthKey, exists := session.Get("auth_key")
	if !exists {
		logger.Warn("an illegal connection", session.Request.RemoteAddr)
		_ = session.WriteBinary([]byte("fuck you!"))
		_ = session.Close()
		return
	}

	nonce := authReq.Nonce
	decodeNonce, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		logger.Info("decode nonce fail", session.Request.RemoteAddr)
		_ = session.WriteBinary([]byte("decode nonce fail"))
		_ = session.Close()
		return
	}

	data, err := cert.SRsaDecryptPri(decodeNonce)
	if err != nil {
		logger.Info("decrypt nonce fail", session.Request.RemoteAddr)
		_ = session.WriteBinary([]byte("decrypt nonce fail"))
		_ = session.Close()
		return
	}

	originalNonce := string(data)
	logger.Debug("save key:", saveAuthKey)
	logger.Debug("reqe key:", originalNonce)

	if !(originalNonce == saveAuthKey) {
		logger.Info("auth fail", session.Request.RemoteAddr)
		_ = session.WriteBinary([]byte("auth fail"))
		_ = session.Close()
		return
	}

	_ = session.WriteBinary([]byte("auth success"))
	session.Set("is_auth", true)

}

// 如果 5 秒钟没有通过验证则清除连接
func delayRemoveSession(session *melody.Session) {
	time.Sleep(5 * time.Second)
	value, exists := session.Get("is_auth")
	if !exists {
		logger.Error("auth fail :is_auth not exist", exists)
		_ = session.CloseWithMsg([]byte("auth fail 01 "))
		return
	}

	b, ok := value.(bool)
	if !ok {
		logger.Error("auth fail :is_auth is not bool ", ok)
		_ = session.CloseWithMsg([]byte("auth fail 02 "))
		return
	}

	if !b {
		logger.Error("auth fail :didnt receive auth info or auth err", b)
		_ = session.CloseWithMsg([]byte("auth fail 03 "))
		return
	}

	isClosed := session.IsClosed()

	if isClosed {
		logger.Info("session auth success but is closed: remote addr", session.Request.RemoteAddr)
	} else {
		logger.Info("session auth success : remote addr", session.Request.RemoteAddr)
	}

}
