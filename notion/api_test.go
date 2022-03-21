package notion

import (
	"testing"
)

func TestClient_GetPage(t *testing.T) {
	type args struct {
		pageId string
	}
	type page struct {
		Object string
		Id     string
	}
	tests := []struct {
		name    string
		args    args
		want    page
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				pageId: "832b5f3f690e4b9ca7b8c916e89ff30e",
			},
			want: page{
				Object: "page",
				Id:     "832b5f3f-690e-4b9c-a7b8-c916e89ff30e",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
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

func TestClient_PostItem(t *testing.T) {
	type args struct {
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				item: &Item{
					Title:  "TEST-NOTION-API",
					DoDate: "2022-04-09",
					URL:    "http://example.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			if err := c.PostItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Client.PostItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}