package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	gounicode "unicode"
	"unicode/utf8"

	pinyin "github.com/mozillazg/go-pinyin"
	"github.com/pborman/uuid"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	PathSeparator = string(os.PathSeparator) //PathSeparator
	DevNull       = os.DevNull               //DevNull
)

func GetENV(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func GetENVToBool(key string) bool {
	envStr := strings.TrimSpace(os.Getenv(key))
	boo, err := StringUtils(envStr).Bool()
	if nil != err {
		boo = false
	}
	return boo
}

func GetENVToInt(key string) (int, error) {
	envStr := strings.TrimSpace(os.Getenv(key))
	return StringUtils(envStr).Int()
}

func GetENVToInt64(key string) (int64, error) {
	envStr := strings.TrimSpace(os.Getenv(key))
	return StringUtils(envStr).Int64()
}

func GetBinPath() (string, error) {
	file, _ := exec.LookPath(os.Args[0])
	bin_path, _ := filepath.Abs(file)
	bin_path, err := filepath.EvalSymlinks(bin_path)
	if nil != err {
		return "", err
	}
	return bin_path, nil
}

func GetBinDir() string {
	binPath, _ := GetBinPath()
	return filepath.Dir(binPath)
}

func GetProjectDir() string {
	dirs := strings.Split(GetBinDir(), PathSeparator)
	pjtPath := strings.Join(dirs, PathSeparator)
	src := PathSeparator + "src"
	srcIndex := strings.Index(pjtPath, src)
	if srcIndex != -1 {
		pjtPath = pjtPath[:srcIndex]
	}
	pjtPath = strings.TrimRight(pjtPath, PathSeparator)
	return pjtPath
}

func ArrayPath(path ...string) string {
	return strings.Join(path, PathSeparator)
}

//RegGob 将自定义 Struct 注册 Gob
func RegGob(o ...interface{}) {
	for _, v := range o {
		gob.Register(v)
	}
	mapGob := make(map[string]interface{})
	gob.Register(mapGob)
}

func DateToStr(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func ValidTime(value string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", value)
	if nil != err {
		return false
	}
	return true
}

func FormatDate(time time.Time, format string) string {
	return time.Format(format)
}

func FormatDateToStdTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", value)
}

func MkdirByFile(file string) error {
	fileDir := filepath.Dir(file)
	if !IsDir(fileDir) {
		if err := os.Mkdir(fileDir, os.ModePerm); nil != err {
			return err
		}
	}
	return nil
}

func WritePidFile(file, pid string) error {
	if err := MkdirByFile(file); nil != err {
		return err
	}

	pidfile, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer pidfile.Close()
	_, err = io.WriteString(pidfile, pid)
	if err != nil {
		return err
	}
	return nil
}

type StringUtils string

func (s *StringUtils) Set(v string) {
	if v != "" {
		*s = StringUtils(v)
	} else {
		s.Clear()
	}
}

func (s *StringUtils) Clear() {
	*s = StringUtils(0x1E)
}

func (s StringUtils) Exist() bool {
	return string(s) != string(0x1E)
}

func (s StringUtils) Bool() (bool, error) {
	v, err := strconv.ParseBool(s.String())
	return bool(v), err
}

func (s StringUtils) Float32() (float32, error) {
	v, err := strconv.ParseFloat(s.String(), 32)
	return float32(v), err
}

func (s StringUtils) Float64() (float64, error) {
	return strconv.ParseFloat(s.String(), 64)
}

func (s StringUtils) Int() (int, error) {
	v, err := strconv.ParseInt(s.String(), 10, 32)
	return int(v), err
}

func (s StringUtils) Int8() (int8, error) {
	v, err := strconv.ParseInt(s.String(), 10, 8)
	return int8(v), err
}

func (s StringUtils) Int16() (int16, error) {
	v, err := strconv.ParseInt(s.String(), 10, 16)
	return int16(v), err
}

func (s StringUtils) Int32() (int32, error) {
	v, err := strconv.ParseInt(s.String(), 10, 32)
	return int32(v), err
}

func (s StringUtils) Int64() (int64, error) {
	v, err := strconv.ParseInt(s.String(), 10, 64)
	return int64(v), err
}

func (s StringUtils) Uint() (uint, error) {
	v, err := strconv.ParseUint(s.String(), 10, 32)
	return uint(v), err
}

func (s StringUtils) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(s.String(), 10, 8)
	return uint8(v), err
}

