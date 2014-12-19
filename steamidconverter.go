//Package steamidconverter helps convert various Steam ID formats
package steamidconverter

import (
	"fmt"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"
	"errors"
	"io/ioutil"
	// "io"
)

type Steam struct {
	apikey string
	client *http.Client
	identifier uint64
}

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
		return 0, err
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return 0, err
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
		return "", err
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return "", err
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
	type ResponseData struct {
		Steamid string `json:"steamid"`
		Success int `json:"success"`
	}
	type Response struct {
		Response ResponseData `json:"response"`
	}
	var dat Response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		return 0, err
	}
	response := dat
	success := response.Response.Success
	if success != 1 {
		return 0, errors.New("JSON not successful")
	}
	sid64 := response.Response.Steamid
	w, err := strconv.ParseUint(sid64, 10, 0)
	if err != nil {
		return 0, err
	}
	return w, nil
}

func (s *Steam) resolveVUrlHelper(url string) string {
	urlsplit := strings.Split(url, "/")
	user := urlsplit[len(urlsplit)-2]
	return fmt.Sprintf("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%v&vanityurl=%v", s.apikey, user)
}