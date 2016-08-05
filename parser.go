package checkedbot

import (
    "github.com/nlopes/slack"
    "fmt"
    "regexp"
    "strings"
)

func (bot *Checkedbot) parseQuery(q string) (t string, matched bool) {
  t, matched = "", false
  matched, _ = regexp.MatchString("(?i).*(show|what are).*(lists|checklists).*", q)
  if matched {
    return "show", matched
  }
  matched, _ = regexp.MatchString("(?i)(Hi|hi|Hello|hello|Hey|hey).*@"+bot.userid, q)
  if matched {
    return "greet", matched
  }
  matched, _ = regexp.MatchString("(i?).*(start|run|execute).*(list|checklist)", q)
  if matched {
    return "start", matched
  }
  matched, _ = regexp.MatchString("(i?).*(cancel|end|abort).*(list|checklist)", q)
  if matched {
    return "abort", matched
  }
  matched, _ = regexp.MatchString("(i?).*(what is|[tT]ell me about).*(list) [0-9]", q)
  if matched {
    return "details", matched
  }
  matched, _ = regexp.MatchString("(i?).*(next|next step|done).*", q)
  if matched {
    return "next", matched
  }
  return "", matched
}


func (bot *Checkedbot) handle(ev *slack.MessageEvent) {
  if ev.SubType == "bot_message" || ev.User == bot.userid {
    return
  }
  var msg string
  if ev.SubType == "message_changed" {
    fmt.Println("Message changed")
    fmt.Printf("New Message: %s\n", ev.SubMessage.Text)
    msg = "I saw that!"
  } else {
    fmt.Printf("Message: %s\n", ev.Text)
    if ! strings.Contains(ev.Text, "@"+bot.userid) {
      return
    }
    query, matched := bot.parseQuery(ev.Text)
    if matched {
      if query == "show" {
        bot.simpleSay("Certainly, Sir! Let me look those up for you.", ev.Channel)
        bot.handleListRequest(ev.Channel)
        return
      } else if query == "greet" {
        msg = "And a good day to you, too, Sir!"
      } else if query == "start" {
        bot.handleStartRequest(ev.Text, ev.Channel)
      } else if query == "abort" {
        //bot.handleAbortRequest(ev.Channel)
      } else if query == "details" {
        bot.handleDetailsRequest(ev.Text, ev.Channel)
      } else if query == "next" {
        bot.handleNextStepRequest(ev.Channel)
      }
    } else {
      msg = "I know you are, what what am I?"
    }
  }
  bot.simpleSay(msg, ev.Channel)
}