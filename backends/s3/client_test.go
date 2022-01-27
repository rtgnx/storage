package s3

import (
	"reflect"
	"strings"
	"testing"
)

var (
	config = ConnectionProfile{
		URL:       "play.min.io",
		AccessKey: "Q3AM3UQ867SPQQA43P2F",
		SecretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
	}
	bucket = "go-storage"
)

func TestS3_WriteBytes(t *testing.T) {
	type args struct {
		prefix string
		b      []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test successful write",
			args:    args{b: []byte("some data"), prefix: "test.data"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s3, err := FromCfg(config, bucket)
			s3.MakeBucket(bucket)

			if err != nil {
				t.Fatal(err)
			}

			if err := s3.WriteBytes(tt.args.prefix, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("S3.WriteBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			hasKey := false
			for _, name := range s3.List("") {
				if strings.Compare(name, tt.args.prefix) == 0 {
					hasKey = true
					break
				}
			}

			if !hasKey {
				t.Fail()
			}
		})
	}
}

func TestS3_ReadBytes(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "successful read",
			args:    args{prefix: "test.data"},
			want:    []byte("some data"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s3, err := FromCfg(config, bucket)
		if err != nil {
			t.Fatal(err)
		}
		// Create Object
		TestS3_WriteBytes(t)

		t.Run(tt.name, func(t *testing.T) {

			got, err := s3.ReadBytes(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3.ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S3.ReadBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