func (s StringUtils) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(s.String(), 10, 16)
	return uint16(v), err
}

func (s StringUtils) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(s.String(), 10, 32)
	return uint32(v), err
}

func (s StringUtils) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(s.String(), 10, 64)
	return uint64(v), err
}

func (s StringUtils) ToTitleLower() string {
	str := strings.ToLower(s.String()[:1]) + s.String()[1:]
	return str
}

func (s StringUtils) ToTitleUpper() string {
	str := strings.ToUpper(s.String()[:1]) + s.String()[1:]
	return str
}

func (s StringUtils) ContainsBool(sep string) bool {
	index := strings.Index(s.String(), sep)
	return index > -1
}

func (s StringUtils) String() string {
	if s.Exist() {
		return string(s)
	}
	return ""
}

func (s StringUtils) MD5() string {
	m := md5.New()
	m.Write([]byte(s.String()))
	return hex.EncodeToString(m.Sum(nil))
}

func (s StringUtils) SHA1() string {
	sha := sha1.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) SHA256() string {
	sha := sha256.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) SHA512() string {
	sha := sha512.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) HMAC_SHA1(key string) string {
	mc := hmac.New(sha1.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) HMAC_SHA256(key string) string {
	mc := hmac.New(sha256.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) HMAC_SHA512(key string) string {
	mc := hmac.New(sha512.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) StdBase64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(s.String()))
}

func (s StringUtils) Base64Encode() string {
	return base64.URLEncoding.EncodeToString([]byte(s.String()))
}

func (s StringUtils) Base64Decode() (string, error) {

	inputVal := s.String()

	seg, err := url.PathUnescape(inputVal)
	if nil != err {
		return "", err
	}

	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	v, err := base64.URLEncoding.DecodeString(seg)
	return string(v), err
}

func (s StringUtils) PinYinToHTML() string {
	var buf bytes.Buffer

	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	_, source := TrimHTML(s.String())

	for i, r := range source {
		py := pinyin.SinglePinyin(r, a)
		_, ok := pinyin.PinyinDict[int(r)]
		if ok {
			txt := fmt.Sprintf("<ruby><rb>%v</rb><rp>(</rp><rt>%v</rt><rp>)</rp></ruby>", source[i:i+3], py[0])
			buf.WriteString(txt)
		} else {
			txt := fmt.Sprintf("<span>%v</span>", py[0])
			buf.WriteString(txt)
		}
	}

	return buf.String()
}

func (s StringUtils) PinYin() map[string][]string {

	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	a.Heteronym = true
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	_, source := TrimHTML(s.String())

	pinYinMap := make(map[string][]string)

	for i, r := range source {
		py := pinyin.SinglePinyin(r, a)
		_, ok := pinyin.PinyinDict[int(r)]
		if ok {
			pinYinStr := source[i : i+3]
			if _, ok := pinYinMap[pinYinStr]; !ok {
				pinYinMap[pinYinStr] = py
			}
		}
	}

	return pinYinMap
}

func (s StringUtils) CleanUP() string {

	source := s.String()
	source = strings.TrimSpace(source)
	source = strings.Replace(source, " ", "", -1)

	buf := bytes.Buffer{}
	for _, r := range source {
		if gounicode.Is(gounicode.Scripts["Han"], r) || gounicode.IsLetter(r) || gounicode.IsNumber(r) {
			str := fmt.Sprintf("%c", r)
			buf.WriteString(str)
		}
	}
	return buf.String()
}

func (s StringUtils) FileNameCleanUP() string {

	source := s.String()
	ext := path.Ext(source)

	s.Set(strings.TrimRight(source, ext))

	source = s.CleanUP()

	return fmt.Sprintf("%s%s", source, ext)
}

//RandURL *
func (s StringUtils) RandURL() string {
	path := strings.TrimSpace(s.String())

	var buf bytes.Buffer
	buf.WriteString(path)

	if strings.Contains(path, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}

	buf.WriteString("_=" + fmt.Sprintf("%v", RandInt64()))
	return buf.String()
}

//TrimHTML *
func (s StringUtils) TrimHTML() (int, string) {
	return TrimHTML(s.String())
}

//RegeMobileNo *
func RegeMobileNo(mobile string) bool {

	mobile = strings.TrimSpace(mobile)

	regx := regexp.MustCompile("((13[0-9])|(14[0-9])|(15[0-9])|(16[0-9])|(17[0-9])|(18[0-9])|(19[0-9]))\\d{8}")

	return regx.MatchString(mobile)
}

//TrimHTML *
func TrimHTML(src string) (int, string) {

	src = HTMLUnicodeToHTMLString(src)

	regx := regexp.MustCompile("(?i)<!--(.*?)-->")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("(?i)<style[\\S\\s]+?</style>")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("(?i)<script[\\S\\s]+?</script>")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("(?i)<(.*?)>")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("\\(((?i)[a-zà-ǜ]+\\))|(（(?i)[a-zà-ǜ]+）)|\\(((?i)[a-zà-ǜ]+）)|(（(?i)[a-zà-ǜ]+\\))")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("([\\r\\n])[\\s]+")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(quot|#34);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(amp|#38);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(lt|#60);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(gt|#62);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(nbsp|#160);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(iexcl|#161);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(cent|#162);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(pound|#163);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&(copy|#169);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("&#(\\d+);")
	src = regx.ReplaceAllString(src, "")

	regx = regexp.MustCompile("\\s+")
	src = regx.ReplaceAllString(src, "")

	return utf8.RuneCountInString(src), src
}

//HTMLUnicodeToHTMLString *
func HTMLUnicodeToHTMLString(htmlUnicode string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(htmlUnicode)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(htmlUnicode[i:])
			break
		}
		if htmlUnicode[i] == '\\' && htmlUnicode[i+1] == 'u' {
			hex := htmlUnicode[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(htmlUnicode[i:x])
			}
			i = x
		} else {
			buf.WriteByte(htmlUnicode[i])
			i++
		}
	}
	return buf.String()
}

