package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"github.com/foobaz/lossypng/lossypng"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// compress string by zlib
func CompressStrByZlib(s string) (string, error) {
	var b bytes.Buffer
	gz := zlib.NewWriter(&b)
	if _, err := gz.Write([]byte(s)); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	res := base64.StdEncoding.EncodeToString(b.Bytes())
	return res, nil
}

// decompression string by zlib
func DeCompressStrByZlib(s string) string {
	data, _ := base64.StdEncoding.DecodeString(s)
	rData := bytes.NewReader(data)
	r, _ := zlib.NewReader(rData)
	b, _ := ioutil.ReadAll(r)
	return string(b)
}

// compress image
func CompressImage(filename string) error {
	return CompressImageSaveOriginal(filename, "")
}

// compress image, save original file to before dif, If the front is empty, it will not be saved
func CompressImageSaveOriginal(filename string, before string) error {
	suffix := strings.ToLower(filepath.Ext(filename))
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		return fmt.Errorf("picture format is not supported: %s", filename)
	}
	isJpg := true
	if suffix == ".png" {
		isJpg = false
	}
	// tmp filename
	newFilename := filename + ".compress"
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot find file %s, err: %v", filename, err)
	}

	beforeFilename := ""
	beforeDir := ""
	if before != "" {
		baseDir, name := filepath.Split(filename)
		if strings.Contains(filename, before) || strings.Contains(baseDir, before) {
			// 当前目录是源文件目录
			return nil
		}
		beforeDir = baseDir + before
		beforeFilename = beforeDir + "/" + name
		_, err := os.Stat(beforeFilename)
		if err == nil {
			return fmt.Errorf("file %s has been compressed, so it will not be compressed again", filename)
		}
	}

	// decode image file
	var img image.Image
	if isJpg {
		img, err = jpeg.Decode(file)
	} else {
		img, err = png.Decode(file)
	}
	if err != nil {
		return fmt.Errorf("decode image err: %v", err)
	}
	file.Close()
	bound := img.Bounds()
	width := bound.Dx()
	height := bound.Dy()
	var compressed image.Image
	if isJpg {
		// compress jpg(Lanczos2)
		compressed = resize.Resize(uint(width), uint(height), img, resize.MitchellNetravali)
	} else {
		// compress png(the quality is 20% of the original)
		compressed = lossypng.Compress(img, lossypng.NoConversion, 20)
	}

	out, err := os.Create(newFilename)
	if err != nil {
		return fmt.Errorf("create tmp file %s, err: %v", newFilename, err)
	}
	defer out.Close()

	// encode image file
	if isJpg {
		err = jpeg.Encode(out, compressed, &jpeg.Options{Quality: 40})
	} else {
		err = png.Encode(out, compressed)
	}
	if err != nil {
		return fmt.Errorf("encode image err: %v", err)
	}
	if beforeDir != "" {
		CreateDirIfNotExists(beforeDir)
		err = os.Rename(filename, beforeFilename)
		if err != nil {
			return fmt.Errorf("save original file to %s err: %v", beforeFilename, err)
		}
	}
	err = os.Rename(newFilename, filename)
	if err != nil {
		return fmt.Errorf("rename %s to %s, err: %v", newFilename, filename, err)
	}
	return nil
}
