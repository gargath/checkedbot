package checkedbot

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
)

func (bot *Checkedbot) present(present bool) (err error) {
  api := bot.api
  var status string
  if present {
    status = "auto"
  } else {
    status = "away"
  }
  err = api.SetUserPresence(status)
  return err
}

func (bot *Checkedbot) cleanup() {
  api := bot.api
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
    <- c
    fmt.Println("Checkedbot signing off.")
    channels, err := api.GetChannels(true)
    if err != nil {
      fmt.Println("Failed to get channel list!")
    }
    for _,channel := range channels {
      if channel.IsMember {
        fmt.Printf("Signing off from channel %s (%s)\n", channel.Name, channel.ID)
        bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("Checkedbot signing off", channel.ID))
      }
    }
    bot.present(false)
    fmt.Println("Signoff finished")
    os.Exit(0)
  }()
}