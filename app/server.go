package app

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/olahol/melody.v1"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"ws-server/common"
	"ws-server/handler"
)

var Server struct {
	Port        string
	MongoDBAddr string
	MongoDBUser string
	MongoDBPwd  string
	MongoDBAuth string
}

var (
	globalSession *mgo.Session
)

// Init server with mode
// @mode gin framework mode
func StartServer(mode, wsPath string) {
	gin.SetMode(mode)
	r := gin.New()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./view/index.html")
	})

	//r.GET("/im/ws", func(c *gin.Context) {
	r.GET(wsPath, func(c *gin.Context) {
		err := m.HandleRequestWithKeys(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("HandlerRequest Fail:", err.Error())
		}
	})

	mHandler(m)

	err := r.Run(":" + Server.Port)
	if err != nil {
		logger.Painc("Start server fail:", err.Error())
	}
}

func mHandler(m *melody.Melody) {
	// onMessage
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		handler.OnMessage(m, s, msg)
	})

	// onClose
	m.HandleClose(func(session *melody.Session, i int, s string) error {
		return handler.OnClose(session, i, s)
	})

	// onConn
	m.HandleConnect(func(session *melody.Session) {
		handler.OnConnect(m, session)
	})

	// onDisconnect
	m.HandleDisconnect(func(session *melody.Session) {
		handler.OnDisconnect(m, session)
	})

	// onErr
	m.HandleError(func(session *melody.Session, e error) {
		handler.OnErr(session, e)
	})

	// onPong
	m.HandlePong(func(session *melody.Session) {
		handler.OnPong(session)
	})

	// onSentMessage
	m.HandleSentMessage(func(session *melody.Session, bytes []byte) {
		handler.OnSentMessage(session, bytes)
	})

	// onMessageBinary
	m.HandleMessageBinary(func(session *melody.Session, bytes []byte) {
		handler.OnMessageBinary(session, bytes)
	})

	// onSentMessageBinary
	m.HandleSentMessageBinary(func(session *melody.Session, bytes []byte) {
		handler.OnSentMessageBinary(session, bytes)
	})

}

// 初始化
func init() {
	// initLogger()
	initServerConf()
	initMongoDb()
}

// 加载日志适配器
func initLogger() {
	//
	bytes, err := ioutil.ReadFile("conf/log.json")
	if err != nil {
		panic("[server/init.go:115] 加载日志配置失败:" + err.Error())
		return
	}

	confJson, err := simplejson.NewJson(bytes)
	if err != nil {
		panic("[server/init.go:121] 加载日志配置失败(not json):" + err.Error())
		return
	}

	fileFullName, err := confJson.Get("File").Get("filename").String()

	suffix := strings.LastIndex(fileFullName, `/`)
	if suffix != -1 {
		filePath := fileFullName[:suffix]
		err := common.CheckAndMakePath(filePath)
		if err != nil {
			panic("[server/init.go:121] 加载日志配置失败(自动创建日志文件夹失败):" + err.Error())
			return
		}
	}

	err = logger.SetLogger("conf/log.json")
	//err := logger.SetLogger("D:/work/im-version1/im-server/conf/log.json")
	if err != nil {
		panic("[server/init.go] 加载日志配置失败:" + err.Error())
		return
	}
	logger.Info("日志加载成功...")
}

// 加载服务器配置
func initServerConf() {
	data, err := ioutil.ReadFile("conf/server.json")

	//data, err := ioutil.ReadFile("D:/work/im-version1/im-server/conf/server.json")
	if err != nil {
		logger.Painc("加载服务器配置失败!", err.Error())
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		logger.Painc("加载服务器配置失败!", err.Error())
	}
	logger.Info("加载服务器配置成功,正在启动服务...")
}

// 初始化mongoDB
func initMongoDb() {
	var err error
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{Server.MongoDBAddr},
		Timeout:  60 * time.Second,
		Database: Server.MongoDBAuth,
		Username: Server.MongoDBUser,
		Password: Server.MongoDBPwd,
	}
	globalSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Painc("Mongodb 连接失败:", err)
	}
	globalSession.SetMode(mgo.Monotonic, true)
	logger.Info("连接mongo数据库成功 address:", Server.MongoDBAddr)
}

func GetDBConn(dbName, cName string) (*mgo.Session, *mgo.Collection) {
	s := globalSession.Copy()
	c := s.DB(dbName).C(cName)
	return s, c
}
