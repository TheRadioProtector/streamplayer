package stream

import (
	"fmt"
	"streamplayer/types"
	"strings"
	"os/exec"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"math/rand"
	"log"
	"runtime"
	"streamplayer/parsers"
)

const CLIENT_ID = "rd5obf4dvmky8pvu46i18wf32kqhacr"
const API = "http://api.twitch.tv/api"
const KRAKEN = "https://api.twitch.tv/kraken"
const USHER = "http://usher.twitch.tv/api"
const PLAYLIST = USHER + "/channel/hls/{channel}.m3u8?player=twitchweb&token={token}&sig={sig}&allow_audio_only=true&allow_source=true&type=any&p={random}"
const AUTHENTICATE = "https://api.twitch.tv/kraken/oauth2/authorize?response_type=token&client_id={clientId}&redirect_uri={redirect}&scope="
const REDIRECT = "http://jakubhyncica.cz/streamplayer/landing.html"

// Player path
var Player 	string
// OAuth is a token for authentication and sign-in over twitch API
var OAuth	string

// PlayStream will check if a given channel is online and if it is, create a playlist file on desktop and play it using VLC
func PlayStream(channel string, quality string) {
	if OAuth != "" {
		log.Println()
	}

	log.Println("Opening: " + channel)
	token := getToken(channel)
	playlist := getPlaylist(channel, token)
	if len(playlist.Streams) == 0 {
		return
	}

	var selectedQuality parsers.Stream
	for _, v := range playlist.Streams {
		if strings.EqualFold(quality, v.Name) {
			selectedQuality = v
		}
	}
	log.Println("Selected: " + selectedQuality.Name + "[" + selectedQuality.Resolution + "]")
	doPlay(selectedQuality)
}

// GenerateAuthToken tries to open a web browser with a twitch URL where user can login and authenticate Streamplayer and generate his unique OAuth token
func GenerateAuthToken() {
	replacer := strings.NewReplacer("{clientId}", CLIENT_ID, "{redirect}", REDIRECT)
	authUrl := replacer.Replace(AUTHENTICATE)

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", authUrl).Start()
	case "darwin":
		err = exec.Command("open", authUrl).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", authUrl).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println("Unable to open a web browser to authenticate")
	}
}

func doPlay(stream parsers.Stream) {
	c := exec.Command(Player, stream.Link)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}

func getPlaylist(channel string, token *types.AccessToken) (parsers.Playlist) {

	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(1000000)
	replacer := strings.NewReplacer("{channel}", channel, "{token}", token.Token, "{sig}", token.Sig, "{random}", strconv.Itoa(p))
	playlistUrl := replacer.Replace(PLAYLIST)
	resp, err := http.Get(playlistUrl)

	if err != nil {
		panic(err.Error())
	}

	content, _ := ioutil.ReadAll(resp.Body)
	s := string(content)
	if strings.Contains(s, "<table") {
		log.Println("Channel " + channel + " is not currently online")
		return parsers.Playlist{}
	}

	return parsers.Parse(content)
}

func getToken(channel string) (*types.AccessToken) {
	resp, err := http.Get(API + "/channels/" + channel + "/access_token" + "?client_id=" + CLIENT_ID)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var s = new(types.AccessToken)
	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err.Error())
	}

	return s
}

func getStreams() {
	resp, err := http.Get(KRAKEN + "/streams" + CLIENT_ID)
	if err != nil {
		panic(err.Error())
	}

	var s = new(types.Streams)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	for _, v := range s.Stream {
		fmt.Println(v.Channel.DisplayName + " playing " + v.Channel.Game)
	}
}