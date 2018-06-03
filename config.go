package main

import (
    "os"
    "encoding/json"
)

type Embed struct {
    Message string
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
    GeneralChannel string
    SendWelcomeDM bool
    RequireAccept bool
    ComplaintReceivedMessage string
    ModChannel string
    WelcomeEmbed Embed
    RoleCommands map[string]string
}


func readConfig() Config {
    file, _ := os.Open("config2.json")
    conf := Config{}
    json.NewDecoder(file).Decode(&conf)
    file.Close()
    return conf
}

