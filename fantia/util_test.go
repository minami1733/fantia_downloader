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
