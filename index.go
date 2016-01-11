package index

import (
	"appengine"
	"appengine/urlfetch"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	ReqId    string   `json:"reqId"`
	Dates    []string `json:"dates"`
	TimeZone string   `json:"timeZone"`
}

func (request *Request) newRequest() (r *Request) {
	r = request
	newDates := make([]string, 0)

	tz := request.TimeZone
	location, _ := time.LoadLocation(tz)
	today := time.Now().In(location)

	length := len(request.Dates)
	for i := 0; i < length; i++ {
		t, err := time.Parse("2006-1-2", request.Dates[i])

		//Check for invalid dateformat.
		if err != nil || t.Year() < 1998 || t.Year() == 0 || int(t.Month()) == 0 || t.Day() <= 0 {
			continue
		}

		//Check for invalid "day".
		ss := strings.Split(request.Dates[i], "-")
		if day, err := strconv.Atoi(ss[2]); err == nil {
			if day <= 0 {
				continue
			}
		} else {
			continue
		}

		//Check for invalid in "this month".
		if t.Year() == today.Year() && t.Month() > today.Month() {
			continue
		} else {
			newDates = append(newDates, request.Dates[i])
		}
	}
	r.Dates = newDates
	return
}

type LastThreeRequest struct {
	ReqId    string `json:"reqId"`
	TimeZone string `json:"timeZone"`
}

func (request *LastThreeRequest) newRequest() (r *Request) {
	tz := request.TimeZone
	location, _ := time.LoadLocation(tz)

	today := time.Now().In(location)
	now := today.Format("2006-1-2")
	beforeYesterday := today.AddDate(0, 0, -2).Format("2006-1-2")
	yesterday := today.AddDate(0, 0, -1).Format("2006-1-2")
	r = &Request{request.ReqId, []string{beforeYesterday, yesterday, now}, request.TimeZone}
	return
}

type MonthRequest struct {
	ReqId    string `json:"reqId"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	TimeZone string `json:"timeZone"`
}

// daysIn returns the number of days in a month for a given year.
func (request *MonthRequest) daysIn() int {
	// This is equivalent to time.daysIn(m, year).
	mon := time.Month(request.Month + 1)
	return time.Date(request.Year, mon, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (request *MonthRequest) newRequest() (r *Request) {
	tz := request.TimeZone
	location, _ := time.LoadLocation(tz)
	today := time.Now().In(location)

	r = &Request{}
	r.ReqId = request.ReqId
	r.Dates = []string{}

	//Check for invalid year, zero objects.
	if request.Year < 1998 || request.Year == 0 || request.Month == 0 {
		return
	}

	//Check for invalid month in this year.
	if request.Year == today.Year() && request.Month > int(today.Month()) {
		return
	}

	daysIn := request.daysIn()
	mon := time.Month(request.Month)
	for day := 1; day <= daysIn; day++ {
		//Check for invalid day in "this month".
		if request.Year == today.Year() && request.Month == int(today.Month()) {
			if day > today.Day() {
				return
			}
		}
		dt := time.Date(request.Year, mon, day, 0, 0, 0, 0, time.UTC)
		r.Dates = append(r.Dates, dt.Format("2006-1-2"))
	}
	return
}

type Meta struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        string `json:"date"`
	Url         string `json:"url"`
	HDUrl       string `json:"hdurl"`
	MediaType   string `json:"media_type"`
}

type Photo struct {
	ReqId       string `json:"reqId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Urls        Urls   `json:"urls"`
}

func (photo *Photo) fromMeta(pMeta *Meta) {
	photo.Title = pMeta.Title
	photo.Description = pMeta.Explanation
	photo.Date = pMeta.Date
	photo.Urls = Urls{pMeta.Url, pMeta.HDUrl}

	photo.ReqId = pMeta.Date
}

type Urls struct {
	Normal string `json:"normal"`
	HD     string `json:"hd"`
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func init() {
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/last_three_list", handleLastThreeList)
	http.HandleFunc("/month_list", handleMonthList)
}

func status(w http.ResponseWriter, reqId string, status int) {
	if reqId == "" {
		reqId = "not provided"
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":%d, "reqId" : "%s"}`, status, reqId)
}

func response(w http.ResponseWriter, reqId string, photo []*Photo) {
	if reqId == "" {
		reqId = "not provided"
	}
	json, _ := json.Marshal(photo)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":200, "reqId" : "%s", "result" : %s}`, reqId, string(json))
}

//getPhoto
//Call API of NASA to getting photos.
func getPhoto(r *http.Request, pDate *string, chPhoto chan *Photo) {
	cxt := appengine.NewContext(r)
	url := fmt.Sprintf(API_APOD, HOST, *pDate)
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		httpClient := urlfetch.Client(cxt)
		r, err := httpClient.Do(req)
		if r != nil {
			defer r.Body.Close()
		}
		if err == nil {
			if bytes, err := ioutil.ReadAll(r.Body); err == nil {
				pMeta := new(Meta)
				json.Unmarshal(bytes, pMeta)
				photo := new(Photo)
				photo.fromMeta(pMeta)
				chPhoto <- photo
			} else {
				cxt.Errorf("getPhoto: %v", err)
				chPhoto <- nil
			}
		} else {
			cxt.Errorf("getPhoto: %v", err)
			chPhoto <- nil
		}
	} else {
		cxt.Errorf("getPhoto: %v", err)
		chPhoto <- nil
	}
}

//showList
func showList(w http.ResponseWriter, r *http.Request, p *Request) {
	dates := p.Dates
	length := len(dates)
	list := make([]*Photo, 0)
	ch := make(chan *Photo, length)
	for i := 0; i < length; i++ {
		dt := dates[i]
		go getPhoto(r, &dt, ch)
		list = append(list, <-ch)
	}
	if list != nil {
		response(w, p.ReqId, list)
	} else {
		status(w, p.ReqId, 500)
	}
}

//handleList
//Get list of photos of specified dates.
func handleList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, fmt.Sprintf("%v", err), 500)
		}
	}()

	req := Request{}
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		if e := json.Unmarshal(bytes, &req); e == nil {
			showList(w, r, req.newRequest())
		} else {
			s := fmt.Sprintf("%v", e)
			status(w, s, 500)
		}
	} else {
		s := fmt.Sprintf("%v", e)
		status(w, s, 500)
	}
}

//handleLastThreeList
//Get list for last three days including today.
func handleLastThreeList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, fmt.Sprintf("%v", err), 500)
		}
	}()

	ltr := LastThreeRequest{}
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		if e := json.Unmarshal(bytes, &ltr); e == nil {
			showList(w, r, ltr.newRequest())
		} else {
			s := fmt.Sprintf("%v", e)
			status(w, s, 500)
		}
	} else {
		s := fmt.Sprintf("%v", e)
		status(w, s, 500)
	}
}

//handleMonthList
//Get list for whole month.
func handleMonthList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, fmt.Sprintf("%v", err), 500)
		}
	}()

	monthRequest := MonthRequest{}
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		if e := json.Unmarshal(bytes, &monthRequest); e == nil {
			showList(w, r, monthRequest.newRequest())
		} else {
			s := fmt.Sprintf("%v", e)
			status(w, s, 500)
		}
	} else {
		s := fmt.Sprintf("%v", e)
		status(w, s, 500)
	}
}
