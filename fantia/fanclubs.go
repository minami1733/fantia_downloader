package fantia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewConfig() *Config {
	return new(Config)
}

func (c *Config) SetOutput(output string) {
	c.Output = &output
}

func NewFantiaDownloader(cfg Config) (*FantiaDownloader, error) {
	return &FantiaDownloader{}, nil
}

var output string
var overwrite bool
var progress bool
var waittime time.Duration
var date *time.Time = nil
var retry int = 5
var retryInterval int

func SetOutput(dir string) {
	output = dir
}
func SetOverwrite(opt bool) {
	overwrite = opt
}
func SetProgress(opt bool) {
	progress = opt
}
func SetDate(opt string) {
	if opt == "" {
		return
	}

	parse, err := time.Parse("2006-01-02", opt)
	if err != nil {
		return
	}
	date = &parse
}
func SetWaitTime(opt int) {
	waittime = time.Millisecond * time.Duration(opt)
}

func SetRetry(opt int) {
	retry = opt
}
func SetRetryInterval(opt int) {
	retryInterval = opt
}

func GetFanclubs(client *http.Client) (*fanclubs, error) {
	resp, err := client.Get(FANTIA_API_FANCLUBS)
	if err != nil {
		return nil, err
	}
	defer HTTPResponseBodyCloser(resp)

	fanclubs := new(fanclubs)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(fanclubs); err != nil {
		return nil, err
	}

	sort.Ints(fanclubs.FanclubIDs)

	return fanclubs, nil
}

func GetFanclub(client *http.Client, fanclub int, exit bool) (string, error) {
	resp, err := client.Get(fmt.Sprintf(FANTIA_API_FANCLUB_INFO, fanclub))
	if err != nil {
		return "", err
	}
	defer HTTPResponseBodyCloser(resp)

	fanclub_data := new(FanclubData)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(fanclub_data); err != nil {
		return "", err
	}

	// 無料プランの判定
	free := true
	for _, plan := range fanclub_data.Fanclub.Plans {
		if plan.Price > 0 && plan.Order.Status == PlanJoined {
			free = !free
		}
	}

	// 無料プラン＆無料プランを省く場合終了
	if exit && free {
		return "", nil
	}

	if free {
		fmt.Printf(" FanclubID: % 10d, Name: %s\n", fanclub, fanclub_data.Fanclub.FanclubNameWithCreatorName)
	} else {
		fmt.Printf("*FanclubID: % 10d, Name: %s\n", fanclub, fanclub_data.Fanclub.FanclubNameWithCreatorName)
	}

	// ファンクラブ名に利用不能な文字がある場合に変換
	name := ForbiddenTextRename(fanclub_data.Fanclub.FanclubNameWithCreatorName)

	// ファンクラブのパス作成(./{output}/{fanclub_name})
	fanclub_path := filepath.Join(output, name)
	_, err = isFolderExist(fanclub_path, true)
	if err != nil {
		return "", err
	}

	// ファンクラブのサムネイル作成
	if icon := fanclub_data.Fanclub.Icon.Original; icon != "" {
		if !strings.Contains(icon, "default") {
			MakeFolderIcon(client, fanclub_path, icon)
		}
	}

	return name, nil
}

