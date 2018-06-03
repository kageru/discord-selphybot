package main

type CommandType int

const (
    CommandTypePrefix      CommandType = 0
    CommandTypeFullMatch   CommandType = 1
    CommandTypeRegex       CommandType = 2
)

type Command struct {
    Input string
    Output string
    Type CommandType
    OutputIsReply bool
    DeleteInput bool
    DMOnly bool
}