func RandInt(start, end int64) int64 {
	if start >= end {
		return end
	}
	rnd := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return rnd.Int63n(end-start) + start
}

func RandInt64() int64 {
	rnd := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return rnd.Int63()
}

func GenCode6() string {
	rnd := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06d", rnd.Int63n(1000000))
	return vCode
}

func GenerateMacAddr() string {
	mac := []int64{
		0x00, 0xf8, 0x3e, RandInt(0x00, 0xff), RandInt(0x00, 0xff), RandInt(0x00, 0xff),
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func AESEncode(msg, key string) (string, error) {
	if len(key) == 16 {
		var iv = []byte(key)[:aes.BlockSize]
		c := make([]byte, len(msg))
		be, err := aes.NewCipher([]byte(key))
		if err != nil {
			return "", err
		}
		e := cipher.NewCFBEncrypter(be, iv)
		e.XORKeyStream(c, []byte(msg))
		b64 := base64.StdEncoding.EncodeToString(c)
		b64 = strings.Replace(b64, "/", "-", -1)
		return b64, nil
	} else {
		return "", fmt.Errorf("%s", "Key length is not equal to 16.")
	}
}

func AESDecode(enmsg, key string) (string, error) {
	if len(key) == 16 {
		enmsg = strings.Replace(enmsg, "-", "/", -1)
		msg, err := base64.StdEncoding.DecodeString(enmsg)
		if nil != err {
			return "", err
		}
		var iv = []byte(key)[:aes.BlockSize]
		d := make([]byte, len(msg))
		var bd cipher.Block
		bd, err = aes.NewCipher([]byte(key))
		if err != nil {
			return "", err
		}
		e := cipher.NewCFBDecrypter(bd, iv)
		e.XORKeyStream(d, msg)
		return string(d), nil
	} else {
		return "", fmt.Errorf("%s", "Key length is not equal to 16.")
	}
}

func (s StringUtils) GenerateRandStr32() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func ArrayToStringMap(values ...string) map[string]string {
	mapVal := make(map[string]string)
	for i := 0; i < len(values); i = i + 2 {
		k, v := values[i], values[i+1]
		mapVal[k] = v
	}
	return mapVal
}

func MergeStringMap(maps ...map[string]string) map[string]string {
	mapVal := make(map[string]string)
	for _, val := range maps {
		for k, v := range val {
			mapVal[k] = v
		}
	}
	return mapVal
}

func ToStr(value interface{}) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', 2, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', 2, 64)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(int64(v), 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

func ConvertUTF8(src []byte) ([]byte, error) {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), GetCharset("UTF-8").NewEncoder()))
	return data, err
}

