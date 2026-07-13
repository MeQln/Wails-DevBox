package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	xqrcode "github.com/tuotoo/qrcode"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// QrEncode 对齐 Tauri qrcode::qr_encode：返回 SVG markup。
// 空文本返回 "文本为空"；过长返回 "文本过长，无法生成"。
// 用 skip2/go-qrcode 取模块矩阵，手写 SVG（黑模块 1x1 rect，矢量，前端可缩放/光栅化）。
func (a *App) QrEncode(text string) (string, error) {
	if text == "" {
		return "", errors.New("文本为空")
	}
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		// skip2 对内容过长返回 error，统一文案
		return "", errors.New("文本过长，无法生成")
	}
	return qrToSVG(qr.Bitmap()), nil
}

// qrToSVG 把含 quiet zone 的模块矩阵渲染为 SVG markup。
// 不设置固定颜色：背景透明，模块使用 currentColor 以继承主题文本色。
func qrToSVG(bitmap [][]bool) string {
	n := len(bitmap)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges">`, n, n, n, n)
	for i, row := range bitmap {
		for j, v := range row {
			if v {
				fmt.Fprintf(&b, `<rect x="%d" y="%d" width="1" height="1" fill="currentColor"/>`, j, i)
			}
		}
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// QrDecode 对齐 Tauri qrcode::qr_decode：解码图片字节，返回二维码内容。
// 图片解析失败返回 "图片解析失败"；未识别到二维码返回 "未识别到二维码"。
// tuotoo/qrcode 的 Decode 接收 io.Reader（内部自行 image.Decode），故先用一次
// image.Decode 验证图片合法性以区分两种错误，再传新 reader 给解码器。
func (a *App) QrDecode(data []byte) (string, error) {
	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		return "", errors.New("图片解析失败")
	}
	qr, err := xqrcode.Decode(bytes.NewReader(data))
	if err != nil {
		return "", errors.New("未识别到二维码")
	}
	return qr.Content, nil
}
