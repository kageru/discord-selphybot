package main

import (
    "github.com/bwmarrin/discordgo"
    "fmt"
)

func getWelcomeEmbed() *discordgo.MessageEmbed {
    return &discordgo.MessageEmbed {
        Author:      &discordgo.MessageEmbedAuthor{},
        Color:       0xffb90f,
        Description: config.WelcomeEmbed.Message,
        Fields: []*discordgo.MessageEmbedField {
            &discordgo.MessageEmbedField {
                Name:   config.WelcomeEmbed.QuestionsTitle,
                Value:  config.WelcomeEmbed.QuestionsText,
                Inline: true,
            },
            &discordgo.MessageEmbedField {
                Name:   config.WelcomeEmbed.BugsTitle,
                Value:  fmt.Sprintf(config.WelcomeEmbed.BugsText, config.AdminID),
                Inline: true,
            },
        },
        Thumbnail: &discordgo.MessageEmbedThumbnail{
            URL: config.WelcomeEmbed.Image,
        },
    }
}
