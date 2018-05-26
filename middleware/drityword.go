package middleware

import (
	"os"
	"strings"
	"sync"

	"github.com/labstack/echo"
)

const (
	//USER_DICT_PATH *
	USER_DICT_PATH = "/tmp/userdict.txt"
)

//DrityWord *
type DrityWord struct {
	UserDictPath string
	DrityWordMap *map[string]string
	Mutex        sync.RWMutex
}

//MwDrityWord Drity word middleware
func (d *DrityWord) MwDrityWord(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		d.Mutex.Lock()
		defer d.Mutex.Unlock()
		c.Set("drityword", d)
		return next(c)
	}
}

//NewDrityWord *
func NewDrityWord(drityWordMap *map[string]string, userDictPath ...string) (*DrityWord, error) {
	userDict := USER_DICT_PATH

	if len(userDictPath) > 0 {
		userDict = strings.TrimSpace(userDictPath[0])
	}

	if err := writeDrityWord(userDict, drityWordMap); nil != err {
		return nil, err
	}

	return &DrityWord{
		UserDictPath: strings.TrimSpace(userDict),
		DrityWordMap: drityWordMap,
	}, nil
}

//writeDrityWord *
func writeDrityWord(userDictPath string, drityWordMap *map[string]string) error {
	f, err := os.OpenFile(userDictPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, drityWord := range *drityWordMap {
		if len(drityWord) > 0 {
			_, err := f.WriteString(drityWord + "\n")
			if nil != err {
				return err
			}
		}
	}

	return nil
}
