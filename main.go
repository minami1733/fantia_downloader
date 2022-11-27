package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"

	"fantia_downloader/fantia"
)

var (
	FantiaSessionID *string
	FanclubID       *string
	DownloadDir     *string
	Date            *string
	ExcludeFreePlan *bool
	Overwrite       *bool
	Progress        *bool
	WaitTime        *int
)

func init() {
	// 引数の設定
	FantiaSessionID = flag.String("key", "", "")
	FanclubID = flag.String("fanclub", "", "")
	DownloadDir = flag.String("output", "./downloads", "")
	Date = flag.String("date", "", "")
	ExcludeFreePlan = flag.Bool("excludeFreePlan", false, "")
	Overwrite = flag.Bool("overwrite", false, "")
	Progress = flag.Bool("progress", false, "")
	WaitTime = flag.Int("wait", 250, "")
	flag.Parse()
}

func main() {
	// ログ周りの設定
	lf, err := os.OpenFile("./fantia_downloader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, lf)
	defer func() {
		log.Println("[E N D] FantiaDownloader")
		lf.Close()
	}()
	log.SetOutput(mw)

	log.Println("[START] FantiaDownloader")

	if *FantiaSessionID == "" {
		log.Fatalln("Key is not specified.")
	}

	// HTTPDefaultにCookieを設定
	HTTPDefaultClient(FantiaSessionID)

	// フォローしているファンクラブの一覧を取得
	fanclubs, err := fantia.GetFanclubs(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	// 出力先フォルダの設定
	fantia.SetOutput(*DownloadDir)

	// 上書きの設定
	fantia.SetOverwrite(*Overwrite)

	// 進捗の表示
	// TODO (verboseにそのうち置き換え)
	fantia.SetProgress(*Progress)

	// 年月日指定
	fantia.SetDate(*Date)

	// 待ち時間
	fantia.SetWaitTime(*WaitTime)

	fantia.SetRetry(10)
	fantia.SetRetryInterval(10)

	var fanclub int
	if *FanclubID != "" {
		fanclub, err = strconv.Atoi(*FanclubID)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// ファンクラブのID順に取得
	for _, fanclub_id := range fanclubs.FanclubIDs {
		if *FanclubID != "" && fanclub_id != fanclub {
			continue
		}

		// ファンクラブの取得
		fanclub_name, err := fantia.GetFanclub(http.DefaultClient, fanclub_id, *ExcludeFreePlan)
		if err != nil {
			log.Fatalln(err)
		}

		// ファンクラブの投稿IDを全て取得
		post_ids := fantia.GetFanclubPage(http.DefaultClient, fanclub_id)

		for _, post_id := range post_ids {
			if brk, err := fantia.GetPost(http.DefaultClient, fanclub_name, post_id); brk {
				break
			} else if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func HTTPDefaultClient(session *string) error {
	jar, err := cookiejar.New(new(cookiejar.Options))
	if err != nil {
		return err
	}

	http.DefaultClient.Jar = jar

	fantia, err := url.Parse(fantia.FANTIA_BASE_URI)
	if err != nil {
		return err
	}

	cookies := []*http.Cookie{{
		Name:  `_session_id`,
		Value: *session,
	}}

	jar.SetCookies(fantia, cookies)
	return nil
}
