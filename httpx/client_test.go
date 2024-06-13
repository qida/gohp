package httpx

import "testing"

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
				url:        "http://192.168.114.60:24722/v1/ai/window/upload_image",
				name_param: "filename",
				file_path:  "D:\\project\\LabelCenter\\test\\data\\image\\1.png",
				params: map[string]string{
					"device_id": "123",
				},
				header: map[string]string{
					"Content-Type": "multipart/form-data",
				},
			},
			wantErr: true,
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
