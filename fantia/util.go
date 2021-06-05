package fantia

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var ForbiddenItemMap map[string]string = map[string]string{
	`\`: `￥`,
	`/`: `／`,
	`:`: `：`,
	`*`: `＊`,
	`?`: `？`,
	`"`: `”`,
	`<`: `＜`,
	`>`: `＞`,
	`|`: `｜`,
}

func ForbiddenTextRename(name string) string {
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

func CutStringToLimit(input string, limit uint) string {
	if len(input) < int(limit) {
		return strings.Trim(input, " ")
	}

	runes := []rune(input)

	for idx := range runes {
		if len(string(runes[0:idx])) > int(limit) {
			input = string(runes[0 : idx-1])
			break
		}
	}
	return strings.Trim(input, " ")
}

func SaveURIToFile(client *http.Client, file, uri string) error {
	if _, err := os.Stat(file); err == nil {
		return nil
	}

	dir, _ := filepath.Split(file)
	_, err := isFolderExist(dir, true)
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

	bw := bufio.NewWriter(f)
	if _, err := bw.ReadFrom(resp.Body); err != nil {
		return err
	}
	if err := bw.Flush(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func ConvertToJpeg(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}

func ConvertToJpegFromPng(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}
func ConvertToJpegFromGif(src string, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := gif.Decode(file)
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: 100}

	return jpeg.Encode(out, img, opts)
}

func DetectFileContentType(file string) string {
	f, _ := os.Open(file)
	defer f.Close()

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
	SaveURIToFile(client, file, uri)

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
		panic(err)
	}
}

func HTTPResponseBodyCloser(resp *http.Response) {
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}

func DirectoryChecker(parent, find string) bool {
	list, err := ioutil.ReadDir(parent)
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
