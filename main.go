package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
	"ysdowner/implements"
	"ysdowner/interfaceList"
	"ysdowner/resovle"
)

var (
	routineNum   *int
	vcpu         *int
	url          *string
	path         *string
	saveFileName *string

	resovler resovle.Resovler
)

func parseCommand() {
	url = flag.String("url", "", "URL")
	path = flag.String("path", "", "Path")
	routineNum = flag.Int("routine", 10, "Routine")
	saveFileName = flag.String("name", "video.ts", "videoName")
	vcpu = flag.Int("vcpu", 2, "VCPU")
	flag.Parse()
}

func openDesFileAndDown() {
	file, err := os.OpenFile(*saveFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var reader interfaceList.Reader
	if *url != "" {
		fromUrl := implements.ReadFromUrl{}
		err = fromUrl.SetUrl(*url)
		if err != nil {
			panic(err)
		}
		defer fromUrl.Close()
		reader = &fromUrl
	} else if *path != "" {
		fromUrl := implements.ReadFromFile{}
		err = fromUrl.SetFileName(*path)
		if err != nil {
			panic(err)
		}
		defer fromUrl.Close()
		reader = &fromUrl
	} else {
		panic("path/url all empty")
	}
	resovler.ReadFromReader(reader)
	resovler.WriteVideo(file, *routineNum)
	resovler.FormatVideoTime("[视频长度: %v 秒]\n", os.Stdout)
}

func main() {
	parseCommand()
	runtime.GOMAXPROCS(*vcpu)
	now := time.Now()
	openDesFileAndDown()
	fmt.Printf("[本次下载耗时: %v ]\n", time.Now().Sub(now))
}
