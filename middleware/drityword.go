package middleware

import (
	"os"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

const (
	//USER_DICT_PATH *
	USER_DICT_PATH = "/tmp/userdict.txt"
)

//Drityword *
type Drityword struct {
	UserDictPath string
	DrityWordMap *map[string]string
	Gorm         *gorm.DB
	Mutex        sync.RWMutex
}

//MwDrityWord Drity word middleware
func (d *Drityword) MwDrityWord(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		d.Mutex.Lock()
		defer d.Mutex.Unlock()
		c.Set("drityword", d)
		return next(c)
	}
}

//NewDrityWord *
func NewDrityWord(db *gorm.DB, userDictPath ...string) (drityWord *Drityword, err error) {
	userDict := USER_DICT_PATH

	if len(userDictPath) > 0 {
		userDict = strings.TrimSpace(userDictPath[0])
	}

	dw := &Drityword{
		UserDictPath: strings.TrimSpace(userDict),
	}

	_, drityWords := FindDrityWords(db)

	drityWordMap := make(map[string]string)

	for _, drityWord := range drityWords {
		drityWordMap[drityWord.MD5] = drityWord.Value
	}

	dw.DrityWordMap = &drityWordMap

	if err = dw.WriteDrityWord(); nil != err {
		return
	}

	return
}

//WriteDrityWord *
func (d *Drityword) WriteDrityWord() error {
	f, err := os.OpenFile(d.UserDictPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range *d.DrityWordMap {
		if len(v) > 0 {
			_, err := f.WriteString(v + "\n")
			if nil != err {
				return err
			}
		}
	}
	return nil
}
