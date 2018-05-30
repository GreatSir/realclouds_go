package middleware

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-ego/gse"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

const (
	//DEFAULT_DICT_DIR *
	DEFAULT_DICT_DIR = "dict_data/dict"

	//USER_DICT_PATH *
	USER_DICT_PATH = "/tmp/userdict.txt"
)

//DrityWord *
type DrityWord struct {
	DefaultDictDir string
	UserDictPath   string
	DrityWordMap   *map[string]string
	Gorm           *gorm.DB
	Segmenter      *gse.Segmenter
	Mutex          sync.RWMutex
}

//MwDrityWord Drity word middleware
func (d *DrityWord) MwDrityWord(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("drityword", d)
		return next(c)
	}
}

//NewDrityWord *
func NewDrityWord(db *gorm.DB, userDictPath ...string) (drityWord *DrityWord, err error) {
	userDict := USER_DICT_PATH

	if len(userDictPath) > 0 {
		userDict = strings.TrimSpace(userDictPath[0])
	}

	err = db.AutoMigrate(&DrityWordDB{}).Error
	if nil != err {
		return nil, err
	}

	drityWord = &DrityWord{
		UserDictPath:   strings.TrimSpace(userDict),
		DefaultDictDir: strings.TrimSpace(DEFAULT_DICT_DIR),
		Segmenter:      new(gse.Segmenter),
		Gorm:           db,
	}

	_, drityWords := FindDrityWords(db)

	drityWordMap := make(map[string]string)

	for _, drityWord := range drityWords {
		drityWordMap[drityWord.MD5] = drityWord.Value
	}

	drityWord.DrityWordMap = &drityWordMap

	if err = drityWord.WriteDrityWord(); nil != err {
		return
	}

	return
}

//WriteDrityWord *
func (d *DrityWord) WriteDrityWord() error {
	f, err := os.OpenFile(d.UserDictPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range *d.DrityWordMap {
		if len(v) > 0 {
			_, err := f.WriteString(v + " 100000\n")
			if nil != err {
				return err
			}
		}
	}

	if err := d.ReloadDict(); nil != err {
		return err
	}
	return nil
}

//ReloadDict *
func (d *DrityWord) ReloadDict() error {

	pd := utils.GetProjectDir()

	paths, err := utils.WalkPaths(utils.ArrayPath(pd, d.DefaultDictDir))
	if nil != err {
		return err
	}

	paths = append(paths, d.UserDictPath)

	pathsStr := strings.Join(paths, ",")

	fmt.Printf("Reload dict,paths: %s\n", pathsStr)

	if err := d.Segmenter.LoadDict(pathsStr); nil != err {
		return err
	}

	return nil
}

//Subscription *
func (d *DrityWord) Subscription(rPool *redis.Pool) error {

	ctx, cancel := context.WithCancel(context.Background())

	err := ListenPubSubChannels(ctx, rPool,
		func() error {
			fmt.Printf("\nDrity word subscription start...\n\n")
			return nil
		},
		func(channel string, message []byte) error {
			msgStr := string(bytes.TrimSpace(message))
			msgStr = strings.ToLower(msgStr)
			channel = strings.TrimSpace(channel)

			if len(msgStr) > 0 && msgStr == "up" {
				fmt.Printf("channel: %s, message: %v\n", channel, msgStr)

				if DRITYWORD_UP_SUBSCRIPTION_KEY == channel {
					_, drityWords := FindDrityWords(d.Gorm)

					drityWordMap := make(map[string]string)

					for _, drityWord := range drityWords {
						drityWordMap[drityWord.MD5] = drityWord.Value
					}

					d.DrityWordMap = &drityWordMap

					if err := d.WriteDrityWord(); nil != err {
						cancel()
						fmt.Printf("\nDrity word subscription error: %v\n", err.Error())
						return err
					}

					fmt.Printf("\nReload drity word at: %v\n", utils.DateToStr(time.Now()))
					cancel()
					return nil
				}
			} else {
				fmt.Printf("nil data\n")
			}
			return nil
		},
		nil,
		[]string{DRITYWORD_UP_SUBSCRIPTION_KEY}, []string{})

	if nil != err {
		fmt.Printf("\nDrity word subscription error: %v\n", err.Error())
		return err
	}
	return nil
}
