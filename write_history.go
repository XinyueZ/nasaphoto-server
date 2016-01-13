package index

import (
	"appengine"
	"firego"

	"fmt"
	"net/http"
	"time"
)

const (
	START_HISTORY = "1998-1-1"
)

func init() {
	http.HandleFunc("/whistory", handleWriteHistory)
}

//saveAddHistoryTime
//Save the last time that saved history.
func saveAddHistoryTime(w http.ResponseWriter, r *http.Request, t map[string]interface{}) {
	f := firego.NewGAE(appengine.NewContext(r), FIRE_URL)
	f.Auth(FIRE_AUTH)
	if e := f.Update(t); e == nil {
		if _, e := f.Push(nil); e == nil {
			status(w, fmt.Sprintf("updated history: %v", t["lastSave"]), 200)
		} else {
			s := fmt.Sprintf("Firebase push error: %v", e)
			status(w, s, 303)
		}
	} else {
		s := fmt.Sprintf("Firebase update command error: %v", e)
		status(w, s, 304)
	}
}

//getAddHistoryTime
//Get last time to write to database.
//Start with 1998-1-1.
func getAddHistoryTime(w http.ResponseWriter, r *http.Request) (t map[string]interface{}, v string) {
	fLastSave := firego.NewGAE(appengine.NewContext(r), FIRE_URL+"lastSave")
	fLastSave.Auth(FIRE_AUTH)
	var lastValue string
	if err := fLastSave.Value(&lastValue); err == nil {
		if lastValue == "" {
			//Insert new time
			t = map[string]interface{}{
				"lastSave": START_HISTORY,
			}
		} else {
			//Update last time
			lastTime, _ := time.Parse("2006-1-2", lastValue)
			t = map[string]interface{}{
				"lastSave": lastTime.AddDate(0, 0, 1).Format("2006-1-2"),
			}
		}
		v = fmt.Sprintf("%v", t["lastSave"])
	}
	return
}

//handleWriteHistory
//Add one history to our database.
func handleWriteHistory(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, fmt.Sprintf("%v", err), 500)
		}
	}()

	full, value := getAddHistoryTime(w, r)

	request := Request{"0", []string{value}, "CET"}
	plist := buildResult(r, &request)
	if len(plist) > 0 {
		f := firego.NewGAE(appengine.NewContext(r), FIRE_URL+"history")
		f.Auth(FIRE_AUTH)
		photo := plist[0]
		if _, e := f.Push(photo); e == nil {
			saveAddHistoryTime(w, r, full)
		}
	}
}
