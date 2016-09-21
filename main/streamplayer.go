package main

import (
	"os"
	"strings"
	"streamplayer/stream"
	"log"
	"encoding/json"
	"flag"
)

const ENVIRONMENT_VARIABLE = "STREAMPLAYER"

type Configuration struct {
	Player		string		`json:"player"`
	OAuth		string		`json:"oAuthToken"`
}

type DefaultConfiguration struct {
	Player	string
}

func main() {
	authenticate := flag.Bool("authenticate", false, "Used to generate an OAuth token connected to Streamplayer")
	flag.Parse()

	if *authenticate {
		stream.GenerateAuthToken()
		return
	}

	channel := os.Args[1]
	quality := os.Args[2]

	if strings.EqualFold(quality, "best") {
		quality = "Source"
	}

	loadConfiguration()

	if (stream.OAuth != "") {
		stream.DoAuthenticate()
	}

	//getStreams()
	stream.PlayStream(channel, quality)
	log.Println("Terminating...")
}

func loadConfiguration() {
	path := os.Getenv(ENVIRONMENT_VARIABLE)
	stream.Playlist = path
	file, err := os.Open(path + "\\conf.json")
	if err != nil {
		log.Println("Error opening configuration file, creating a new one with default values")
		createDefault(path)
	}

	if stream.Player != "" {
		return
	}

	decoder := json.NewDecoder(file)
	configuration := new(Configuration)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("Error decoding configuration file, creating a new one with default values")
		createDefault(path)
	} else {
		stream.Player = configuration.Player
		stream.OAuth = configuration.OAuth
	}
}

func createDefault(path string) {
	defaultPath := "C:\\Program Files\\VideoLAN\\VLC\\vlc.exe"
	defaultOAuth := ""

	stream.Player = defaultPath
	stream.OAuth = defaultOAuth

	defaultConfig := Configuration{defaultPath, defaultOAuth}
	file, _ := os.Create(path + "\\conf.json")
	encoder := json.NewEncoder(file)
	encoder.Encode(&defaultConfig)
}