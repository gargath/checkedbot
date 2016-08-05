package checkedbot

import (
  "fmt"
  "github.com/nlopes/slack"
  "os"
  "io/ioutil"
  "strings"
)

type Checkedbot struct {
  api *slack.Client
  rtm *slack.RTM
  userid string
  current Execution
}

func (bot *Checkedbot) simpleSay(msg string, chnl string) {
  bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(msg, chnl))
}

func initialize() *Checkedbot {
  println("Checkedbot starting")  

  i := Execution{}
  i.List.Id = -1
  file, e := ioutil.ReadFile("./config")
  if e != nil {
    fmt.Println("Failed to read config file.")
    os.Exit(1)
  }
  key := strings.Split(string(file),"\n")
  api := slack.New(key[0])
  api.SetDebug(false)
  b := &Checkedbot{}
  b.current = i
  b.api = api
  b.cleanup()
  users, err := api.GetUsers()
  if err != nil {
    fmt.Println("ERROR: Failed to identify own ID: ", err)
    os.Exit(1)
  }
  for _, user := range users {
    if user.Name == "checkedbot" {
      b.userid = user.ID
    }
  }
  if b.userid == "" {
    fmt.Println("ERROR: Failed to identify own ID.")
    os.Exit(1)
  }
  return b
}

func Start() {
  bot := initialize()
  
  err := bot.present(true)
  if err != nil {
    fmt.Println("WARN: Failed to udpate presence")
  }

  rtm := bot.api.NewRTM()
  bot.rtm = rtm
  go rtm.ManageConnection()
  
  channels, err := bot.api.GetChannels(true)
    if err != nil {
      fmt.Println("Failed to get channel list!")
    }
    for _,channel := range channels {
      if channel.IsMember {
        fmt.Printf("Announcing in channel %s (%s)\n", channel.Name, channel.ID)
        bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("Checkedbot online and ready", channel.ID))
      }
    }
  
  for  {
    select {
      case msg := <-rtm.IncomingEvents:
        switch ev := msg.Data.(type) {
          case *slack.MessageEvent:
            bot.handle(ev)
          default:
            fmt.Printf("=====Event: %+v\n", ev)
        }
    }
  }

}

