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

type Config struct {
    Token string
    Welcome string
    WelcomeChannel string
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

func genericReply(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    winks, _ := regexp.MatchString("([()|DoO];|;[()|DoOpP])", m.Content)
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
    if m.Content == "test_welcome()" {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(config.Welcome, m.Author.ID))
    }
}

func onJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
    fmt.Println("user joined")

    fmt.Println(member.User.Bot)
    if !member.User.Bot {
        s.ChannelMessageSend(config.WelcomeChannel, fmt.Sprintf(config.Welcome, member.User.ID))
    }
    fmt.Println(fmt.Sprintf(config.Welcome, member.User.ID))
}

