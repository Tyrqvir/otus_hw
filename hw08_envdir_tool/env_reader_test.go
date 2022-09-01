package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const tmpFolder = "./tmpFolder"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makeFile(fileName string, body string) {
	file, err := os.Create(filepath.Join(tmpFolder, fileName))
	check(err)

	defer file.Close()

	_, err = file.WriteString(body)
	check(err)
}

func _before() {
	err := os.Mkdir(tmpFolder, os.ModePerm)
	check(err)

	makeFile("HELLO", "world")
	makeFile("WITH_SPACES", " with spaces    ")
	makeFile("EMPTY", "")
	makeFile("WITH_TERMINAL_NULLS", "val\x00")
	makeFile("WITH_EQUAL=", "with equal")
}

func _after() {
	err := os.RemoveAll(tmpFolder)
	check(err)
}

func TestReadDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name: "Should read files config from folder correctly",
			args: args{
				dir: tmpFolder,
			},
			want: map[string]EnvValue{
				"HELLO": {
					Value:      "world",
					NeedRemove: false,
				},
				"WITH_SPACES": {
					Value:      " with spaces",
					NeedRemove: false,
				},
				"EMPTY": {
					Value:      "",
					NeedRemove: true,
				},
				"WITH_TERMINAL_NULLS": {
					Value:      "val\n",
					NeedRemove: false,
				},
				"WITH_EQUAL": {
					Value:      "with equal",
					NeedRemove: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_before()
			defer _after()
			got, err := ReadDir(tt.args.dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
