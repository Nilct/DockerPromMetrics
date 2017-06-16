package main

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	humanize "github.com/dustin/go-humanize"
)

// https://github.com/moby/moby/tree/master/client
// https://godoc.org/github.com/moby/moby/client doc
//"github.com/moby/moby/api/types"
//"github.com/moby/moby/client"

// Cli Environment client for Docker
var Cli *client.Client

type imageinfo struct {
	repository     string
	tag            string
	rawSize        int64
	humanSize      string
	dateOfCreation string
	rawDuration    int64
	humanDuration  string
}
type byRepoTag []imageinfo

func (b byRepoTag) Len() int {
	return len(b)
}

func (b byRepoTag) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byRepoTag) Less(i, j int) bool {
	return b[i].repository+b[i].tag < b[j].repository+b[j].tag
}

func listEvents() []imageinfo {
	images, _ := Cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	var description []imageinfo
	description = make([]imageinfo, len(images))
	for i, image := range images {
		response, _, err := Cli.ImageInspectWithRaw(context.Background(), image.ID)
		if err != nil {
			panic(err)
		}
		s := strings.Split(image.RepoTags[0], ":")
		repo, tag := s[0], s[1]
		size := image.Size
		date := response.Created
		duration := image.Created
		description[i] = imageinfo{repo, tag, size, humanize.Bytes(uint64(image.Size)), date, duration, humanize.Time(time.Unix(image.Created, 0))}
	}
	sort.Sort(byRepoTag(description))
	return description
}

func initEventScraping() {
	Cli, err := client.NewEnvClient() // client pour les metriques docker
	if err != nil {
		panic(err)
	}
}
