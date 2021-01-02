package utils

import (
	"github.com/bwmarrin/discordgo"
)

func IsGenshin(game *discordgo.Game) bool {
	switch game.Name {
	case "原神":
		return true
	case "Genshin Impact":
		return true
	}
	return false
}

func IsOwner(s *discordgo.Session, guildID string, userID string) (bool, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return false, err
	}
	return g.OwnerID == userID, nil
}
