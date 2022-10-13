package resovle

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unsafe"
	"ysdowner/interfaceList"
)

type Resovler struct {
	duration  float64
	VideoList []string
	dreader   interfaceList.Reader
}

func (rec *Resovler) bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func (rec *Resovler) ReadFromReader(reader interfaceList.Reader) error {
	rec.dreader = reader
	Reader := bufio.NewReader(reader.GetReader())
	rec.VideoList = make([]string, 0, 500)
	for {
		if bytes, isPre, err := Reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			} else if isPre == false {
				return nil
			}
		} else {
			line := rec.bytesToString(bytes)
			if line == "#EXT-X-ENDLIST" {
				return nil
			}
			if index := strings.Index(line, "#EXTINF:"); index == -1 {
				continue
			} else {
				if durlen, err := strconv.ParseFloat(line[8:len(line)-1], 64); err != nil {
					continue
				} else {
					rec.duration += durlen
					if bytes, _, err := Reader.ReadLine(); err == nil && len(bytes) != 0 {
						rec.VideoList = append(rec.VideoList, string(bytes))
					}
				}
			}
		}
	}
	return nil
}

func (rec *Resovler) WriteVideo(writer io.Writer, routine int) (int64, error) {
	signalToSwitch := make(chan int)
	for r := 0; r < routine; r++ {
		if r == 0 {
			go goroutine(r, rec.VideoList[0:len(rec.VideoList)/routine], signalToSwitch)
		} else if r <= routine-1 {
			go goroutine(r, rec.VideoList[r*(len(rec.VideoList)/routine):(r+1)*(len(rec.VideoList)/routine)], signalToSwitch)
		} else {
			go goroutine(r, rec.VideoList[r*(len(rec.VideoList)/routine):], signalToSwitch)
		}
	}
	t := 0
	for {
		if t >= routine {
			break
		}
		t += <-signalToSwitch
	}
	var file *os.File
	for i := 0; i < routine; i++ {
		sprintln := fmt.Sprint("swap", i)
		file, _ = os.OpenFile(sprintln, os.O_RDONLY, 0666)
		io.Copy(writer, file)
		file.Close()
		os.Remove(sprintln)
	}
	return 0, nil
}

func goroutine(fileTag int, workList []string, signalToSwitch chan<- int) {
	b := make([]byte, 33)
	sprintln := fmt.Sprint("swap", fileTag)
	file, err := os.OpenFile(sprintln, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := 0; i < len(workList); i++ {
		resp, err := http.Get(workList[i])
		if err != nil {
			panic(err)
		}
		resp.Body.Read(b)
		io.Copy(file, resp.Body)
	}
	signalToSwitch <- 1
	fmt.Printf("[第%d协程-任务完成]\n", fileTag+1)
}

func (rec *Resovler) FormatVideoTime(format string, desWriter io.Writer) error {
	_, err := fmt.Fprintf(desWriter, format, rec.duration)
	return err
}

func (rec *Resovler) Close() error {
	return rec.dreader.Close()
}
