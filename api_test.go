package main

import (
	"net/http"
	"testing"
)

func TestClient_GetPage(t *testing.T) {
	type fields struct {
		BaseURL    string
		HTTPClient *http.Client
	}
	type args struct {
		pageId string
	}
	type page struct {
		Object string
		Id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    page
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				BaseURL: "https://api.notion.com/v1",
				HTTPClient: new(http.Client)},
			args: args{
				pageId: "832b5f3f690e4b9ca7b8c916e89ff30e",
			},
			want: page{
				Object: "page", 
				Id: "832b5f3f-690e-4b9c-a7b8-c916e89ff30e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				BaseURL:    tt.fields.BaseURL,
				HTTPClient: tt.fields.HTTPClient,
			}
			got, err := c.GetPage(tt.args.pageId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.GetPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Object != tt.want.Object {
				t.Fatalf("Client.GetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
