package index

import (
	"appengine"
	"appengine/urlfetch"

	"encoding/base64"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
)

type Request struct {
	ReqId string   `json:"reqId"`
	Dates []string `json:"dates"`
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

	bys := []byte(pMeta.Url)
	photo.ReqId = base64.StdEncoding.EncodeToString(bys)
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

func handleList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			status(w, fmt.Sprintf("%v", err), 500)
		}
	}()

	req := Request{}
	if bytes, e := ioutil.ReadAll(r.Body); e == nil {
		if e := json.Unmarshal(bytes, &req); e == nil {
			// cxt := appengine.NewContext(r)
			// cxt.Infof("dates:%v", len(req.Dates))
			showList(w, r, &req)
		} else {
			s := fmt.Sprintf("%v", e)
			status(w, s, 500)
		}
	} else {
		s := fmt.Sprintf("%v", e)
		status(w, s, 500)
	}
}
