# nasaphoto-server

Give list of photo of APOD from NASA.

Example host on GAE: http://orbital-stage-648.appspot.com

API:

  API| Method|Comment
--------|--------- |---------
   [/list](#1-by-different-date)|POST  | Get list of photos with dates.
   [/month_list](#2-one-month)|POST  |Get list of photos of a month.
  [/last_three_list](#3-last-3-month-including-today)|POST | Get list of photos of last 3 days including today.


# Response

```json
{
  "status": 200,
  "reqId": "sadfadsf-a345345-as456456-353456adsfa",
  "result": [
    {
      "reqId": "2016-1-3",
      "title": "A Starry Night of Iceland",
      "description": "On some nights, the sky is the best show in town. On this night, the sky was not only the best show in town, but a composite image of the sky won an international competition for landscape astrophotography. The featured winning image was taken in 2011 over Jökulsárlón, the largest glacial lake in Iceland.  The photographer combined six exposures to capture not only two green auroral rings, but their reflections off the serene lake. Visible in the distant background sky is the band of our Milky Way Galaxy and the Andromeda galaxy. A powerful coronal mass ejection from the Sun caused auroras to be seen as far south as Wisconsin, USA.  Solar activity over the past week has resulted in auroras just over the past few days.   Follow APOD on: Facebook,  Google Plus, or Twitter",
      "date": "2016-1-3",
      "urls": {
        "normal": "http://apod.nasa.gov/apod/image/1601/aurora_vetter_1080.jpg",
        "hd": "http://apod.nasa.gov/apod/image/1601/aurora_vetter_2000.jpg"
      }
    },
    {
      "reqId": "2016-1-4",
      "title": "Earthset from the Lunar Reconnaissance Orbiter",
      "description": "On the Moon, the Earth never rises -- or sets.  If you were to sit on the surface of the Moon, you would see the Earth just hang in the sky. This is because the Moon always keeps the same side toward the Earth. Curiously, the featured image does picture the Earth setting over a lunar edge.  This was possible because the image was taken from a spacecraft orbiting the Moon - specifically the Lunar Reconnaissance Orbiter (LRO). In fact, LRO orbits the Moon so fast that, from the spacecraft, the Earth appears to set anew about every two hours. The featured image captured one such Earthset about three months ago.  By contrast, from the surface of the Earth, the Moon sets about once a day -- with the primary cause being the rotation of the Earth. LRO was launched in 2009 and, while creating a detailed three dimensional map of the Moon's surface, is also surveying the Moon for water and possible good landing spots for future astronauts.   Free APOD Lectures: Editor to speak this coming weekend in Philadelphia and New York City",
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
