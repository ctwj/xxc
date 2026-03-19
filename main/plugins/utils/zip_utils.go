package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RepackageZip 重新打包 ZIP 文件
// originalData: 原始 ZIP 数据
// deleteFiles: 要删除的文件列表
// addFiles: 要添加的本地文件列表
// renameRules: 文件名重命名规则（旧名: 新名）
func RepackageZip(originalData []byte, deleteFiles, addFiles, renameRules []string) ([]byte, error) {
	// 读取原始 ZIP
	reader, err := zip.NewReader(bytes.NewReader(originalData), int64(len(originalData)))
	if err != nil {
		return nil, err
	}

	// 创建删除文件的映射
	deleteMap := make(map[string]bool)
	for _, f := range deleteFiles {
		deleteMap[f] = true
	}

	// 创建重命名规则映射
	renameMap := make(map[string]string)
	for _, rule := range renameRules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) == 2 {
			renameMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// 创建新的 ZIP buffer
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// 复制原始文件（跳过要删除的）
	for _, file := range reader.File {
		if deleteMap[file.Name] {
			continue
		}

		// 应用重命名规则
		newName := file.Name
		if replacement, ok := renameMap[file.Name]; ok {
			newName = replacement
		}

		// 打开原文件
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		// 创建新文件
		dst, err := writer.Create(newName)
		if err != nil {
			return nil, err
		}

		// 复制内容
		if _, err := io.Copy(dst, src); err != nil {
			return nil, err
		}
	}

	// 添加新文件
	for _, filePath := range addFiles {
		// 读取本地文件
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// 添加到 ZIP
		dst, err := writer.Create(filepath.Base(filePath))
		if err != nil {
			return nil, err
		}

		if _, err := dst.Write(data); err != nil {
			return nil, err
		}
	}

	// 关闭 writer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ListZipFiles 列出 ZIP 中的所有文件
func ListZipFiles(data []byte) ([]string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range reader.File {
		files = append(files, file.Name)
	}

	return files, nil
}

// ExtractFileFromZip 从 ZIP 中提取指定文件
func ExtractFileFromZip(data []byte, filename string) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, file := range reader.File {
		if file.Name == filename {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, os.ErrNotExist
}

// IsZipFile 检查数据是否为有效的 ZIP 文件
func IsZipFile(data []byte) bool {
	// ZIP 文件的魔数是 0x50 0x4B 0x03 0x04
	if len(data) < 4 {
		return false
	}
	return data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04
}