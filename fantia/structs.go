package fantia

import (
	"net/http"
	"regexp"
	"time"
)

const (
	PlanJoined = "joined"
)

const (
	RIGHT_ARROWS = "-> "

	COOKIE_SESSION_ID = `_session_id`

	FANTIA_BASE_URI = "https://fantia.jp"

	FANTIA_API_FANCLUBS     = "https://fantia.jp/api/v1/me/fanclubs"
	FANTIA_API_FANCLUB_INFO = "https://fantia.jp/api/v1/fanclubs/%d"
	FANTIA_API_POST_INFO    = "https://fantia.jp/api/v1/posts/%d"
	FANTIA_POST_CSRF_TOKEN  = "https://fantia.jp/posts/%d"

	FANTIA_FANCLUB_POSTS = "https://fantia.jp/fanclubs/%d/posts?page=%d"

	POST_DIR_FORMAT = "2006-01-02_150405"

	FOLDER_JPEG = "folder.jpeg"

	IMAGE_BITMAP = "image/bmp"
	IMAGE_JPEG   = "image/jpeg"
	IMAGE_GIF    = "image/gif"
	IMAGE_PNG    = "image/png"
)

var (
	Posts       = regexp.MustCompile(`"/posts/([0-9]+)"`)
	PostItemExt = regexp.MustCompile(`^.+/.+\.(.+)\?.+$`)
	IconExt     = regexp.MustCompile(`.+\/.+\.(\w+)$`)
	ExtJpeg     = regexp.MustCompile(`(?i)^(jpg|jpeg)$`)
	ExtPng      = regexp.MustCompile(`(?i)^(png)$`)
	ExtGif      = regexp.MustCompile(`(?i)^(gif)$`)
)

type Config struct {
	SessionID *string
	Output    *string
}

type FantiaDownloader struct {
	config *Config
	client *http.Client
}

type fanclubs struct {
	Result     bool  `json:"result"`
	FanclubIDs []int `json:"fanclub_ids"`
}

