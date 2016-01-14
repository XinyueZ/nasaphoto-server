# nasaphoto-server(APOD+)

Give list of photo of APOD(Astronomy Picture of the Day) from NASA after 1998(inc.)
by calling different date-time scopes.

# Reason for this API
The default APOD-API doesn't provide list.

# Example on GAE
http://nasa-photo-dev.appspot.com
http://orbital-stage-648.appspot.com

# Specification

  API| Method|Comment
--------|--------- |---------
  [/list](#1-by-different-date)|POST  | Get list of photos with different single date.
  [/month_list](#2-one-month)|POST  |Get list of photos of a month.
  [/last_three_list](#3-last-3-month-including-today)|POST | Get list of photos of last 3 days including today.


# cron Task

An auto called task(cron) for archive history in Firebase.

Default period is 15 minutes, checkout [cron.yaml](https://github.com/XinyueZ/nasaphoto-server/blob/master/cron.yaml).

API| Method|Comment
--------|--------- |---------
/whistory|GET  | Checkout [write_histroy.go](https://github.com/XinyueZ/nasaphoto-server/blob/master/write_history.go) for building history from 1998-1-1.

# Config file

After checkout it's still lack of a config.go
that gives some base value of whole project(inc. key to NASA API, Firebase's Auth).

```go
package index

const (
	KEY      = "vvvvvvvvvv"
	HOST     = "https://api.nasa.gov"
	API_APOD = "%s/planetary/apod?date=%s&hd=true&concept_tags=false&api_key=" + KEY

	FIRE_URL  = "https://xxxxxx.firebaseio.com/"
	FIRE_AUTH = "yyyyyyy"
)
```

Value|Comment
--------|---------
KEY|[Key for NASA API](https://api.nasa.gov/index.html#apply-for-an-api-key).
HOST|Host of NASA API, must be https://api.nasa.gov
API_APOD|[API location of  APOD](https://api.nasa.gov/api.html#apod).
FIRE_URL | The location of Firebase to save history.
AUTH | Auth for Firebase.

# Response

```json
{
  "status": 200,
  "reqId": "sadfadsf-a345345-as456456-353456adsfa",
  "result": [
    {
      "reqId": "2016-1-3",
      "title": "A Starry Night of Iceland",
      "description": "On some nights, the sky is the best show in town. On this night, the sky was not only the best show in town.....",
      "date": "2016-1-3",
      "urls": {
        "normal": "http://apod.nasa.gov/apod/image/1601/aurora_vetter_1080.jpg",
        "hd": "http://apod.nasa.gov/apod/image/1601/aurora_vetter_2000.jpg"
      }
    },
    {
      "reqId": "2016-1-4",
      "title": "Earthset from the Lunar Reconnaissance Orbiter",
      "description": "On the Moon, the Earth never rises -- or sets.  If you were to sit on the surface of the Moon, you would see the Earth just hang in the sky. This is because the Moon always keeps the same side toward the Earth. Curiously... ",
      "date": "2016-1-4",
      "urls": {
        "normal": "http://apod.nasa.gov/apod/image/1601/Earthrise_LRO_960.jpg",
        "hd": "http://apod.nasa.gov/apod/image/1601/Earthrise_LRO_5634.jpg"
      }
    },

  ......
  ]
}

```

# Request

1. Give list of different date.

2. Give list of one month.

3. Give list of last three days.


# Example

## 1. By different date

Request body for /list:

```json
{
    "reqId" : "sadfadsf-a345345-as456456-353456adsfa",
    "dates" : [
        "2016-1-3",
        "2016-1-4",
        "2016-1-5",
        "2016-1-6",
        "2016-1-7",
        "2016-1-8"
    ],
    "timeZone" : "CET"
}
```

## 2. One month.

Request body for /month_list:

```json
{
    "reqId" : "sadfadsf-a345345-as456456-353456adsfa",
    "year" : 2015,
    "month" : 1,
    "timeZone" : "CET"
}
```

## 3. Last 3 month including today.

Request body for /last_three_list:

```json
{
    "reqId" : "sadfadsf-a345345-as456456-353456adsfa",
    "timeZone" : "CET"
}
```


# Status and Error

Code|  Comment
--------| ---------
300|No result to request.   
301|POST with invalid body.
302|POST body has invalid content.
303|Firebase push error.
304|Firebase update command error.
500|Critical Error.


```json
//Error
{
    "status" : 500,
    "message" : "error for some reasons."
}
//Success
{
    "status" : 200,
    .....
}
```


# Liscense
======
```json
			Copyright Xinyue Zhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
