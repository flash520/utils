/**
 * @Author: koulei
 * @Description: TODO
 * @File:  baseInfo
 * @Version: 1.0.0
 * @Date: 2021/7/5 11:57
 */

package files

import (
	"mime/multipart"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// 返回值说明：
//	7z、exe、doc 类型会返回 application/octet-stream  未知的文件类型
//	jpg	=>	image/jpeg
//	png	=>	image/png
//	ico	=>	image/x-icon
//	bmp	=>	image/bmp
//  xlsx、docx 、zip	=>	application/zip
//  tar.gz	=>	application/x-gzip
//  txt、json、log等文本文件	=>	text/plain; charset=utf-8   备注：就算txt是gbk、ansi编码，也会识别为utf-8

// GetMimeByFilepath 从文件路径获取文件类型
func GetMimeByFilepath(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("打开文件 %s 文件错误, err: %v\n", path, path)
		return ""
	}
	defer func() { _ = f.Close() }()

	buf := make([]byte, 32)
	if _, err = f.Read(buf); err != nil {
		log.Errorf(err.Error())
		return ""
	}

	return http.DetectContentType(buf)
}

// GetMimeByFp 从文件指针获取文件类型
func GetMimeByFp(file multipart.File) string {
	buf := make([]byte, 32)
	if _, err := file.Read(buf); err != nil {
		log.Errorf(err.Error())
		return ""
	}

	return http.DetectContentType(buf)
}
