package main

import (
    "fmt"
    "os/signal"
    "os"
    "syscall"
    "log"
    "github.com/bwmarrin/discordgo"
)

var config = readConfig()
var commands []Command

func main() {
    dg, err := discordgo.New("Bot " + config.Token)
    if err != nil {
        fmt.Println("error: ", err)
        return
    }

    dg.AddHandler(evaluateMessage)
    dg.AddHandler(onJoin)
    err = dg.Open()
    if err != nil {
        fmt.Println("no connection, ", err)
        return
    }

    f, err := os.OpenFile("selphybot.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("error opening log file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)

    // Moderation
    registerCommand(Command{Trigger: "^[^`]*([()|DoO];|;[()|DoOpP]|:wink:)[^`]*$", Output: "<@%s> Oboe!", DeleteInput: true, OutputIsReply: true, Type: CommandTypeRegex})

    // Misc commands
    registerCommand(Command{Trigger: "o/", Output: "\\o", Type: CommandTypeFullMatch})
    registerCommand(Command{Trigger: "\\o", Output: "o/", Type: CommandTypeFullMatch})
    registerCommand(Command{Trigger: "<:selphyDango:441001954542616576>", Output: ":notes: Dango, Dango, Dango, Dango, Dango Daikazoku :notes:", Type: CommandTypeFullMatch})


    registerCommand(Command{Trigger: "!welcome", OutputEmbed: getWelcomeEmbed(), Type: CommandTypeFullMatch, DMOnly: true})
    registerCommand(Command{Trigger: "<@%s> <3", Output: "<@%s> <3", Type: CommandTypeFullMatch, AdminOnly: true, OutputIsReply: true, RequiresMention: true})
    registerCommand(Command{Trigger: "!complain", Type: CommandTypePrefix, DMOnly: true, Function: redirectComplaint})



    fmt.Println("bot running. selphyWoo")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

/*func undelete(s *discordgo.Session, m *discordgo.MessageDelete) {
    channel, _ := s.State.Channel(m.ChannelID)
    message, _ := s.State.Message(m.ChannelID, m.ID)
    log.Println(fmt.Sprintf("Someone deleted a message in %s: “%s”", channel.Name, messageToString(message)))
}*/

