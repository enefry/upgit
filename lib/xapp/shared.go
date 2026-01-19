package xapp

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
	"github.com/pluveto/upgit/lib/xpath"
)

const UserAgent = "UPGIT/0.2"
const DefaultBranch = "master"

// case insensitive
const ClipboardPlaceholder = ":clipboard"
const ClipboardFilePlaceholder = ":clipboard-file"

var MaxUploadSize = int64(5 * 1024 * 1024)
var ConfigFilePath string

func Rename(path string, time time.Time) (ret string) {

	base := xpath.Basename(path)
	ext := filepath.Ext(path)
	md5HashStr := fmt.Sprintf("%x", md5.Sum([]byte(base)))
	replacementList := []string{
		"{year}", time.Format("2006"),
		"{month}", time.Format("01"),
		"{day}", time.Format("02"),
		"{hour}", time.Format("15"),
		"{minute}", time.Format("04"),
		"{second}", time.Format("05"),
		"{unixts}", fmt.Sprint(time.Unix()),
		"{unixtsms}", fmt.Sprint(time.UnixMicro()),
		"{ext}", ext,
		"{fullname}", base + ext,
		"{filename}", base,
		"{fname}", base,
		"{filenamehash}", md5HashStr,
		"{fnamehash}", md5HashStr,
		"{fnamehash4}", md5HashStr[:4],
		"{fnamehash8}", md5HashStr[:8],
	}

	if len(AppCfg.HmacKey) > 0 && len(AppCfg.HmacFormat) > 0 && strings.Contains(AppCfg.Rename, "{hmac}") {
		r := strings.NewReplacer(replacementList...)
		param := r.Replace(AppCfg.HmacFormat)
		h := hmac.New(sha256.New, []byte(AppCfg.HmacKey))
		h.Write([]byte(param))
		hmacHashStr := hex.EncodeToString(h.Sum(nil))
		if AppCfg.HmacLen > 0 && AppCfg.HmacLen < len(hmacHashStr) {
			hmacHashStr = hmacHashStr[:AppCfg.HmacLen]
		}
		replacementList = append(replacementList, "{hmac}", hmacHashStr)
	}

	r := strings.NewReplacer(replacementList...)
	ret = r.Replace(AppCfg.Rename)
	return
}
func ReplaceUrl(path string) (ret string) {
	var rules []string
	for k, v := range AppCfg.Replacements {
		rules = append(rules, k, v)
	}
	r := strings.NewReplacer(rules...)
	ret = r.Replace(path)
	return
}

func LoadUploaderConfig[T any](uploaderId string) (ret T, err error) {
	var mCfg map[string]interface{}
	bytes, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		return
	}
	err = toml.Unmarshal(bytes, &mCfg)
	if err != nil {
		return
	}
	cfgMap := mCfg["uploaders"].(map[string]interface{})[uploaderId]
	var cfg_ T
	mapstructure.Decode(cfgMap, &cfg_)
	ret = cfg_
	return
}
