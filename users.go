package checkedbot

import "errors"
import "github.com/nlopes/slack"


func (bot *Checkedbot) FindUser(username string) (usr *slack.User, err error) {
  api := bot.api
  users, err := api.GetUsers()
  if err != nil {
    return nil, err
  }
  for _, user := range users {
    if user.Name == username {
      return &user, nil
    }
  }
  return nil, errors.New("User not found")
}

func (bot *Checkedbot) OpenChannel(username string) (channel_id string, err error) {
  api := bot.api
  user, err := bot.FindUser(username)
  if err != nil {
    return "", err
  }
  id := user.ID
  _, _, channel_id, err = api.OpenIMChannel(id)
  if err != nil {
    return "", err
  }
  return channel_id, nil
}