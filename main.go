package main

import (
    "fmt"
    "encoding/json"
    //"strings"
    "os"
    "os/signal"
    "syscall"
    "regexp"
    "github.com/bwmarrin/discordgo"
)

type Embed struct {
    Message string
    RulesTitle string
    RulesText string
    RulesText2 string
    QuestionsTitle string
    QuestionsText string
    BugsTitle string
    BugsText string
    Image string
}

type Config struct {
    AdminID string
    ServerID string
    LockedRoleID string
    Token string
    WelcomeChannel string
    WelcomeEmbed Embed
}

var config = readConfig()

func readConfig() Config {
    file, _ := os.Open("config.json")
    conf := Config{}
    _ = json.NewDecoder(file).Decode(&conf)
    file.Close()
    return conf
}



func main() {
    dg, err := discordgo.New("Bot " + config.Token)
    if err != nil {
        fmt.Println("error: ", err)
        return
    }

    dg.AddHandler(genericReply)
    dg.AddHandler(onJoin)
    err = dg.Open()
    if err != nil {
        fmt.Println("no connection, ", err)
        return
    }

    fmt.Println("bot running. selphyWoo")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

func onJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
    if !member.User.Bot {
        s.GuildMemberRoleAdd(config.ServerID, member.User.ID, config.LockedRoleID)
    }
}

func unlockUser(s *discordgo.Session, id string) {
    s.GuildMemberRoleRemove(config.ServerID, id, config.LockedRoleID)
}

func genericReply(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    if m.Author.ID == config.AdminID {
        replyGodmode(s, m)
    } else if m.ChannelID == config.WelcomeChannel {
        s.ChannelMessageDelete(m.ChannelID, m.ID)
        if m.Content == "!accept" {
            unlockUser(s, m.Author.ID)
        }
        return
    }

    // In case this doesn’t work with your font: the last character is a winking emoji.
    winks, _ := regexp.MatchString("([()|DoO];|;[()|DoOpP]|:wink:|😉)", m.Content)
    if winks {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> faggot", m.Author.ID))
        s.ChannelMessageDelete(m.ChannelID, m.ID)
        return
    }

    // As per our majesty’s command:
    if m.Content == "\\o" {
        s.ChannelMessageSend(m.ChannelID, "o/")
    } else if m.Content == "o/" {
        s.ChannelMessageSend(m.ChannelID, "\\o")
    }
}


// Admin stuff down here. This is very server-specific

func replyGodmode(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Content == "print_rules()" {
        embedColor := 0xffb90f // kageru gold
        embed := &discordgo.MessageEmbed{
                Author:      &discordgo.MessageEmbedAuthor{},
                Color:       embedColor,
                Description: config.WelcomeEmbed.Message,
                Fields: []*discordgo.MessageEmbedField{
                    &discordgo.MessageEmbedField{
                        Name:   config.WelcomeEmbed.QuestionsTitle,
                        Value:  config.WelcomeEmbed.QuestionsText,
                        Inline: true,
                    },
                    &discordgo.MessageEmbedField{
                        Name:   config.WelcomeEmbed.BugsTitle,
                        Value:  fmt.Sprintf(config.WelcomeEmbed.BugsText, config.AdminID),
                        Inline: true,
                    },
                },
                Thumbnail: &discordgo.MessageEmbedThumbnail{
                    URL: config.WelcomeEmbed.Image,
                },
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        embed = &discordgo.MessageEmbed{
                Author:      &discordgo.MessageEmbedAuthor{},
                Color:       embedColor,
                Fields: []*discordgo.MessageEmbedField{
                    &discordgo.MessageEmbedField{
                        Name:   config.WelcomeEmbed.RulesTitle,
                        Value:  config.WelcomeEmbed.RulesText,
                        Inline: true,
                    },
                },
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        embed = &discordgo.MessageEmbed{
                Author:      &discordgo.MessageEmbedAuthor{},
                Color:       embedColor,
                Fields: []*discordgo.MessageEmbedField{
                    &discordgo.MessageEmbedField{
                        Name:   config.WelcomeEmbed.RulesTitle,
                        Value:  config.WelcomeEmbed.RulesText2,
                        Inline: true,
                    },
                },
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
}

