package fantia

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var ForbiddenItemMap map[string]string = map[string]string{
	`\`:   `￥`,
	`/`:   `／`,
	`:`:   `：`,
	`*`:   `＊`,
	`?`:   `？`,
	`"`:   `”`,
	`<`:   `＜`,
	`>`:   `＞`,
	`|`:   `｜`,
	`...`: `…`,
}

var bs *regexp.Regexp = regexp.MustCompile("\b")

func ForbiddenTextRename(name string) string {
	if bs.MatchString(name) {
		name = bs.ReplaceAllString(name, "")
	}

	for key, val := range ForbiddenItemMap {
		if !strings.Contains(name, key) {
			continue
		}
		name = strings.ReplaceAll(name, key, val)
	}
	return name
}

func isFolderExist(dir string, make bool) (bool, error) {
	var exist bool
	_, err := os.Stat(dir)
	exist = !os.IsNotExist(err)
	if !make {
		return exist, err
	}

	if exist {
		return exist, nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, err
	}
	return true, nil
}

func removeDuplicateInt(source []int) []int {
	results := make([]int, 0, len(source))
	encountered := map[int]bool{}
	for i := 0; i < len(source); i++ {
		if !encountered[source[i]] {
			encountered[source[i]] = true
			results = append(results, source[i])
		}
	}
	return results
}

func RetrySaveURIToFile(client *http.Client, file, uri string) error {
	var err error
	for i, interval := 1, retryInterval; i <= retry; i++ {
		err = SaveURIToFile(client, file, uri)
		if err == nil {
			return nil
		}
		log.Printf("Retry wait time is %d seconds. error to %v. ", interval, err)

		time.Sleep(time.Duration(interval) * time.Second)

		if interval*2 > 300 {
			interval = 300
		} else {
			interval *= 2
		}
	}
	return err
}

func SaveURIToFile(client *http.Client, file, uri string) (err error) {
	if _, err := os.Stat(file); err == nil {
		return nil
	}

	dir, _ := filepath.Split(file)
	_, err = isFolderExist(dir, true)
	if err != nil {
		return err
	}

	resp, err := client.Get(uri)
	if err != nil {
		return err
	}
	defer HTTPResponseBodyCloser(resp)

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP Response status not successed: %s\n", resp.Status)
		log.Printf("=> Save %s from %s\n", file, uri)
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer CloseFile(f, err)

	bw := bufio.NewWriter(f)
	if _, err := bw.ReadFrom(resp.Body); err != nil {
		return err
	}
	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func CloseFile(f *os.File, err error) {
	err = f.Close()
}

func ConvertToJpeg(src string, dest string) (err error) {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer CloseFile(file, err)

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer CloseFile(out, err)

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}

func ConvertToJpegFromPng(src string, dest string) (err error) {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer CloseFile(file, err)

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer CloseFile(out, err)

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}

func ConvertToJpegFromGif(src string, dest string) (err error) {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer CloseFile(file, err)

	img, err := gif.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer CloseFile(out, err)

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}

func DetectFileContentType(file string) string {
	var err error
	f, _ := os.Open(file)
	defer CloseFile(f, err)

	buffer := make([]byte, 512)
	f.Read(buffer)

	contentType := http.DetectContentType(buffer)

	f.Seek(0, 0)
	f.Fd()

	return contentType
}

func MakeFolderIcon(client *http.Client, path, uri string) {
	icon_path := filepath.Join(path, FOLDER_JPEG)
	ext := IconExt.ReplaceAllString(uri, `$1`)

	if _, err := os.Stat(icon_path); err == nil {
		return
	}
	file := filepath.Join(path, fmt.Sprintf("folder.%s", ext))
	RetrySaveURIToFile(client, file, uri)

	img := DetectFileContentType(file)

	logs := func(img, file, icon, uri string, err error) {
		if err != nil {
			log.Printf("URI : %s\n", uri)
			log.Printf("IMG : %s\n", img)
			log.Printf("FILE: %s\n", file)
			log.Printf("ICON: %s\n", icon_path)
			log.Printf("ERR : %s\n", err)
		}
	}

	var err error
	switch {
	case img == IMAGE_JPEG && icon_path == file:
	case img == IMAGE_JPEG && icon_path != file:
		err = os.Rename(file, icon_path)
		logs(img, file, icon_path, uri, err)
	case img == IMAGE_PNG:
		err = ConvertToJpegFromPng(file, icon_path)
		logs(img, file, icon_path, uri, err)

		err = os.Remove(file)
	case img == IMAGE_GIF:
		err = ConvertToJpegFromGif(file, icon_path)
		logs(img, file, icon_path, uri, err)

		err = os.Remove(file)
	default:
		err = ConvertToJpeg(file, icon_path)
		logs(img, file, icon_path, uri, err)

		err = os.Remove(file)
	}
	if err != nil {
		AddErrorLog(fmt.Sprintf("%s: %v", time.Now(), err))
	}
}

func HTTPResponseBodyCloser(resp *http.Response) {
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}

func DirectoryChecker(parent, find string) bool {
	list, err := os.ReadDir(parent)
	if err != nil {
		return false
	}

	_, last := filepath.Split(find)
	for _, v := range list {
		if v.Name() == last {
			return true
		}
	}
	return false
}
