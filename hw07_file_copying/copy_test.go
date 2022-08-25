package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}

	input := computedPathByName("input.txt")
	tmpOutput, _ := os.CreateTemp("", "tests")

	tests := []struct {
		name    string
		args    args
		output  string
		wantErr bool
	}{
		{
			name: "Test with offset 0 and limit 0",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   0,
				limit:    0,
			},
			output:  computedPathByName("out_offset0_limit0.txt"),
			wantErr: false,
		},
		{
			name: "Test with offset 0 and limit 10",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   0,
				limit:    10,
			},
			output:  computedPathByName("out_offset0_limit10.txt"),
			wantErr: false,
		},
		{
			name: "Test with offset 0 and limit 1000",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   0,
				limit:    1000,
			},
			output:  computedPathByName("out_offset0_limit1000.txt"),
			wantErr: false,
		},
		{
			name: "Test with offset 0 and limit 10000",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   0,
				limit:    10000,
			},
			output:  computedPathByName("out_offset0_limit10000.txt"),
			wantErr: false,
		},
		{
			name: "Test with offset 100 and limit 1000",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   100,
				limit:    1000,
			},
			output:  computedPathByName("out_offset100_limit1000.txt"),
			wantErr: false,
		},
		{
			name: "Test with offset 6000 and limit 1000",
			args: args{
				fromPath: input,
				toPath:   tmpOutput.Name(),
				offset:   6000,
				limit:    1000,
			},
			output:  computedPathByName("out_offset6000_limit1000.txt"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
			sourceFileInfo, _ := os.Stat(tt.output)
			targetFileInfo, _ := os.Stat(tt.args.toPath)
			require.Equal(t, sourceFileInfo.Size(), targetFileInfo.Size())
		})
	}
}

func computedPathByName(name string) string {
	return fmt.Sprintf("%s/%s", "testData", name)
}
