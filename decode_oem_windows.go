//go:build windows

package main

import (
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

var procGetOEMCP = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetOEMCP")

// getOEMCP 调用 Windows API GetOEMCP 获取系统 OEM 代码页编号。
func getOEMCP() uint32 {
	ret, _, _ := procGetOEMCP.Call()
	return uint32(ret)
}

// decodeOEM 在 Windows 上通过 GetOEMCP 获取系统 OEM 代码页，
// 将非 UTF-8 的原始字节转为 UTF-8，解决 ping 等系统命令输出乱码问题。
//
// GetOEMCP 返回当前系统 OEM 代码页编号：
//   936 → GBK（中文简体）
//   932 → ShiftJIS（日语）
//   949 → EUC-KR（韩语）
//   950 → Big5（中文繁体）
//   437 → IBM-437（英语美国）
//   850 → MS-DOS Latin-1（西欧）
//   866 → Cyrillic（俄语）
//   ……
func decodeOEM(raw []byte) []byte {
	cp := getOEMCP()
	dec := decoderForCP(cp)
	if dec == nil {
		return raw
	}
	utf8Bytes, err := dec.Bytes(raw)
	if err != nil {
		return raw
	}
	return utf8Bytes
}

// decoderForCP 返回 OEM 代码页对应的 UTF-8 decoder。
// 不支持的代码页返回 nil，由调用方处理回退。
func decoderForCP(cp uint32) *encoding.Decoder {
	switch cp {
	case 437:
		return charmap.CodePage437.NewDecoder()
	case 850:
		return charmap.CodePage850.NewDecoder()
	case 852:
		return charmap.CodePage852.NewDecoder()
	case 855:
		return charmap.CodePage855.NewDecoder()
	case 858:
		return charmap.CodePage858.NewDecoder()
	case 860:
		return charmap.CodePage860.NewDecoder()
	case 862:
		return charmap.CodePage862.NewDecoder()
	case 863:
		return charmap.CodePage863.NewDecoder()
	case 865:
		return charmap.CodePage865.NewDecoder()
	case 866:
		return charmap.CodePage866.NewDecoder()
	case 932:
		return japanese.ShiftJIS.NewDecoder()
	case 936:
		return simplifiedchinese.GBK.NewDecoder()
	case 949:
		return korean.EUCKR.NewDecoder()
	case 950:
		return traditionalchinese.Big5.NewDecoder()
	default:
		// 65001 (UTF-8) 及未映射代码页由 toUTF8 的 UTF-8 验证和 GBK 兜底处理
		return nil
	}
}