package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	// https://github.com/moby/moby/tree/master/client
	// https://godoc.org/github.com/moby/moby/client doc
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

func (ii *imageinfo) print() {
	fmt.Printf("%s\t%s\t%d\t%s\t%s\t%d\t%s\n", ii.repository, ii.tag, ii.rawSize, ii.humanSize, ii.dateOfCreation, ii.rawDuration, ii.humanDuration)
}

var (
	addr               = flag.String("192.168.1.129", ":8082", "The address to listen on for HTTP requests.")
	timeBetweenMetrics = flag.Duration("timeinsec.between", 10*time.Second, "Time between two logs")
)

// https://prometheus.io/docs/concepts/metric_types/
var (
	imageDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "codingame",
			Subsystem: "my_computer",
			Name:      "docker_image_duration",
			Help:      "Docker image existence duration (since ??)",
		},
		[]string{"repository", "tag"},
	)
	imageSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "codingame",
			Subsystem: "my_computer",
			Name:      "docker_image_size",
			Help:      "Docker image size (in Mo)",
		},
		[]string{"repository", "tag"},
	)
)

func listEvents(cli *client.Client) []imageinfo {
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	var description []imageinfo
	description = make([]imageinfo, len(images))
	for i, image := range images {
		response, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
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

func main() {
	flag.Parse()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	prometheus.MustRegister(imageDuration)
	prometheus.MustRegister(imageSize)
	go func() {
		for {
			imageDuration.Reset()
			imageSize.Reset()
			infos := listEvents(cli)
			for _, in := range infos {
				imageDuration.WithLabelValues(in.repository, in.tag).Set(float64(in.rawDuration))
				imageSize.WithLabelValues(in.repository, in.tag).Set(float64(in.rawSize / 1000000))
			}
			fmt.Printf("ok\n")
			//rpcDurations.WithLabelValues("normal").Observe(v)
			//rpcDurationsHistogram.Observe(v)
			time.Sleep(*timeBetweenMetrics)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
