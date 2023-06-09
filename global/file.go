package global

import (
	"path"
	"strings"
)

// IsImage 判断文件是否为图片
func IsImage(filename string) bool {
	// 判断文件是否为图片
	isImage := false
	ext := path.Ext(filename)
	ext = strings.ToLower(ext)
	switch ext {
	case ".bmp", ".png", ".jpg", ".jpeg", ".gif":
		isImage = true
	}
	return isImage
}

// IsScript 判断文件是否为脚本
func IsScript(filename string) bool {
	// 判断文件是否为图片
	isImage := false
	ext := path.Ext(filename)
	ext = strings.ToLower(ext)
	switch ext {
	case ".js", ".php", ".sh", ".shell", ".py", ".rs":
		isImage = true
	}
	return isImage
}

// IsAllowFile 判断文件是否为允许的文件
func IsAllowFile(filename string) bool {
	isImage := false
	ext := path.Ext(filename)
	ext = strings.ToLower(ext)
	switch ext {
	case ".bmp", ".png", ".jpg", ".jpeg", ".gif", ".txt", ".ico", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".doc", ".docx":
		isImage = true
	}
	return isImage
}
