package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ImageInfo 返回给前端的图片信息
type ImageInfo struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	SizeBytes int64  `json:"size_bytes"`
	DataBase64 string `json:"data_base64"`
}

// ImageFmt 目标格式枚举
type ImageFmt string

const (
	FmtPNG  ImageFmt = "png"
	FmtJPEG ImageFmt = "jpeg"
	FmtWebP ImageFmt = "webp"
	FmtBMP  ImageFmt = "bmp"
	FmtGIF  ImageFmt = "gif"
	FmtTIFF ImageFmt = "tiff"
)

func parseImageFormat(s string) (ImageFmt, bool) {
	switch strings.ToLower(s) {
	case "png":
		return FmtPNG, true
	case "jpeg", "jpg":
		return FmtJPEG, true
	case "webp":
		return FmtWebP, true
	case "bmp":
		return FmtBMP, true
	case "gif":
		return FmtGIF, true
	case "tiff", "tif":
		return FmtTIFF, true
	}
	return "", false
}

// decodeImage 打开文件并解码图片，注册所有支持的格式解码器
func decodeImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 注册 bmp/tiff 解码器后，image.Decode 能识别这些格式
	img, fmtName, err := image.Decode(f)
	if err != nil {
		return nil, "", fmt.Errorf("解码图片失败: %w", err)
	}
	return img, fmtName, nil
}

// encodeImage 将图片编码为指定格式，写入 writer
func encodeImage(w io.Writer, img image.Image, imgFmt ImageFmt, quality int) error {
	switch imgFmt {
	case FmtPNG:
		return png.Encode(w, img)
	case FmtJPEG:
		q := quality
		if q < 1 {
			q = 85
		}
		if q > 100 {
			q = 100
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: q})
	case FmtGIF:
		return gif.Encode(w, img, nil)
	case FmtBMP:
		return bmp.Encode(w, img)
	case FmtTIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate})
	default:
		return fmt.Errorf("不支持的输出格式: %s", imgFmt)
	}
}

// readImageInfo 读取图片并返回 ImageInfo（通用函数）
func readImageInfo(sourcePath string) (*ImageInfo, error) {
	path := sourcePath
	img, decodedFmt, err := decodeImage(path)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 读取原始文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 从扩展名推断格式，回退到解码器识别的格式
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if _, ok := parseImageFormat(ext); !ok {
		ext = decodedFmt
	}

	info := &ImageInfo{
		Width:      width,
		Height:     height,
		Format:     ext,
		SizeBytes:  int64(len(data)),
		DataBase64: base64.StdEncoding.EncodeToString(data),
	}
	return info, nil
}

// ImageRead 读取图片并返回基本信息 + base64 数据（用于前端预览）
func (a *App) ImageRead(sourcePath string) (*ImageInfo, error) {
	return readImageInfo(sourcePath)
}

// ImageConvert 转换图片格式
// sourcePath: 源文件路径, targetFmt: 目标格式, outputPath: 输出路径
func (a *App) ImageConvert(sourcePath string, targetFmt string, outputPath string) (*ImageInfo, error) {
	fmtVal, ok := parseImageFormat(targetFmt)
	if !ok {
		return nil, fmt.Errorf("不支持的格式: %s", targetFmt)
	}

	// WebP 编码暂不支持，返回明确错误
	if fmtVal == FmtWebP {
		return nil, fmt.Errorf("WebP 编码暂不支持，请选择 PNG 或 JPEG 格式")
	}

	img, _, err := decodeImage(sourcePath)
	if err != nil {
		return nil, err
	}

	// 防止源文件被覆盖：若输出路径等于源路径，先写入临时文件再替换
	var needTemp bool
	var tempPath string
	absSrc, _ := filepath.Abs(sourcePath)
	absOut, _ := filepath.Abs(outputPath)
	if absSrc == absOut {
		needTemp = true
		tempPath = outputPath + ".tmp"
		outputPath = tempPath
	}

	// 写入输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	if err := encodeImage(outFile, img, fmtVal, 85); err != nil {
		return nil, fmt.Errorf("转换失败: %w", err)
	}
	outFile.Close()

	// 若使用了临时文件，替换回原路径
	if needTemp {
		if err := os.Rename(tempPath, absSrc); err != nil {
			return nil, fmt.Errorf("替换源文件失败: %w", err)
		}
		outputPath = absSrc
	}

	return readImageInfo(outputPath)
}

// ImageCompress 压缩图片
// sourcePath: 源文件路径, quality: 质量 (1-100), outputPath: 输出路径
func (a *App) ImageCompress(sourcePath string, quality int, outputPath string) (*ImageInfo, error) {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}

	img, _, err := decodeImage(sourcePath)
	if err != nil {
		return nil, err
	}

	// 根据输出文件扩展名决定格式
	ext := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	fmtVal, ok := parseImageFormat(ext)
	if !ok {
		fmtVal = FmtJPEG
	}

	// 对于 WebP 格式，Go 标准库不支持编码，降级为 PNG
	if fmtVal == FmtWebP {
		// 修正扩展名为 .png
		outputPath = strings.TrimSuffix(outputPath, "."+ext) + ".png"
		fmtVal = FmtPNG
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	if err := encodeImage(outFile, img, fmtVal, quality); err != nil {
		return nil, fmt.Errorf("压缩失败: %w", err)
	}
	outFile.Close()

	return readImageInfo(outputPath)
}

// 注册 bmp/tiff 解码器（init 自动注册 webp）
func init() {
	image.RegisterFormat("bmp", "BM", bmp.Decode, bmp.DecodeConfig)
	image.RegisterFormat("tiff", "II*\x00", tiff.Decode, tiff.DecodeConfig)
	image.RegisterFormat("tiff", "MM\x00*", tiff.Decode, tiff.DecodeConfig)
}