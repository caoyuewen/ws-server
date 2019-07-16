package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/wonderivan/logger"
	"io/ioutil"
)

// 两对秘钥, 服务器用服务器私钥接收客户端用服务器公钥加密的数据
// 客户端用客户端私钥接收服务器端用客户端公钥加密的数据
// 服务器端只用到了 服务器私钥 和 客户端公钥
// 客户端只用到了 服务器公钥 和 客户端私钥
const (
	//D:/fido/myframework/ws-server/cert/client_pri.pem
	cPubKey = "D:/fido/myframework/ws-server/cert/client_pub.pem"
	sPriKey = "D:/fido/myframework/ws-server/cert/server_pri.pem"
	sPubKey = "D:/fido/myframework/ws-server/cert/server_pub.pem"
	cPriKey = "D:/fido/myframework/ws-server/cert/client_pri.pem"

)

var cPubPerm []byte
var sPriPerm []byte

var sPubPerm []byte
var cPriPerm []byte

func init() {
	var err error
	sPubPerm, err = ioutil.ReadFile(sPubKey)
	if err != nil {
		logger.Error("app rsa priKey fail!")
		panic(err.Error())
	}

	cPriPerm, err = ioutil.ReadFile(cPriKey)
	if err != nil {
		logger.Error("app rsa pubKey fail!")
		panic(err.Error())
	}

	sPriPerm, err = ioutil.ReadFile(sPriKey)
	if err != nil {
		logger.Error("app rsa priKey fail!")
		panic(err.Error())
	}

	cPubPerm, err = ioutil.ReadFile(cPubKey)
	if err != nil {
		logger.Error("app rsa pubKey fail!")
		panic(err.Error())
	}

	//cPubPerm = []byte(publicKey)
	//sPriPerm = []byte(privateKey)

}

// 使用 client公钥加密
func CRsaEncryptPub(origData []byte) (data []byte, err error) {
	//pubKey, err := ioutil.ReadFile(cPubKey)
	//
	////pubKey, err := ioutil.ReadFile(cPubKey)
	//if err != nil {
	//	return
	//}

	//解密pem格式的公钥
	block, _ := pem.Decode(cPubPerm)
	if block == nil {
		return nil, errors.New("public key error")
	}

	//解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	//pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//类型断言
	pub := pubInterface.(*rsa.PublicKey)

	data, err = rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

	return
}

// 使用 client私钥解密 PKCS1
func CRsaDecryptPri(cipherText []byte) (data []byte, err error) {
	//priKey, err := ioutil.ReadFile(sPriKey)
	//if err != nil {
	//	return
	//}

	//解密
	block, _ := pem.Decode(cPriPerm)
	if block == nil {
		return
	}

	//解析PKCS1格式的私钥
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	data, err = rsa.DecryptPKCS1v15(rand.Reader, key, cipherText)
	return
}



// 使用 server私钥解密 PKCS1
func SRsaDecryptPri(cipherText []byte) (data []byte, err error) {
	//priKey, err := ioutil.ReadFile(sPriKey)
	//if err != nil {
	//	return
	//}

	//解密
	block, _ := pem.Decode(sPriPerm)
	if block == nil {
		return
	}

	//解析PKCS1格式的私钥
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	data, err = rsa.DecryptPKCS1v15(rand.Reader, key, cipherText)
	return
}


// 使用 server公钥加密
func SRsaEncryptPub(origData []byte) (data []byte, err error) {
	//pubKey, err := ioutil.ReadFile(cPubKey)
	//
	////pubKey, err := ioutil.ReadFile(cPubKey)
	//if err != nil {
	//	return
	//}

	//解密pem格式的公钥
	block, _ := pem.Decode(sPubPerm)
	if block == nil {
		return nil, errors.New("public key error")
	}

	//解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	//pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//类型断言
	pub := pubInterface.(*rsa.PublicKey)

	data, err = rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

	return
}



// 获取一个10000-99999的随机数字符串
func GetSNonce(sNonce string) (ensNonce string, err error) {
	//sNonce := strconv.Itoa(RndNum(100000, 999999))
	//sNonce := uuid.New().String()

	//bytes, err := base64.StdEncoding.DecodeString(sNonce)
	//if err != nil {
	//	return
	//}
	data, err := CRsaEncryptPub([]byte(sNonce))
	if err != nil {
		return
	}
	ensNonce = base64.StdEncoding.EncodeToString(data)
	return
}



