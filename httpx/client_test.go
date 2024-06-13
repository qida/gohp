package httpx

import (
	"os"
	"testing"
)

func TestClientHttp_UploadFile(t *testing.T) {
	type args struct {
		url        string
		name_param string
		file_path  string
		params     map[string]string
		header     map[string]string
	}
	tests := []struct {
		name    string
		tr      *ClientHttp
		args    args
		wantErr bool
	}{

		{
			name: "test",
			tr:   NewClientHttp(),
			args: args{
				url:        "http://192.168.114.60:22722/v1/ai/window/upload_image",
				name_param: "filename",
				file_path:  "D:\\project\\LabelCenter\\test\\data\\image\\7.jpg",
				params: map[string]string{
					"device_id": "123",
				},
				header: map[string]string{
					"Content-Type": "multipart/form-data",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UploadFile(tt.args.url, tt.args.name_param, tt.args.file_path, tt.args.params, tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("ClientHttp.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientHttp_UploadFileFromBytes(t *testing.T) {
	data, err := os.ReadFile("D:\\project\\LabelCenter\\test\\data\\image\\7.jpg")
	if err != nil {
		return
	}
	type args struct {
		url        string
		name_param string
		file_name  string
		file_data  []byte
		params     map[string]string
		header     map[string]string
	}
	tests := []struct {
		name    string
		tr      *ClientHttp
		args    args
		wantErr bool
	}{
		{
			name: "test",
			tr:   NewClientHttp(),
			args: args{
				url:        "http://192.168.114.60:22722/v1/ai/window/upload_image",
				name_param: "filename",
				file_name:  "test.jpg",
				file_data:  data,
				params: map[string]string{
					"device_id": "456",
				},
				header: map[string]string{
					"Content-Type": "multipart/form-data",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.UploadFileFromBytes(tt.args.url, tt.args.name_param, tt.args.file_name, tt.args.file_data, tt.args.params, tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("ClientHttp.UploadFileFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
