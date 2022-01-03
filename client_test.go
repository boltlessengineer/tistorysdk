package tistorysdk_test

import (
	"testing"

	"github.com/boltlessengineer/tistorysdk"
)

func TestGetAuthCodeURL(t *testing.T) {
	want := `https://www.tistory.com/oauth/authorize?client_id=%7B%7BAppID%7D%7D&redirect_uri=%7B%7BRedirectUri%7D%7D&response_type=code&state=%7B%7Bstate%7D%7D`
	tistory := tistorysdk.NewClient("{{AppID}}", "{{SecretKey}}", "{{RedirectUri}}")
	got := tistory.GetAuthCodeURL("{{state}}")
	if got.String() != want {
		t.Errorf("GetAuthCodeURL() got = %v, want = %v", got, want)
	}
}
