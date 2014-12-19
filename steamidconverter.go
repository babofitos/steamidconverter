//Package steamidconverter helps convert various Steam ID formats
package steamidconverter

import (
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"
)

type Steam struct {
	apikey string
	client *http.Client
	identifier uint64
}

type ResponseData struct {
	Steamid string `json:"steamid"`
	Success int `json:"success"`
}

type Response struct {
	Response ResponseData `json:"response"`
}

//New returns a new Steam instance storing the api key
func New(apikey string) *Steam {
	client := new(http.Client)
	//Steam identifier
	var identifier uint64 = 76561197960265728
	return &Steam{apikey, client, identifier}
}

//ConvertToText converts a Steam64 ID to STEAM_0:Y:Z format
func (s *Steam) ConvertToText(w uint64) string {
	y := w % 2
	v := s.identifier
	z := (w - y - v) / 2
	return fmt.Sprintf("STEAM_0:%d:%d", y, z)
}

//ConvertTo64 converts STEAM_X:Y:Z to a Steam64 format
func (s *Steam) ConvertTo64(text string) (uint64, error) {
	sidSlice := strings.Split(text, ":")
	z, err := strconv.ParseUint(sidSlice[2], 10, 0)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint(%q) => _, %q", sidSlice[2], err)
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint(%q) => _, %q", sidSlice[1], err)
	}
	v := s.identifier
	w := z*2 + v + y
	return w, nil
}

//ConvertToSteam3 converts STEAM_X:Y:Z to a Steam3 [U:1:W] format
func (s *Steam) ConvertToSteam3(text string) (string, error) {
	sidSlice := strings.Split(text, ":")
	z, err := strconv.ParseUint(sidSlice[2], 10, 0)
	if err != nil {
		return "", fmt.Errorf("strconv.ParseUint(%q) => _, %q", sidSlice[2], err)
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return "", fmt.Errorf("strconv.ParseUint(%q) => _, %q", sidSlice[1], err)
	}
	w := z*2 + y
	return fmt.Sprintf("[U:1:%d]", w), nil
}

//ConvertVanityTo64 converts a vanity URL into a Steam64 format
func (s *Steam) ConvertVanityTo64(vanityUrl string) (uint64, error) {
	apiReqUrl := s.resolveVUrlHelper(vanityUrl)
	resp, err := s.client.Get(apiReqUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var response Response
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&response); err != nil {
		return 0, err
	}
	if success := response.Response.Success; success != 1 {
		return 0, fmt.Errorf("Response.success != 1 in struct %#v", response)
	}
	sid64 := response.Response.Steamid
	w, err := strconv.ParseUint(sid64, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseUint(%d) => _, %q", sid64, err)
	}
	return w, nil
}

func (s *Steam) resolveVUrlHelper(url string) string {
	urlsplit := strings.Split(url, "/")
	user := urlsplit[len(urlsplit)-2]
	return fmt.Sprintf("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%v&vanityurl=%v", s.apikey, user)
}