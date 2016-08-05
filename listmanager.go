package checkedbot

import (
  "net/http"
  "fmt"
  "encoding/json"
  "strings"
  "strconv"
  "regexp"
  "io/ioutil"
)

type Checklist struct {
  Id int
  Name string
  Url string
  Created_at string
}

type Listversion struct {
  Id int
}

type Liststep struct {
  Id int
  Position int
  Description string
}

type Execution struct {
  List Checklist
  Currentstep int
  Steps []Liststep
}

func getListDetails(listId int) (l Checklist, err error) {
  client := &http.Client {}
  req, _ := http.NewRequest("GET", "http://localhost:3004/checklists/"+strconv.Itoa(listId)+"?key=awesome", nil)
  req.Header.Set("Accept", "application/json")
  rsp, err := client.Do(req)
  if err != nil {
    fmt.Printf("Query failed: %v\n", err)
    return Checklist{}, err
  }
  defer rsp.Body.Close()
  body, _ := ioutil.ReadAll(rsp.Body)
  var c Checklist
  err = json.Unmarshal(body, &c)
  if err != nil {
    fmt.Printf("Parse failed: %v\n", err)
    return Checklist{}, err
  }
  return c, nil
}

func getListSteps(listId int) (steps []Liststep, err error) {
  client := &http.Client {}
  req, _ := http.NewRequest("GET", "http://localhost:3004/checklists/"+strconv.Itoa(listId)+"/versions?key=awesome", nil)
  req.Header.Set("Accept", "application/json")
  rsp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer rsp.Body.Close()
  body, _ := ioutil.ReadAll(rsp.Body)
  var v []Listversion
  err = json.Unmarshal(body, &v)
  if err != nil {
    return nil, err
  }
  latest := 0
  for _, version := range v {
    if version.Id > latest {
      latest = version.Id
    }
  }
  req, _ = http.NewRequest("GET", "http://localhost:3004/checklists/"+strconv.Itoa(listId)+"/versions/"+strconv.Itoa(latest)+"/steps?key=awesome", nil)
  req.Header.Set("Accept", "application/json")
  rsp, err = client.Do(req)
  if err != nil {
    return nil, err
  }
  defer rsp.Body.Close()
  body, _ = ioutil.ReadAll(rsp.Body)
  var s []Liststep
  err = json.Unmarshal(body, &s)
  if err != nil {
    return nil, err
  }
  fmt.Printf("Body: %v\n", string(body))
  fmt.Printf("Stuff I parsed: %#v\n", s)
  return s, nil
}


func (bot *Checkedbot) handleListRequest(chnl string) {
  client := &http.Client {}
  req, _ := http.NewRequest("GET", "http://localhost:3004/checklists/?key=awesome", nil)
  req.Header.Set("Accept", "application/json")
  rsp, err := client.Do(req)
  if err != nil {
    bot.simpleSay("I'm awfully sorry, but something seems to have gone awry.", chnl)
    fmt.Printf("Query failed: %v\n", err)
    return
  }
  defer rsp.Body.Close()
  body, _ := ioutil.ReadAll(rsp.Body)
  var c []Checklist
  err = json.Unmarshal(body, &c)
  if err != nil {
    bot.simpleSay("I'm awfully sorry, but something seems to have gone awry.", chnl)
    fmt.Printf("Parse failed: %v\n", err)
    return
  }
  fmt.Printf("Body: %v\n", string(body))
  fmt.Printf("Stuff I parsed: %#v\n", c)
  msg := "Here's your lists:\n\n"
  for _, list := range c {
    msg += "* ["+strconv.Itoa(list.Id)+"] "+list.Name+"\n"
  }
  bot.simpleSay(msg, chnl)
}

func (bot *Checkedbot) handleStartRequest(q string, chnl string) {
  if bot.current.List.Id != -1 {
    bot.simpleSay("Awfully sorry, Sir, but we are already executing list "+ strconv.Itoa(bot.current.List.Id) + " right now.\nMay I suggest finishing or aborting that first?", chnl)
    return
  }
  words := strings.Fields(q)
  var idx int
  for i, word := range words {
    matched, _ :=  regexp.MatchString("(i?)list", word)
    if matched {
      idx, _ = strconv.Atoi(words[i+1])
      break
    }
  }
  e := Execution{}
  e.List.Id = idx
  list, err := getListDetails(idx)
  if err != nil {
    bot.simpleSay("I'm awfully sorry, but something seems to have gone awry.", chnl)
    fmt.Printf("Getting details failed: %v\n", err)
    return
  }
  e.List = list
  steps, err := getListSteps(idx)
  if err != nil {
    bot.simpleSay("I'm awfully sorry, but something seems to have gone awry.", chnl)
    fmt.Printf("Getting details failed: %v\n", err)
    return
  }
  e.Steps = steps
  e.Currentstep = 0
  bot.current = e
  msg := "Certainly, Sir!\nTo *"+e.List.Name+"*, we need to follow these steps:\n\n"
  for i, step := range e.Steps {
    if e.Currentstep == i {
      msg += "  *"+strconv.Itoa(i+1)+". "+step.Description+"*\n"
    } else {
      msg += "  "+strconv.Itoa(i+1)+". "+step.Description+"\n"
    }
  }
  bot.simpleSay(msg, chnl)
}

func (bot *Checkedbot) handleNextStepRequest(chnl string) {
  if bot.current.List.Id == -1 {
    bot.simpleSay("Begging your pardon, Sir, but we seem to not currently be running through a list. Would Sir care to start one?", chnl)
    return
  }
  bot.current.Currentstep += 1
  if bot.current.Currentstep < len(bot.current.Steps) {
    msg := "Right, so the next step will be:\n\n  "+strconv.Itoa(bot.current.Currentstep+1)+". "+bot.current.Steps[bot.current.Currentstep].Description+"\n"
    bot.simpleSay(msg, chnl)
  } else {
    e := Execution{}
    e.Currentstep = -1
    e.List.Id = -1
    bot.current = e
    bot.simpleSay("And we're done! Glad to have been of service.", chnl)
  }
}

func (bot *Checkedbot) handleDetailsRequest(q string, chnl string) {
  words := strings.Fields(q)
  var idx int
  for i, word := range words {
    matched, _ :=  regexp.MatchString("(i?)list", word)
    if matched {
      idx, _ = strconv.Atoi(words[i+1])
      break
    }
  }
  fmt.Println("List ID I heard: " + strconv.Itoa(idx))
  list, err := getListDetails(idx)
  msg := "Here's the details for list "+strconv.Itoa(list.Id)+":\n\n"
  msg += "* ["+strconv.Itoa(list.Id)+"] "+list.Name+", created at "+list.Created_at+"\n\nSteps:\n"
  steps, err := getListSteps(list.Id);
  if err != nil {
    bot.simpleSay("I'm awfully sorry, but something seems to have gone awry.", chnl)
    fmt.Printf("Getting steps failed: %v\n", err)
    return
  }
  for _, step := range steps {
    msg += "     "+strconv.Itoa(step.Position+1)+". "+step.Description+"\n"
  }
  bot.simpleSay(msg, chnl)
}