func GenerateUUID() string {
	return uuid.NewRandom().String()
}

func ObjectToJSON(i interface{}) (string, error) {
	b, err := json.Marshal(i)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

func JSONToObject(data string, i interface{}) error {
	return json.Unmarshal([]byte(data), i)
}

func GetCharset(charset string) encoding.Encoding {
	switch strings.ToUpper(charset) {
	case "GB18030":
		return simplifiedchinese.GB18030
	case "GB2312", "HZ-GB2312":
		return simplifiedchinese.HZGB2312
	case "GBK":
		return simplifiedchinese.GBK
	case "BIG5":
		return traditionalchinese.Big5
	case "EUC-JP":
		return japanese.EUCJP
	case "ISO2022JP":
		return japanese.ISO2022JP
	case "SHIFTJIS":
		return japanese.ShiftJIS
	case "EUC-KR":
		return korean.EUCKR
	case "UTF8", "UTF-8":
		return encoding.Nop
	case "UTF16-BOM", "UTF-16-BOM":
		return unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	case "UTF16-BE-BOM", "UTF-16-BE-BOM":
		return unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	case "UTF16-LE-BOM", "UTF-16-LE-BOM":
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	case "UTF16", "UTF-16":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "UTF16-BE", "UTF-16-BE":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "UTF16-LE", "UTF-16-LE":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	//case "UTF32", "UTF-32":
	//	return simplifiedchinese.GBK
	default:
		return nil
	}
}

func IsFile(p string) bool {
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	} else {
		return fi.IsDir()
	}
}

func MkDirAll(dir string) error {
	if IsDir(dir) {
		if err := os.RemoveAll(dir); nil != err {
			return err
		}
	}
	if err := os.MkdirAll(dir, os.ModePerm); nil != err {
		return err
	}
	return nil
}

func GetFileSizeToUnit(fileSize int64) string {
	ff_size := float64(fileSize)
	var (
		fs   string
		pb_s float64 = 1024 << 40
		tb_s float64 = 1024 << 30
		gb_s float64 = 1024 << 20
		mb_s float64 = 1024 << 10
	)
	if ff_size > pb_s {
		f := ff_size / pb_s
		fs = ToStr(f) + " PB"
	} else if ff_size > tb_s {
		f := ff_size / tb_s
		fs = ToStr(f) + " TB"
	} else if ff_size > gb_s {
		f := ff_size / gb_s
		fs = ToStr(f) + " GB"
	} else if ff_size > mb_s {
		f := ff_size / mb_s
		fs = ToStr(f) + " MB"
	} else {
		f := ff_size / 1024
		fs = ToStr(f) + " KB"
	}
	return fs
}
