package common

import (
	"github.com/wonderivan/logger"
	"os"
)

// 删除一个文件 或者一个文件夹下的所有文件
func DeleteAllFile(path string) error {
	err := os.RemoveAll(path) //删除文件test.txt // 或者删除文件夹及该文件夹下的所有文件
	if err != nil {
		//如果删除失败则输出 file remove Error!
		//输出错误详细信息
		logger.Error("删除文件失败:", err.Error())
	} else {
		//如果删除成功则输出 file remove OK!
		logger.Info(":file remove OK:", path)
	}
	return err

}

// 检查文件目录是否存在 不存在则创建
func CheckAndMakePath(path string) error {
	b, err := IsPathExist(path)
	if err != nil {
		// log
		logger.Error("判断文件失败：", err.Error())
		return err
	}

	if !b {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			// log
			logger.Error("创建文件失败：", err.Error())
			return err
		}
	}
	return nil
}

// 判断文件路径是否存在
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}