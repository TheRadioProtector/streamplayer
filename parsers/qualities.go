package parsers

import (
	"io"
	"bufio"
	"strings"
)

type Qualities struct {
	Info 		Info
	Quality[] 	Quality
}
type Quality struct {
	Name 		string
	Resolution	string
	Video		string
	Playlist	string
}

type Info struct {
	Time 		string
	IP		string
}

// Parse takes a body of a GET response and assumes it contains a M3U extended playlist from a twitch channel
func Parse (body io.ReadCloser) (Qualities) {

	var scanner = bufio.NewScanner(body)

	scanner.Scan()
	firstLine := scanner.Text()
	if !strings.HasPrefix(firstLine, "#EXTM3U") {
		panic("This is not a valid M3U playlist, got firstline" + firstLine)
	}

	var info Info
	qualities := make([]Quality, 10, 10)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		var name, res, video, playlist string

		if strings.HasPrefix(line, "#EXT-X-TWITCH-INFO") {
			var time, ip string

			parts := strings.Split(line, ",")
			for _, v := range parts {
				if strings.Contains(v, "SERVER-TIME") {
					time = v[strings.Index(v,"\"")+1:strings.LastIndex(v, "\"")]
				}

				if strings.Contains(v, "USER-IP") {
					ip = v[strings.Index(v,"\"")+1:strings.LastIndex(v, "\"")]
				}
			}

			info = Info {time, ip}
		}

		if strings.HasPrefix(line, "#EXT-X-MEDIA") {
			parts := strings.Split(line, ",")
			for _, v := range parts {
				if strings.Contains(v, "NAME") {
					name = v[strings.Index(v,"\"")+1:strings.LastIndex(v, "\"")]
				}
			}

			scanner.Scan()
			line = scanner.Text()

			parts = strings.Split(line, ",")

			for _, v := range parts {
				if strings.Contains(v, "RESOLUTION") {
					res = v[strings.Index(v,"=")+1:]
				}

				if strings.Contains(v, "VIDEO") {
					video = v[strings.Index(v,"\"")+1:strings.LastIndex(v, "\"")]
				}
			}

			scanner.Scan()
			line = scanner.Text()

			playlist = line
		}
		qualities[i] = Quality{name, res, video, playlist}
		i++
	}

	final := Qualities{
		info,
		qualities[:i+1]}


	return final
}