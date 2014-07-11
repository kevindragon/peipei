package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const DistPath string = `C:\home\peipei\2014最新商家资料\鹏荣达开票资料`
const LongTimeFormat string = "2006-01-02 15:04:05"
const AbbrevTimeFormat string = "2006-1-2 15:4:5"

func getFileTime() syscall.Filetime {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	st, err := time.ParseInLocation(LongTimeFormat,
		"2014-06-23 10:09:29", loc)
	if err != nil {
		fmt.Println(err)
	}
	et, err := time.ParseInLocation(LongTimeFormat,
		"2014-06-23 10:09:30", loc)
	if err != nil {
		fmt.Println(err)
	}

	gapTimestamp := et.Unix() - st.Unix()
	rand.Seed(time.Now().UnixNano())

	randGap := rand.Int63n(gapTimestamp)

	distTime := uniformTime(time.Unix(st.Unix()+randGap, 0))

	ftime := syscall.NsecToFiletime(distTime.UnixNano())

	return ftime
}

func modifyFileTime(path string, ftime syscall.Filetime) error {
	fd, err := syscall.Open(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(path, err)
	}

	err = syscall.SetFileTime(fd, &ftime, &ftime, &ftime)
	if err != nil {
		fmt.Println(path, err)
	}
	syscall.Close(fd)

	return err
}

func fileHandler(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		fmt.Println(path, "is a dir")
		return nil
	}

	ftime := getFileTime()
	modifyFileTime(path, ftime)

	return nil
}

func uniformTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	rand.Seed(time.Now().UnixNano())
	startYmdStr := fmt.Sprintf("%d-%d-%d 09:30:00", t.Year(), t.Month(), t.Day())
	endYmdStr := fmt.Sprintf("%d-%d-%d 18:30:00", t.Year(), t.Month(), t.Day())

	daySt, _ := time.ParseInLocation(AbbrevTimeFormat, startYmdStr, loc)
	dayEt, _ := time.ParseInLocation(AbbrevTimeFormat, endYmdStr, loc)
	dayGap := dayEt.Unix() - daySt.Unix()
	randGap := rand.Int63n(dayGap)

	//fmt.Println(t, startYmdStr, daySt, dayEt, dayGap, randGap)

	if t.Unix() < daySt.Unix() || t.Unix() > dayEt.Unix() {
		return time.Unix(daySt.Unix()+randGap, 0)
	}
	return t
}

func main() {
	filepath.Walk(DistPath, fileHandler)
}
