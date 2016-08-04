package checkedbot

import (
  "fmt"
  "github.com/nlopes/slack"
  "os"
)

type Checkedbot struct {
  api *slack.Client
  rtm *slack.RTM
  userid string
}

func initialize() *Checkedbot {
  println("Checkedbot starting")  

  api := slack.New("xoxb-66195202069-gePVrKqmbSWDVyl29Ry4FS1q")
  api.SetDebug(false)
  b := &Checkedbot{}
  b.api = api
  b.cleanup()
  users, err := api.GetUsers()
  if err != nil {
    fmt.Println("ERROR: Failed to identify own ID.")
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
  for  {
    select {
      case msg := <-rtm.IncomingEvents:
        switch ev := msg.Data.(type) {
          case *slack.MessageEvent:
            if ev.SubType == "bot_message" || ev.User == bot.userid {
              break
            }
            var msg string
            if ev.SubType == "message_changed" {
              fmt.Println("Message changed")
              fmt.Printf("New Message: %s\n", ev.SubMessage.Text)
              msg = "I saw that!"
            } else {
              fmt.Printf("Message: %s\n", ev.Text)
              msg = "I know you are, what what am I?"
            }
            rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
          default:
            fmt.Printf("=====Event: %+v\n", ev)
        }
    }
  }

}

