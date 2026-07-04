package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// HashResult 对齐 Tauri hash::HashResult，字段名小写与前端 hash.ts 一致。
type HashResult struct {
	Size   uint64 `json:"size"`
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
	Sha384 string `json:"sha384"`
	Sha512 string `json:"sha512"`
}

// computeHash 一次性计算所有哈希（用于文本 / 字节）。
func computeHash(data []byte) HashResult {
	md5sum := md5.Sum(data)
	sha1sum := sha1.Sum(data)
	sha256sum := sha256.Sum256(data)
	sha384sum := sha512.Sum384(data)
	sha512sum := sha512.Sum512(data)
	return HashResult{
		Size:   uint64(len(data)),
		Md5:    hex.EncodeToString(md5sum[:]),
		Sha1:   hex.EncodeToString(sha1sum[:]),
		Sha256: hex.EncodeToString(sha256sum[:]),
		Sha384: hex.EncodeToString(sha384sum[:]),
		Sha512: hex.EncodeToString(sha512sum[:]),
	}
}

// HashText 对齐 Tauri hash_text。
func (a *App) HashText(text string) HashResult {
	return computeHash([]byte(text))
}

// HashBytes 对齐 Tauri hash_bytes（前端拖拽读到的字节）。
func (a *App) HashBytes(bytes []byte) HashResult {
	return computeHash(bytes)
}

// HashFile 对齐 Tauri hash_file：64KB 分块流式读取，避免大文件占内存。
func (a *App) HashFile(path string) (HashResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return HashResult{}, err
	}
	defer f.Close()

	var md5h hash.Hash = md5.New()
	var sha1h hash.Hash = sha1.New()
	var sha256h hash.Hash = sha256.New()
	var sha384h hash.Hash = sha512.New384()
	var sha512h hash.Hash = sha512.New()
	hashes := []hash.Hash{md5h, sha1h, sha256h, sha384h, sha512h}

	buf := make([]byte, 64*1024)
	var size uint64
	for {
		n, rerr := f.Read(buf)
		if n > 0 {
			size += uint64(n)
			for _, h := range hashes {
				h.Write(buf[:n])
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return HashResult{}, rerr
		}
	}

	return HashResult{
		Size:   size,
		Md5:    hex.EncodeToString(md5h.Sum(nil)),
		Sha1:   hex.EncodeToString(sha1h.Sum(nil)),
		Sha256: hex.EncodeToString(sha256h.Sum(nil)),
		Sha384: hex.EncodeToString(sha384h.Sum(nil)),
		Sha512: hex.EncodeToString(sha512h.Sum(nil)),
	}, nil
}
