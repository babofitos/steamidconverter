package steamidconverter

import (
	"testing"
	"net/http"
	"net/url"
	"net/http/httptest"
	"encoding/json"
	"strconv"
)

var textTests = []struct {
	in uint64
	out string
}{
	{76561197960430077, "STEAM_0:1:82174"},
	{76561198013163430, "STEAM_0:0:26448851"},
	{76561197960271553, "STEAM_0:1:2912"},
	{76561197964812891, "STEAM_0:1:2273581"},
}

var sixtyFourTests = []struct {
	in string
	out uint64
}{
	{"STEAM_0:1:82174", 76561197960430077},
	{"STEAM_0:0:26448851", 76561198013163430},
	{"STEAM_0:1:2912", 76561197960271553},
	{"STEAM_0:1:2273581", 76561197964812891},
}

var steam3Tests = []struct {
	in string
	out string
}{
	{"STEAM_0:1:82174", "[U:1:164349]"},
	{"STEAM_0:0:26448851", "[U:1:52897702]"},
	{"STEAM_0:1:2912", "[U:1:5825]"},
	{"STEAM_0:1:2273581", "[U:1:4547163]"},
}

var vanityTests = []struct {
	in string
	out uint64
}{
	{"http://steamcommunity.com/id/panvertigo/", 76561198000670105},
	{"http://steamcommunity.com/id/FireSlash/", 76561197972495328},
}

var apikey string = ""
var identifier uint64 = 76561197960265728

func TestConvertToText(t *testing.T) {
	s := New(apikey)
	for _, tt := range textTests {
		sidText := s.ConvertToText(tt.in)
		if sidText != tt.out {
			t.Errorf("ConvertToText(%d) => %q, want %q", tt.in, sidText, tt.out)
		}
	}
}

func TestConvertTo64(t *testing.T) {
	s := New(apikey)
	for _, tt := range sixtyFourTests {
		sid64, err := s.ConvertTo64(tt.in)
		if err != nil {
			t.Errorf("ConvertTo64(%q) => _, %q", tt.in, err)
		}
		if sid64 != tt.out {
			t.Errorf("ConvertTo64(%q) => %d, want %d", tt.in, sid64, tt.out)
		}
	}
}

func TestConvertToSteam3(t *testing.T) {
	s := New(apikey)
	for _, tt := range steam3Tests {
		sid3, err := s.ConvertToSteam3(tt.in)
		if err != nil {
			t.Errorf("ConvertToSteam3(%q) => _, %q", tt.in, err)
		}
		if sid3 != tt.out {
			t.Errorf("ConvertToSteam3(%q) => %q, want %q", tt.in, sid3, tt.out)
		}
	}
}

func TestConvertVanityTo64(t *testing.T) {
	for _, tt := range vanityTests {
		type ResponseData struct {
			Steamid string `json:"steamid"`
			Success int `json:"success"`
		}
		type Response struct {
			Response ResponseData `json:"response"`
		}
		var dat Response
		dat.Response.Steamid = strconv.FormatUint(tt.out, 10)
		dat.Response.Success = 1
		b, err := json.Marshal(dat)
		if err != nil {
			t.Errorf("ConvertVanity(%q) => _, %q", tt.in, err)
		}
		//test server to respond json
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
    	w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		}
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		proxy, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("ConvertVanity(%q) => _, %q", tt.in, err)
		}
		//set client transport to proxy with test server url
		transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
		client := &http.Client{Transport: transport}
		//literal instead of constructor to pass in our custom client
		s := &Steam{apikey, client, identifier}
		w, err := s.ConvertVanityTo64(tt.in)
		if err != nil {
			t.Errorf("ConvertVanityTo64(%q) => _, %q", tt.in, err)
		}
		if w != tt.out {
			t.Errorf("ConvertVanityTo64(%q) => %d, want %d", tt.in, w, tt.out)
		}
	}
}