package main

import (
    "github.com/bwmarrin/discordgo"
    "log"
    "fmt"
)

func unlockUser(s *discordgo.Session, id string) {
    s.GuildMemberRoleRemove(config.ServerID, id, config.LockedRoleID)
    log.Printf("Removed lock from user: %s", userToString(getUser(s, id)))
}

func userToString(u *discordgo.User) string {
    return fmt.Sprintf("%s#%s (ID: %s)", u.Username, u.Discriminator, u.ID)
}

func roleName(s *discordgo.State, rid string) string {
    role, _ := s.Role(config.ServerID, rid)
    return role.Name
}

func channelToString(c *discordgo.Channel) string {
    return fmt.Sprintf("%s (ID: %s) on %s", c.Name, c.ID, c.GuildID)
}

func messageToString(m *discordgo.Message) string {
    return fmt.Sprintf("<%s#%s>: %s", m.Author.Username, m.Author.Discriminator, m.Content)
}

func getChannel(s *discordgo.State, cid string) *discordgo.Channel {
    channel, _ := s.Channel(cid)
    return channel
}

func getUser(s *discordgo.Session, uid string) *discordgo.User {
    user, _ := s.User(uid)
    return user
}

