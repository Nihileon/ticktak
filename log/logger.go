package log

import (
    "github.com/Sirupsen/logrus"
    "os"
    "time"
)

var logger *logrus.Logger

func GetLogger() *logrus.Logger {
    return logger
}

func Init() {
    file := "./log/" + time.Now().Format("2006-01-02") + ".log "
    logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
    if err!=nil && os.IsNotExist(err){
        logFile, _ = os.Create(file)
    }
    if nil != err {
        panic(err)
    }
    logger = logrus.StandardLogger()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetOutput(logFile)
    logger.SetLevel(logrus.InfoLevel)
}