type FanclubData struct {
	Fanclub struct {
		ID   int `json:"id"`
		User struct {
			ID                     int    `json:"id"`
			ToranoanaIdentifyToken string `json:"toranoana_identify_token"`
			Name                   string `json:"name"`
			Image                  struct {
				Small  string `json:"small"`
				Medium string `json:"medium"`
				Large  string `json:"large"`
			} `json:"image"`
			ProfileText interface{} `json:"profile_text"`
			HasFanclub  bool        `json:"has_fanclub"`
		} `json:"user"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
			URI  struct {
				Fanclub  string `json:"fanclub"`
				Products string `json:"products"`
				Posts    string `json:"posts"`
			} `json:"uri"`
		} `json:"category"`
		Name                       string `json:"name"`
		CreatorName                string `json:"creator_name"`
		FanclubName                string `json:"fanclub_name"`
		FanclubNameWithCreatorName string `json:"fanclub_name_with_creator_name"`
		FanclubNameOrCreatorName   string `json:"fanclub_name_or_creator_name"`
		Title                      string `json:"title"`
		Cover                      struct {
			Thumb    string `json:"thumb"`
			Medium   string `json:"medium"`
			Main     string `json:"main"`
			Ogp      string `json:"ogp"`
			Original string `json:"original"`
		} `json:"cover"`
		Icon struct {
			Thumb    string `json:"thumb"`
			Main     string `json:"main"`
			Original string `json:"original"`
		} `json:"icon"`
		IsJoin        bool `json:"is_join"`
		FanCount      int  `json:"fan_count"`
		PostsCount    int  `json:"posts_count"`
		ProductsCount int  `json:"products_count"`
		URI           struct {
			Show     string `json:"show"`
			Posts    string `json:"posts"`
			Plans    string `json:"plans"`
			Products string `json:"products"`
		} `json:"uri"`
		IsBlocked   bool   `json:"is_blocked"`
		Comment     string `json:"comment"`
		RecentPosts []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Comment string `json:"comment"`
			Rating  string `json:"rating"`
			Thumb   struct {
				Thumb    string `json:"thumb"`
				Medium   string `json:"medium"`
				Large    string `json:"large"`
				Main     string `json:"main"`
				Ogp      string `json:"ogp"`
				Micro    string `json:"micro"`
				Original string `json:"original"`
			} `json:"thumb"`
			ThumbMicro     string `json:"thumb_micro"`
			ShowAdultThumb bool   `json:"show_adult_thumb"`
			PostedAt       string `json:"posted_at"`
			LikesCount     int    `json:"likes_count"`
			Liked          bool   `json:"liked"`
			IsContributor  bool   `json:"is_contributor"`
			URI            struct {
				Show string      `json:"show"`
				Edit interface{} `json:"edit"`
			} `json:"uri"`
			IsPulishOpen    bool      `json:"is_pulish_open"`
			IsBlog          bool      `json:"is_blog"`
			ConvertedAt     time.Time `json:"converted_at"`
			FanclubBrand    int       `json:"fanclub_brand"`
			SpecialReaction struct {
				Reaction    string `json:"reaction"`
				Kind        string `json:"kind"`
				DisplayType string `json:"display_type"`
			} `json:"special_reaction"`
			RedirectURLFromSave string `json:"redirect_url_from_save"`
		} `json:"recent_posts"`
		RecentProducts []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			Category struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				URI  string `json:"uri"`
			} `json:"category"`
			Thumb struct {
				Thumb    string `json:"thumb"`
				Small    string `json:"small"`
				Main     string `json:"main"`
				Ogp      string `json:"ogp"`
				Micro    string `json:"micro"`
				Original string `json:"original"`
			} `json:"thumb"`
			ShowAdultThumb bool        `json:"show_adult_thumb"`
			Stock          interface{} `json:"stock"`
			Price          int         `json:"price"`
			Likes          struct {
				Count   int  `json:"count"`
				HasLike bool `json:"has_like"`
			} `json:"likes"`
			URI        string `json:"uri"`
			ThumbMicro string `json:"thumb_micro"`
			Rating     string `json:"rating"`
		} `json:"recent_products"`
		Plans []struct {
			ID          int         `json:"id"`
			Price       int         `json:"price"`
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Limit       int         `json:"limit"`
			Thumb       string      `json:"thumb"`
			VacantSeat  interface{} `json:"vacant_seat"`
			Order       struct {
				Status     string `json:"status"`
				IsOneclick bool   `json:"is_oneclick"`
				URI        string `json:"uri"`
			} `json:"order"`
		} `json:"plans"`
		Background    interface{} `json:"background"`
		PointTopUsers []struct {
			ID             int    `json:"id"`
			SupportComment string `json:"support_comment"`
			SupportImage   struct {
				Medium   interface{} `json:"medium"`
				Main     interface{} `json:"main"`
				Original interface{} `json:"original"`
			} `json:"support_image"`
			SupportPoint int `json:"support_point"`
			ExtraPayPlan int `json:"extra_pay_plan"`
			User         struct {
				ID                     int    `json:"id"`
				ToranoanaIdentifyToken string `json:"toranoana_identify_token"`
				Name                   string `json:"name"`
				Image                  struct {
					Small  string `json:"small"`
					Medium string `json:"medium"`
					Large  string `json:"large"`
				} `json:"image"`
				ProfileText string      `json:"profile_text"`
				HasFanclub  interface{} `json:"has_fanclub"`
				Fanclub     interface{} `json:"fanclub"`
			} `json:"user"`
		} `json:"point_top_users"`
		SupportPoint      int           `json:"support_point"`
		SupportPointGoals []interface{} `json:"support_point_goals"`
	} `json:"fanclub"`
}

type PostData struct {
	Post struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Rating  string `json:"rating"`
		Thumb   struct {
			Thumb    string `json:"thumb"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Main     string `json:"main"`
			Ogp      string `json:"ogp"`
			Micro    string `json:"micro"`
			Original string `json:"original"`
		} `json:"thumb"`
		ThumbMicro     string `json:"thumb_micro"`
		ShowAdultThumb bool   `json:"show_adult_thumb"`
		PostedAt       string `json:"posted_at"`
		LikesCount     int    `json:"likes_count"`
		Liked          bool   `json:"liked"`
		IsContributor  bool   `json:"is_contributor"`
		URI            struct {
			Show string      `json:"show"`
			Edit interface{} `json:"edit"`
		} `json:"uri"`
		IsPulishOpen        bool        `json:"is_pulish_open"`
		IsBlog              bool        `json:"is_blog"`
		ConvertedAt         time.Time   `json:"converted_at"`
		FanclubBrand        int         `json:"fanclub_brand"`
		SpecialReaction     interface{} `json:"special_reaction"`
		RedirectURLFromSave string      `json:"redirect_url_from_save"`
		Fanclub             struct {
			ID   int `json:"id"`
			User struct {
				ID                     int    `json:"id"`
				ToranoanaIdentifyToken string `json:"toranoana_identify_token"`
				Name                   string `json:"name"`
				Image                  struct {
					Small  string `json:"small"`
					Medium string `json:"medium"`
					Large  string `json:"large"`
				} `json:"image"`
				ProfileText string `json:"profile_text"`
				HasFanclub  bool   `json:"has_fanclub"`
			} `json:"user"`
			Category struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
				URI  struct {
					Fanclub  string `json:"fanclub"`
					Products string `json:"products"`
					Posts    string `json:"posts"`
				} `json:"uri"`
			} `json:"category"`
			Name                       string `json:"name"`
			CreatorName                string `json:"creator_name"`
			FanclubName                string `json:"fanclub_name"`
			FanclubNameWithCreatorName string `json:"fanclub_name_with_creator_name"`
			FanclubNameOrCreatorName   string `json:"fanclub_name_or_creator_name"`
			Title                      string `json:"title"`
			Cover                      struct {
				Thumb    string `json:"thumb"`
				Medium   string `json:"medium"`
				Main     string `json:"main"`
				Ogp      string `json:"ogp"`
				Original string `json:"original"`
			} `json:"cover"`
			Icon struct {
				Thumb    string `json:"thumb"`
				Main     string `json:"main"`
				Original string `json:"original"`
			} `json:"icon"`
			IsJoin        bool `json:"is_join"`
			FanCount      int  `json:"fan_count"`
			PostsCount    int  `json:"posts_count"`
			ProductsCount int  `json:"products_count"`
			URI           struct {
				Show     string `json:"show"`
				Posts    string `json:"posts"`
				Plans    string `json:"plans"`
				Products string `json:"products"`
			} `json:"uri"`
			IsBlocked   bool `json:"is_blocked"`
			RecentPosts []struct {
				ID      int    `json:"id"`
				Title   string `json:"title"`
				Comment string `json:"comment"`
				Rating  string `json:"rating"`
				Thumb   struct {
					Thumb    string `json:"thumb"`
					Medium   string `json:"medium"`
					Large    string `json:"large"`
					Main     string `json:"main"`
					Ogp      string `json:"ogp"`
					Micro    string `json:"micro"`
					Original string `json:"original"`
				} `json:"thumb"`
				ThumbMicro     string `json:"thumb_micro"`
				ShowAdultThumb bool   `json:"show_adult_thumb"`
				PostedAt       string `json:"posted_at"`
				LikesCount     int    `json:"likes_count"`
				Liked          bool   `json:"liked"`
				IsContributor  bool   `json:"is_contributor"`
				URI            struct {
					Show string      `json:"show"`
					Edit interface{} `json:"edit"`
				} `json:"uri"`
				IsPulishOpen        bool        `json:"is_pulish_open"`
				IsBlog              bool        `json:"is_blog"`
				ConvertedAt         time.Time   `json:"converted_at"`
				FanclubBrand        int         `json:"fanclub_brand"`
				SpecialReaction     interface{} `json:"special_reaction"`
				RedirectURLFromSave string      `json:"redirect_url_from_save"`
			} `json:"recent_posts"`
			RecentProducts []interface{} `json:"recent_products"`
			Plans          []struct {
				ID          int         `json:"id"`
				Price       int         `json:"price"`
				Name        string      `json:"name"`
				Description string      `json:"description"`
				Limit       int         `json:"limit"`
				Thumb       string      `json:"thumb"`
				VacantSeat  interface{} `json:"vacant_seat"`
				Order       struct {
					Status     string `json:"status"`
					IsOneclick bool   `json:"is_oneclick"`
					URI        string `json:"uri"`
				} `json:"order"`
			} `json:"plans"`
		} `json:"fanclub"`
		Tags []struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"tags"`
		Status       string `json:"status"`
		PostContents []struct {
			ID             int         `json:"id"`
			Title          string      `json:"title"`
			VisibleStatus  string      `json:"visible_status"`
			PublishedState string      `json:"published_state"`
			Category       string      `json:"category"`
			Comment        interface{} `json:"comment"`
			EmbedURL       interface{} `json:"embed_url"`
			ContentType    interface{} `json:"content_type"`
			Comments       struct {
				GetURL    string `json:"get_url"`
				PostURI   string `json:"post_uri"`
				DeleteURI string `json:"delete_uri"`
			} `json:"comments"`
			CommentsReactions struct {
				PostURI   string `json:"post_uri"`
				DeleteURI string `json:"delete_uri"`
				GetURL    string `json:"get_url"`
			} `json:"comments_reactions"`
			EmbedAPIURL string `json:"embed_api_url"`
			Reactions   struct {
				GetURL    string `json:"get_url"`
				PostURI   string `json:"post_uri"`
				DeleteURI string `json:"delete_uri"`
			} `json:"reactions"`
			ReactionTypesURL  string `json:"reaction_types_url"`
			PostContentPhotos []struct {
				ID  int `json:"id"`
				URL struct {
					Thumb    string `json:"thumb"`
					Medium   string `json:"medium"`
					Large    string `json:"large"`
					Main     string `json:"main"`
					Micro    string `json:"micro"`
					Original string `json:"original"`
				} `json:"url"`
				Comment         interface{} `json:"comment"`
				ShowOriginalURI string      `json:"show_original_uri"`
				IsConverted     bool        `json:"is_converted"`
			} `json:"post_content_photos,omitempty"`
			PostContentPhotosMicro []string `json:"post_content_photos_micro"`
			Plan                   struct {
				ID          int    `json:"id"`
				Price       int    `json:"price"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Limit       int    `json:"limit"`
				Thumb       string `json:"thumb"`
			} `json:"plan"`
			Product          interface{} `json:"product"`
			OnsaleBacknumber interface{} `json:"onsale_backnumber"`
			BacknumberLink   interface{} `json:"backnumber_link"`
			JoinStatus       interface{} `json:"join_status"`
			ParentPost       struct {
				Title    string    `json:"title"`
				URL      string    `json:"url"`
				Date     time.Time `json:"date"`
				Deadline time.Time `json:"deadline"`
			} `json:"parent_post"`
			IsConverted bool   `json:"is_converted,omitempty"`
			Filename    string `json:"filename,omitempty"`
			DownloadURI string `json:"download_uri,omitempty"`
		} `json:"post_contents"`
		Deadline          string `json:"deadline"`
		PublishReservedAt string `json:"publish_reserved_at"`
		Comments          struct {
			PostURI   string `json:"post_uri"`
			DeleteURI string `json:"delete_uri"`
			GetURL    string `json:"get_url"`
		} `json:"comments"`
		BlogComment       string `json:"blog_comment"`
		CommentsReactions struct {
			PostURI   string `json:"post_uri"`
			DeleteURI string `json:"delete_uri"`
			GetURL    string `json:"get_url"`
		} `json:"comments_reactions"`
		Reactions struct {
			PostURI   string `json:"post_uri"`
			DeleteURI string `json:"delete_uri"`
			GetURL    string `json:"get_url"`
		} `json:"reactions"`
		ReactionTypesURL string `json:"reaction_types_url"`
		OgpAPIURL        string `json:"ogp_api_url"`
		Links            struct {
			Previous struct {
				ID      int    `json:"id"`
				Title   string `json:"title"`
				Comment string `json:"comment"`
				Rating  string `json:"rating"`
				Thumb   struct {
					Thumb    string `json:"thumb"`
					Medium   string `json:"medium"`
					Large    string `json:"large"`
					Main     string `json:"main"`
					Ogp      string `json:"ogp"`
					Micro    string `json:"micro"`
					Original string `json:"original"`
				} `json:"thumb"`
				ThumbMicro     string `json:"thumb_micro"`
				ShowAdultThumb bool   `json:"show_adult_thumb"`
				PostedAt       string `json:"posted_at"`
				LikesCount     int    `json:"likes_count"`
				Liked          bool   `json:"liked"`
				IsContributor  bool   `json:"is_contributor"`
				URI            struct {
					Show string      `json:"show"`
					Edit interface{} `json:"edit"`
				} `json:"uri"`
				IsPulishOpen        bool        `json:"is_pulish_open"`
				IsBlog              bool        `json:"is_blog"`
				ConvertedAt         time.Time   `json:"converted_at"`
				FanclubBrand        int         `json:"fanclub_brand"`
				SpecialReaction     interface{} `json:"special_reaction"`
				RedirectURLFromSave string      `json:"redirect_url_from_save"`
			} `json:"previous"`
			Next interface{} `json:"next"`
		} `json:"links"`
		IsFanclubTipAccept bool `json:"is_fanclub_tip_accept"`
		IsFanclubJoined    bool `json:"is_fanclub_joined"`
	} `json:"post"`
}
