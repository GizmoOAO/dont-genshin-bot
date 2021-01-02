package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	settings sync.Map
)

type BotSetting struct {
	GuildId string
	// MainChannelId 發送消息的文字頻道的ID, 未設定不發送消息
	MainChannelId string
	// SendMessage 游玩原神后對用戶發送的消息, 未設定不發送消息
	SendMessage string
	// DefaultUserRoleId 默認用戶的身分組ID, Ban之後會移除, 解除Ban后會重新賦予
	DefaultUserRoleId string
	// BanRoleId Ban掉用戶之後賦予的身分組的ID
	BanRoleId string
}

func LoadSettings() {
	p := os.Getenv("SETTINGS")
	if _, err := os.Stat(p); err != nil {
		if err = os.MkdirAll(p, 0666); err != nil {
			panic(err)
		}
	}
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, _ := filepath.Abs(filepath.Join(p, info.Name()))
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		var setting BotSetting
		if err = json.Unmarshal(data, &setting); err != nil {
			return err
		}
		guildId := strings.ReplaceAll(info.Name(), ".json", "")
		settings.Store(guildId, setting)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func saveSetting(guildId string, s BotSetting) {
	s.GuildId = guildId
	filename := filepath.Join(os.Getenv("SETTINGS"), guildId+".json")
	data, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return
	}
	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		log.Println(err)
	}
}

func GetSetting(guildId string) (setting BotSetting, ok bool) {
	s, ok := settings.Load(guildId)
	if !ok {
		return
	}
	return s.(BotSetting), true
}

func SetSetting(guildId string, s BotSetting) {
	settings.Store(guildId, s)
	saveSetting(guildId, s)
}

func GuildIds() []string {
	var ids []string
	settings.Range(func(k, v interface{}) bool {
		ids = append(ids, k.(string))
		return false
	})
	return ids
}
