package app

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"dont-genshin-bot/utils"
	"github.com/bwmarrin/discordgo"
)

const CmdKey = "!dgb"

var genshinPlayer sync.Map

func SettingHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	// 不處理Bot的消息和DM
	if e.Author.Bot ||
		e.GuildID == "" ||
		e.ChannelID == "" {
		return
	}
	arr := strings.Split(e.Content, " ")
	if arr[0] != CmdKey {
		return
	}
	isOwner, err := utils.IsOwner(s, e.GuildID, e.Author.ID)
	if err != nil {
		log.Println(err)
		return
	} else if !isOwner {
		_, _ = s.ChannelMessageSend(e.ChannelID, "Unauthorized")
		return
	}
	if len(arr) == 1 {
		botHelp(s, e.ChannelID)
		return
	}
	r := arr[1:]
	if len(r) >= 1 {
		switch {
		case r[0] == "help" && len(r) == 1:
			botHelp(s, e.ChannelID)
		case r[0] == "setting" && len(r) == 1:
			st, ok := utils.GetSetting(e.GuildID)
			if !ok {
				_, _ = s.ChannelMessageSend(e.ChannelID, "Setting not found.")
				return
			}
			_, _ = s.ChannelMessageSendEmbed(e.ChannelID, &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0x00ff00,
				Title:  "Now settings",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "mainChannel",
						Value: func() string {
							if st.MainChannelId != "" {
								return st.MainChannelId
							}
							return "Not set"
						}(),
						Inline: false,
					},
					{
						Name: "sendMessage",
						Value: func() string {
							if st.SendMessage != "" {
								return st.SendMessage
							}
							return "Not set"
						}(),
						Inline: false,
					},
					{
						Name: "defaultUserRole",
						Value: func() string {
							if st.DefaultUserRoleId != "" {
								return st.DefaultUserRoleId
							}
							return "Not set"
						}(),
						Inline: false,
					},
					{
						Name: "banRole",
						Value: func() string {
							if st.BanRoleId != "" {
								return st.BanRoleId
							}
							return "Not set"
						}(),
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			})
		case r[0] == "setting" && len(r) == 3 && r[1] == "mainChannel":
			st, ok := utils.GetSetting(e.GuildID)
			if !ok {
				st = utils.BotSetting{}
			}
			cid := r[2]
			if cid == "now" {
				cid = e.ChannelID
			}
			st.MainChannelId = cid
			utils.SetSetting(e.GuildID, st)
			_, _ = s.ChannelMessageSend(e.ChannelID, "Successfully updated.")
		case r[0] == "setting" && len(r) == 3 && r[1] == "defaultUserRole":
			st, ok := utils.GetSetting(e.GuildID)
			if !ok {
				st = utils.BotSetting{}
			}
			rid := r[2]
			if rid == "null" {
				rid = ""
			}
			st.DefaultUserRoleId = rid
			utils.SetSetting(e.GuildID, st)
			_, _ = s.ChannelMessageSend(e.ChannelID, "Successfully updated.")
		case r[0] == "setting" && len(r) == 3 && r[1] == "banRole":
			st, ok := utils.GetSetting(e.GuildID)
			if !ok {
				st = utils.BotSetting{}
			}
			st.BanRoleId = r[2]
			utils.SetSetting(e.GuildID, st)
			_, _ = s.ChannelMessageSend(e.ChannelID, "Successfully updated.")
		case r[0] == "setting" && len(r) >= 3 && r[1] == "sendMessage":
			st, ok := utils.GetSetting(e.GuildID)
			if !ok {
				st = utils.BotSetting{}
			}
			var msg string
			if r[2] == "default" {
				msg = "Don't play Genshin!"
			} else {
				msg = strings.Join(r[2:], " ")
			}
			st.SendMessage = msg
			utils.SetSetting(e.GuildID, st)
			_, _ = s.ChannelMessageSend(e.ChannelID, "Successfully updated.")
		}
	}
}

func DontGenshinHandler(s *discordgo.Session, e *discordgo.PresenceUpdate) {
	ids := utils.GuildIds()
	_, hp := genshinPlayer.Load(e.User.ID)
	if playingGenshin(e.Activities) {
		for _, id := range ids {
			if e.GuildID != id {
				continue
			}
			_, err := s.GuildMember(id, e.User.ID)
			if err == nil {
				setting, ok := utils.GetSetting(id)
				if ok {
					if setting.SendMessage != "" && !hp {
						msg := fmt.Sprintf("<@%s> %s", e.User.ID, setting.SendMessage)
						if setting.MainChannelId != "" {
							_, _ = s.ChannelMessageSend(setting.MainChannelId, msg)
						}
					}
					if !hp {
						ban(s, e.User.ID, setting)
					}
				}
			}
		}
	} else if hp {
		for _, id := range ids {
			_, err := s.GuildMember(id, e.User.ID)
			if err == nil {
				setting, ok := utils.GetSetting(id)
				if ok {
					unBan(s, e.User.ID, setting)
				}
			}
		}
	}
}

// playingGenshin 是否正在游玩原神
func playingGenshin(activities []*discordgo.Game) bool {
	for _, game := range activities {
		if game.Type == discordgo.GameTypeGame && utils.IsGenshin(game) {
			return true
		}
	}
	return false
}

func unBan(s *discordgo.Session, userId string, setting utils.BotSetting) {
	if setting.DefaultUserRoleId != "" {
		_ = s.GuildMemberRoleAdd(setting.GuildId, userId, setting.DefaultUserRoleId)
	}
	if setting.BanRoleId != "" {
		_ = s.GuildMemberRoleRemove(setting.GuildId, userId, setting.BanRoleId)
	}
	genshinPlayer.Delete(userId)
}

func ban(s *discordgo.Session, userId string, setting utils.BotSetting) {
	if setting.DefaultUserRoleId != "" {
		_ = s.GuildMemberRoleRemove(setting.GuildId, userId, setting.DefaultUserRoleId)
	}
	if setting.BanRoleId != "" {
		_ = s.GuildMemberRoleAdd(setting.GuildId, userId, setting.BanRoleId)
	}
	genshinPlayer.Store(userId, struct{}{})
}

// botHelp 發送機器人幫助信息
func botHelp(s *discordgo.Session, channelId string) {
	_, _ = s.ChannelMessageSendEmbed(channelId, &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00,
		Title:       "Don't Genshin Bot",
		Description: "幫助戒原神Discord機器人, 自動Ban掉正在玩原神的用戶",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "help",
				Value:  "獲取機器人指令列表.",
				Inline: false,
			},
			{
				Name:   "setting",
				Value:  "獲取當前伺服器設定.",
				Inline: false,
			},
			{
				Name:   "setting mainChannel <channel_Id | now>",
				Value:  "設定發送消息的文字頻道的ID, 未設定不發送消息. \n設定為now則為當前文字頻道.",
				Inline: false,
			},
			{
				Name:   "setting sendMessage <message | default>",
				Value:  "游玩原神后對用戶發送的消息, 未設定不發送消息. \n設定為default則為`Don't play Genshin!`.",
				Inline: false,
			},
			{
				Name:   "setting defaultUserRole <role_Id | null>",
				Value:  "設定默認用戶的身分組ID, Ban之後會移除, 解除Ban后會重新賦予. \n設定為null則為無默認身分組.",
				Inline: false,
			},
			{
				Name:   "setting banRole <role_Id>",
				Value:  "設定Ban掉用戶之後賦予的身分組的ID.",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
