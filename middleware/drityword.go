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
func NewDrityWord(drityWordMap *map[string]string, userDictPath ...string) (drityWord *DrityWord, err error) {
	userDict := USER_DICT_PATH

	if len(userDictPath) > 0 {
		userDict = strings.TrimSpace(userDictPath[0])
	}

	drityWord = &DrityWord{
		UserDictPath: strings.TrimSpace(userDict),
		DrityWordMap: drityWordMap,
	}

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

	for _, drityWord := range *d.DrityWordMap {
		if len(drityWord) > 0 {
			_, err := f.WriteString(drityWord + "\n")
			if nil != err {
				return err
			}
		}
	}
	return nil
}
