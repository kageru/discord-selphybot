package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "strings"
    "log"
    "regexp"
)

func onJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
    if !member.User.Bot && config.RequireAccept {
        s.GuildMemberRoleAdd(config.ServerID, member.User.ID, config.LockedRoleID)
    }
    if !member.User.Bot && config.SendWelcomeDM {
        dm, err := s.UserChannelCreate(member.User.ID)
        if err != nil {
            log.Println(fmt.Sprintf("Error creating DM with %s", userToString(member.User), err))
        } else {
            embed := getWelcomeEmbed()
            _, err = s.ChannelMessageSendEmbed(dm.ID, embed)
            if err != nil {
                log.Println(fmt.Sprintf("Error sending DM to %s", userToString(member.User), err))
            }
        }
        if err != nil {
            // if any of the preceding operations produced an error
            log.Printf("Sending welcome @mention at %s", userToString(member.User))
            s.ChannelMessageSend(config.GeneralChannel, fmt.Sprintf("Wilkommen <@%s>. Bitte aktiviere vorübergehend DMs für diesen Server und sende eine Nachricht mit !welcome an mich.", member.User.ID))
        }
    }
    log.Printf("User joined: %s", userToString(member.User))
}

func onDM(s *discordgo.Session, m *discordgo.MessageCreate) {
    log.Printf("Received DM from %s with content: “%s”", userToString(m.Author), m.Content)
    fmt.Sprintf("Received DM from %s with content: “%s”", userToString(m.Author), m.Content)
    Member, _ := s.GuildMember(config.ServerID, m.Author.ID)
    dm, _ := s.UserChannelCreate(Member.User.ID)
    if m.Content == "!welcome" {
        s.ChannelMessageSendEmbed(dm.ID, getWelcomeEmbed())
        return
    }
    if strings.HasPrefix(m.Content, "!complain") {
        redirectComplaint(s, m)
        s.ChannelMessageSend(dm.ID, config.ComplaintReceivedMessage)
        return
    }
    for comm, role := range config.RoleCommands {
        if m.Content == comm {
            for _, irole := range config.RoleCommands {
                for _, mrole := range Member.Roles {
                    if irole == mrole {
                        s.ChannelMessageSend(dm.ID, "Baka, du kannst nur eine der Rollen haben.")
                        log.Printf("Denied Role %s to %s. User already has %s", roleName(s.State, irole), userToString(m.Author), roleName(s.State, irole))
                        return
                    }
                }
            }
            log.Printf("Giving Role %s to %s", roleName(s.State, role), userToString(m.Author))
            s.ChannelMessageSend(dm.ID, "Haaai, Ryoukai desu~")
            s.GuildMemberRoleAdd(config.ServerID, m.Author.ID, role)
        }
    }
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        log.Printf("<Self> %s", m.Content)
        return
    }

    if getChannel(s.State, m.ChannelID).Type == discordgo.ChannelTypeDM {
        onDM(s, m)
    }

    if m.Author.ID == config.AdminID {
        //replyGodmode(s, m)
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
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> Oboe!", m.Author.ID))
        s.ChannelMessageDelete(m.ChannelID, m.ID)
        channel := getChannel(s.State, m.ChannelID)
        log.Printf("Deleted message by %s in %s. Content: “%s”", userToString(m.Author), channelToString(channel), m.Content)
        return
    }

    // As per our majesty’s command:
    if m.Content == "\\o" {
        s.ChannelMessageSend(m.ChannelID, "o/")
        log.Printf("o/ at %s", userToString(m.Author))
    } else if m.Content == "o/" {
        s.ChannelMessageSend(m.ChannelID, "\\o")
        log.Printf("\\o at %s", userToString(m.Author))
    } else if m.Content == "ayy" {
        s.ChannelMessageSend(m.ChannelID, "lmao")
        log.Printf("ayy lmao at %s", userToString(m.Author))
    }
}


// Admin stuff down here. This is very server-specific

func adminshit (s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Content == fmt.Sprintf("<@%s> <3", s.State.User.ID) {
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> <3", config.AdminID))
    }
/*
    if m.Content == "print_rules()" {
        channel := getChannel(s.State, m.ChannelID)
        log.Printf("print_rules() triggered by %s in %s.", userToString(m.Author), channelToString(channel))
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
*/
}


