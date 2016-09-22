package parsers

import (
	"bytes"
	"bufio"
	"strings"
	"regexp"
	"strconv"
)

const HEADER = "#EXTM3U"
const INFO = "#EXT-X-TWITCH-INFO"

type Playlist struct {
	Info	Info
	Streams []Stream

}

type Info struct {
	Node			string
	ManifestNodeType	string
	ManifestNode		string
	Suppress		bool
	UserIP			string
	ServingID		string
	Cluster			string
	ABS			bool
	BroadcastID		string
	StreamTime		float64
	ManifestCluster		string
}

type Stream struct {
	Type		string
	GroupID 	string
	Name		string
	Autoselect	bool
	Default		bool
	ProgramID	int
	Bandwidth	int
	Resolution	string
	Codecs		string
	Video		string
	Link		string
}

// Parse takes a slice of bytes, expecting them to be a m3u8 twitch channel playlist. It parses this content and returns
// a structure describing all the info extracted from this.
func Parse(content []byte) (Playlist){

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Scan()

	firstLine := scanner.Text()
	if !strings.HasPrefix(firstLine, HEADER) {
		panic("This is not a valid m3u8 playlist")
	}

	info := parseInfo(scanner)

	i := 0
	var streams [10]Stream

	for scanner.Scan() {
		streams[i] = parseStream(scanner)
		i++
	}

	return Playlist{info, streams[:i]}
}

func parseInfo(scanner *bufio.Scanner) (Info) {
	scanner.Scan()
	firstLine := scanner.Text()
	if !strings.HasPrefix(firstLine, INFO) {
		panic("Line \"" + firstLine + "\" is not a valid info")
	}

	firstLine = strings.Replace(firstLine, "\"", "", -1)
	matcher := regexp.MustCompile("[\\w\\-]*[=][\\w\\d\\-.]*");
	matches := matcher.FindAllString(firstLine, -1)

	node := extractValue(matches[0])
	manifestNodeType := extractValue(matches[1])
	manifestNode := extractValue(matches[2])
	suppress, _ := strconv.ParseBool(extractValue(matches[3]))
	ip := extractValue(matches[5])
	id := extractValue(matches[6])
	cluster := extractValue(matches[7])
	abs, _ := strconv.ParseBool(extractValue(matches[8]))
	broadcast := extractValue(matches[9])
	time, _ := strconv.ParseFloat(extractValue(matches[10]), 64)
	manifestCluster := extractValue(matches[11])

	return Info{node, manifestNodeType, manifestNode, suppress, ip, id, cluster, abs, broadcast, time, manifestCluster}
}

func parseStream(scanner *bufio.Scanner) (Stream) {
	line := strings.Replace(scanner.Text(), "\"", "", -1)
	matcher := regexp.MustCompile("[\\w\\-]*[=][\\w\\d\\-.]*");
	matches := matcher.FindAllString(line, -1)
	scanner.Scan()
	line = strings.Replace(scanner.Text(), "\"", "", -1)
	matches = append(matches, matcher.FindAllString(line, -1)...)

	scanner.Scan()
	url := scanner.Text()

	linkType := extractValue(matches[0])
	groupID := extractValue(matches[1])
	name := extractValue(matches[2])
	isAudio := name == "Audio"
	autoselect, _ := strconv.ParseBool(extractValue(matches[3]))
	isDefault, _ := strconv.ParseBool(extractValue(matches[4]))
	programID, _ := strconv.Atoi(extractValue(matches[5]))
	bandwidth, _ := strconv.Atoi(extractValue(matches[6]))
	var resolution, codecs, video string
	if isAudio {
		codecs = extractValue(matches[7])
		video = extractValue(matches[8])
	} else {
		resolution = extractValue(matches[7])
		codecs = extractValue(matches[8])
		video = extractValue(matches[9])
	}

	stream := Stream{linkType, groupID, name, autoselect, isDefault, programID, bandwidth, resolution, codecs, video, url}
	return stream
}

func extractValue(property string) (string) {
	return property[strings.Index(property, "=")+1:]
}
