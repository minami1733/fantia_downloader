package fantia

import (
	"reflect"
	"testing"
)

func TestForbiddenTextRename(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ForbiddenTextRename(tt.args.name); got != tt.want {
				t.Errorf("ForbiddenTextRename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFolderExist(t *testing.T) {
	type args struct {
		dir  string
		make bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isFolderExist(tt.args.dir, tt.args.make)
			if (err != nil) != tt.wantErr {
				t.Errorf("isFolderExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isFolderExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDuplicateInt(t *testing.T) {
	type args struct {
		source []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicateInt(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicateInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CutStringToLimit(t *testing.T) {
	type args struct {
		input string
		limit uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"long name",
			args{
				"downloads/Dikk0Fantia (ディッコ)/2018-10-29_085551 - ツイ垢の凍結解除されました＋予告やラフ置き場/001 - アルトリア描く予定でしたが、今日夢にアビィちゃんが出てきたので、 今週はふぁ〇きゅーとか汚い言葉をスラムで言いまくってたら路地裏に連れ込まれて FXXKされちゃうアビィちゃん描きます(＊'ω'＊)",
				245,
			},
			"downloads/Dikk0Fantia (ディッコ)/2018-10-29_085551 - ツイ垢の凍結解除されました＋予告やラフ置き場/001 - アルトリア描く予定でしたが、今日夢にアビィちゃんが出てきたので、 今週はふぁ〇",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CutStringToLimit(tt.args.input, tt.args.limit); got != tt.want {
				t.Errorf("CutStringToLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
