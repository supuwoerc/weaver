package utils

import "os"

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) (bool, error) {
	dir, err := IsDir(path)
	return !dir, err
}