func GetFanclubPage(client *http.Client, fanclub int) []int {
	var results []int
	for page := 1; ; page++ {
		posts := func() []string {
			resp, err := client.Get(fmt.Sprintf(FANTIA_FANCLUB_POSTS, fanclub, page))
			if err != nil {
				panic(err)
			}
			defer HTTPResponseBodyCloser(resp)

			if resp.StatusCode != 200 {
				return nil
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			posts := Posts.FindAllStringSubmatch(string(data), -1)
			if len(posts) == 0 {
				return nil
			}

			var posts_data []string
			for _, ids := range posts {
				if len(ids) == 2 {
					posts_data = append(posts_data, ids[1])
				}
			}

			return posts_data
		}()

		if len(posts) == 0 {
			break
		}

		for _, post_id := range posts {
			id, err := strconv.Atoi(post_id)
			if err != nil {
				return nil
			}
			results = append(results, id)
		}
	}

	sort.Ints(results)

	reverse := removeDuplicateInt(results)

	for i := 0; i < len(reverse)/2; i++ {
		reverse[i], reverse[len(reverse)-i-1] = reverse[len(reverse)-i-1], reverse[i]
	}

	return reverse
}

func GetPost(client *http.Client, parent string, post_id int) (bool, error) {
	// 投稿を取得
	resp, err := client.Get(fmt.Sprintf(FANTIA_API_POST_INFO, post_id))
	if err != nil {
		return false, err
	}
	defer HTTPResponseBodyCloser(resp)

	post := new(PostData)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(post); err != nil {
		return false, err
	}

	// 投稿内の投稿日時を取得
	post_date, err := time.Parse(time.RFC1123Z, post.Post.PostedAt)
	if err != nil {
		return false, err
	}

	if date != nil {
		if !post_date.After(*date) {
			return true, nil
		}
	}

	log.Printf("% 4s%s - %s\n", RIGHT_ARROWS, post_date.Format(LOG_FORMAT_DATE), post.Post.Title)

	// 投稿のタイトルからDirectoryに利用できない文字を置換
	post_title := ForbiddenTextRename(post.Post.Title)

	// 投稿のDirectory名を{yyyy-mm-dd_hhmmss_{POST_TITLE}}に変換
	post_dir_title := fmt.Sprintf("%s-%s", post_date.Format(POST_DIR_FORMAT), post_title)

	// 投稿名のDirectory作成 {output}/{fanclub_name}/{yyyy-mm-dd_hhmmss_{POST_TITLE}}
	post_root_dir := CutStringToLimit(filepath.Join(output, parent, post_dir_title), 230)

	for retry := 0; retry < 3; retry++ {
		if _, err := isFolderExist(post_root_dir, true); err != nil {
			panic(err)
		}
		if DirectoryChecker(filepath.Join(output, parent), post_root_dir) {
			break
		}
	}

	if thumb := post.Post.Thumb.Original; thumb != "" {
		MakeFolderIcon(client, post_root_dir, thumb)
	}

	for idx, post_content := range post.Post.PostContents {
		if post_content.Title == "" {
			post_content.Title = "No Title"
		}
		// セクション名を加工
		post_content_title := fmt.Sprintf("%03d-%s", idx+1, ForbiddenTextRename(post_content.Title))

		// 最終的なディレクトリ名を作成 {output}/{fanclub_name}/{yyyy-mm-dd_hhmmss_{POST_TITLE}}/{001-contents_name}
		post_content_root_dir := filepath.Join(post_root_dir, post_content_title)

		log.Printf("% 6s%s\n", RIGHT_ARROWS, post_content_title)

		post_content_root_dir = CutStringToLimit(post_content_root_dir, 245)
		for retry := 0; retry < 3; retry++ {
			if _, err := isFolderExist(post_content_root_dir, true); err != nil {
				panic(err)
			}
			if DirectoryChecker(post_root_dir, post_content_root_dir) {
				break
			}
		}

		if len(post_content.PostContentPhotos) > 0 {
			for idx, photo := range post_content.PostContentPhotos {
				ext := PostItemExt.ReplaceAllString(photo.URL.Original, `$1`)
				fname := fmt.Sprintf("%03d.%s", idx+1, ext)
				path := filepath.Join(post_content_root_dir, fname)

				if err := RetrySaveURIToFile(client, path, photo.URL.Original); err != nil {
					panic(err)
				}

				if progress {
					log.Printf("% 8s%s\n", RIGHT_ARROWS, fname)
				}
			}
		}

		if post_content.Filename != "" && post_content.DownloadURI != "" {
			path := filepath.Join(post_content_root_dir, post_content.Filename)
			RetrySaveURIToFile(client, path, FANTIA_BASE_URI+post_content.DownloadURI)

			if progress {
				log.Printf("% 8s%s\n", RIGHT_ARROWS, post_content.Filename)
			}
		}

		time.Sleep(waittime)
	}
	return false, nil
}
