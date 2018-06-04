package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "strings"
    "log"
    "regexp"
)

type CommandType int

const (
    CommandTypePrefix      CommandType = 0
    CommandTypeFullMatch   CommandType = 1
    CommandTypeRegex       CommandType = 2
    CommandTypeContains    CommandType = 3
)

type Command struct {
    Trigger string                          // must be specified
    Output string                           // no output if unspecified
    OutputEmbed *discordgo.MessageEmbed     // no embed output if unspecified
    Type CommandType                        // defaults to Prefix
    OutputIsReply bool                      // defaults to false
    RequiresMention bool                    // defaults to false
    DeleteInput bool                        // defaults to false
    DMOnly bool                             // defaults to false
    AdminOnly bool                          // defaults to false
    // for custom commands that go beyond prints and deletions
    Function func(*discordgo.Session, *discordgo.MessageCreate)
}


func registerCommand(command Command) {
    if command.Trigger == "" {
        fmt.Println("Cannot register a command with no trigger. Skipping.")
        return
    }
    commands = append(commands, command)
}

func evaluateMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        log.Printf("<Self> %s", m.Content)
        return
    }
    fmt.Println(m.Content)
    for _, command := range commands {
        if command.RequiresMention {
            command.Trigger = fmt.Sprintf(command.Trigger, s.State.User.ID)
        }
        switch command.Type {
        case CommandTypePrefix:
            if strings.HasPrefix(m.Content, command.Trigger) {
                executeCommand(s, m, command)
                return
            }
        case CommandTypeFullMatch:
            if m.Content == command.Trigger {
                executeCommand(s, m, command)
                return
            }
        case CommandTypeRegex:
            match, _ := regexp.MatchString(command.Trigger, m.Content)
            if match {
                executeCommand(s, m, command)
                return
            }
        case CommandTypeContains:
            if strings.Contains(m.Content, command.Trigger) {
                executeCommand(s, m, command)
                return
            }
        }
    }
}

func executeCommand(session *discordgo.Session, message *discordgo.MessageCreate, command Command) {
    if (!command.DMOnly || (getChannel(session.State, message.ChannelID).Type == discordgo.ChannelTypeDM)) &&
        (!command.AdminOnly || (message.Author.ID == config.AdminID)) {
        log.Printf("Executed command %s triggered by user %s", command.Trigger, userToString(message.Author))
        fmt.Printf("Executed command %s triggered by user %s", command.Trigger, userToString(message.Author))
        if command.Function == nil {
            // simple reply
            if command.OutputEmbed == nil {
                messageContent := generateReply(message, command)
                session.ChannelMessageSend(message.ChannelID, messageContent)
            } else {
                session.ChannelMessageSendEmbed(message.ChannelID, command.OutputEmbed)
            }
            if command.DeleteInput {
                session.ChannelMessageDelete(message.ChannelID, message.ID)
            }
        } else {
            // execute custom function
            command.Function(session, message)
        }
    } else {
        log.Printf("Denied command %s to user %s.", command.Trigger, userToString(message.Author))
    }
}

func generateReply(message *discordgo.MessageCreate, command Command) string {
    output := command.Output
    if command.OutputIsReply {
        output = fmt.Sprintf(output, message.Author.ID)
    }
    return output
}

func redirectComplaint(s *discordgo.Session, m *discordgo.MessageCreate) {
    embed := &discordgo.MessageEmbed {
            Author:      &discordgo.MessageEmbedAuthor{},
            Color:       0xbb0000,
            Description: m.Content,
    }
    s.ChannelMessageSendEmbed(config.ModChannel, embed)
}

