// Who's On It. "Email tracking service"
package whosonit

import (
    "appengine"
    "appengine/datastore"
//    "appengine/mail"
//    "appengine/user"
//    "bytes"
    "fmt"
    "http"
    "time"
    "strconv"
    )

type MailEvent struct {
  Sender string
  Owner string
  RecieptDate datastore.Time
  OwnerDate datastore.Time
  ClosedDate datastore.Time
  Subject string
  Body string
}

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/accept", accept_event)
  http.HandleFunc("/close", close_event)
  http.HandleFunc("/show", show_event)
  http.HandleFunc("/_ah/mail/", new_mail_event)
  http.HandleFunc("/form", show_form)
  http.HandleFunc("/test_form", handle_test_form)
}

func show_form(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, test_form)
}

func handle_test_form(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  event := MailEvent{
    Body: r.FormValue("Body"),
    Sender: r.FormValue("Sender"),
    RecieptDate: datastore.SecondsToTime(time.Seconds()),
  }
  _, err := datastore.Put(c, datastore.NewIncompleteKey("Events"), &event)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/", http.StatusFound)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  query := datastore.NewQuery("Events").Filter("ClosedDate <",
              time.Seconds()).Order("ClosedDate")
  events := make([]MailEvent,0, 100)
  if _, err := query.GetAll(c, &events); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  if err := rootTemplate.Execute(w, events); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}

func find_event(c appengine.Context, sender string, timestring string, t *[]MailEvent) (bool, string) {
  timestamp, _ := strconv.Atoi64(timestring)
  query := datastore.NewQuery("Events").Filter("Sender =", sender).Filter("RecieptDate =", timestamp).Limit(1)
  if count, _ := query.Count(c); count < 1 {
    err := fmt.Sprintf("No results found for Sender = %s and Date = %s", sender, timestamp)
  return false, err
  } else if count > 1 {
    err := fmt.Sprintf("Too many results found for Sender = %s and Date = %s", sender, timestamp)
    return false, err
  }
  _, err := query.GetAll(c, t)
  if err != nil {
    return false, err.String()
  }
  return true, ""
}

func show_event(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  event := make([]MailEvent, 0, 1)
  status, err_string := find_event(c, r.FormValue("sender"),
                                   r.FormValue("date"), &event)
  if ! status {
    fmt.Fprintf(w, err_string)
    return
  }
  if err := showTemplate.Execute(w, event); err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}

func find_single_event(c appengine.Context, sender string, timestring string) (MailEvent, *datastore.Key, string) {
  timestamp, _ := strconv.Atoi64(timestring)
  event := MailEvent{}
  key := datastore.NewIncompleteKey("key")
  query := datastore.NewQuery("Events").Filter("Sender =", sender).Filter("RecieptDate =", timestamp).Limit(1).KeysOnly()
  if count, _ := query.Count(c); count < 1 {
    err := fmt.Sprintf("No results found for Sender = %s and Date = %s", sender, timestamp)
    return event, key, err
  } else if count > 1 {
    err := fmt.Sprintf("Too many results found for Sender = %s and Date = %s", sender, timestamp)
    return event, key, err
  }
  keys, err := query.GetAll(c, nil)
  if err != nil {
    return event, key, err.String()
  }

  err = datastore.Get(c, keys[0], &event)
  if err != nil {
    return event, key, err.String()
  }
  return event, keys[0], ""
}

func accept_event(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  event, key, err_string := find_single_event(c, r.FormValue("sender"), 
                                              r.FormValue("date"))

  if len(err_string) > 0 {
    fmt.Fprintf(w, err_string)
    return
  }
  if event.OwnerDate > 0 {
    fmt.Fprintf(w, "This message is already owned by %s.", event.Owner)
    return
  }
  event.OwnerDate = datastore.SecondsToTime(time.Seconds())
  event.Owner = "Someone"
  _, err := datastore.Put(c, key, &event)
  if err != nil {
    fmt.Fprintf(w, err.String())
    return
  } else {
    target_url := fmt.Sprintf("/show?sender=%s&date=%d", event.Sender,
                              event.RecieptDate)
    http.Redirect(w, r, target_url, http.StatusFound)
  }
}

func close_event(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  event, key, err_string := find_single_event(c, r.FormValue("sender"), 
                                              r.FormValue("date"))

  if len(err_string) > 0 {
    fmt.Fprintf(w, err_string)
    return
  }
  // TODO: Check to make sure the closer matches the owner
  if event.ClosedDate > 0 {
    fmt.Fprintf(w, "This message is already closed.")
    return
  }
  event.ClosedDate = datastore.SecondsToTime(time.Seconds())
  _, err := datastore.Put(c, key, &event)
  if err != nil {
    fmt.Fprintf(w, err.String())
    return
  } else {
    target_url := fmt.Sprintf("/show?sender=%s&date=%d", event.Sender,
                              event.RecieptDate)
    http.Redirect(w, r, target_url, http.StatusFound)
  }
}

func new_mail_event(w http.ResponseWriter, r *http.Request) {
//  c := appengine.NewContext(r)
//  defer r.Body.Close()
//  var b bytes.Buffer
//  if _, err := b.ReadFrom(r.Body); err != nil {
//    c.Errorf("Error reading body: %v", err)
//    return
//  }
//  //c.Infof("Received mail: %v", b)
}

