# Superman-Detector
This is an api that identifies logins by a user that occur from locations that are farther apart than a normal person can reasonably travel on an airplane. These locations are determined by looking up the source IP addresses of successful logins in a GeoIP database.

The api takes a POST request to the `/loginRequest` endpoint with a json payload that contains details of a users login attempt.

A response is returned in json format that contains location details of the current login attempt along with previous and subsequent location details of login attempts if they exist.

## Getting Started :
1. First step is to clone this repo into a directory of your choosing:
   
   `$ git clone https://github.com/mysteryboy73/Superman-Detector.git`

2. Once you've cloned the repo navigate to the src directory inside this repo:

    `$ cd */Superman-Detector/src`

3. This directory contains a *Dockerfile* and *docker-compose.yml* that we can use issuing the following command:

    `$ docker-compose up -d`
    > This command will build a docker image named *superman-detector* via the commands listed in the *Dockerfile*. It will then start a container using that image which will start a server listening on port *1210*. It will also volume map the location */db/geolite2.db* to the container to utilize the *sqlite* db.

4. When this container is running we can start to make requests to *http://localhost:1210/loginRequest*

## Making A Request :

**URL** : `http://localhost:1210/loginRequest`

**Method** : `POST`

**Request Body :**
> The request body is a json payload containing details about the users login attempt

```
{
	"username": "Bob",
	"unix_timestamp": 1611692790,
	"event_uuid": "926C4CDA-A5F0-19AB-C097-A61431CB4BFA",
	"ip_address": "1.0.212.0/23"
}
```
**Response :**

>The endpoint should return an HTTP Response code of *200 Ok* along with a json formatted response body including the following details

```
{
    "currentGeo": {
        "lat": 13.4,
        "lon": 100,
        "radius": 500
    },
    "travelToCurrentGeoSuspicious": true,
    "traveFromCurrentGeoSuspicious": false,
    "precedingIpAccess": {
        "ip": "1.10.160.0/22",
        "speed": 376,
        "lat": 16.4792,
        "lon": 104.6583,
        "radius": 500,
        "timestamp": 1611689190
    },
    "subsequentIpAccess": {
        "ip": "1.0.75.0/24",
        "speed": 2495,
        "lat": 34.4593,
        "lon": 132.4731,
        "radius": 1,
        "timestamp": 1611696390
    }
}
```

>Note : Not every request will return a payload with all the details. To recieve all details the current request being made must have a previous and subsequent login attempt to be able to calculate the neccessary fields.

## Assumptions / Design Decisions :

* With not be super framiliar with Go I tried to stick to what the standard library has to offer so I could get a better handle on what it is trying to accomplish without using a bunch of other libraries.

* Sticking with not trying to use a lot of other libraries I have a lack of an ORM for my repository struct. This means I'm writing sql statements directly in the code which can be inefficient, hard to maintain, and hard to understand.

* While writing the code my C# object oriented background seem to have influence on how I wrote things. Getting more comfortable with Go it seems that it's better suited to be a functional language.

* My code probably does not do the best job with error handling. I struggled with decided when errors should just be logged and when they should completely interrupt the app running.

* Thinking about how a users login attempts could continue to grow over time it seemed to me that you could be dealing with a large dataset at some point. That is why I thought it was better to keep the logic to get previous and subsequent login attempts within the database. My thought was if we could filtered our result set back before we brought it back then all we would have to do is map it to some struct.

## Third Party Library References : 

While I tried to leverage a lot of the built in tools the Go language has to offer I did use a few other packages to get the job done

* github.com/DATA-DOG/go-sqlmock - This package was great for creating a mock database that I could use during unit testing of my interactions with the database.
  
* github.com/jmoiron/sqlx - Sqlx provides extensions to Go's native database/sql library that allowed me to leverage some functions that made mapping db data to structs simple.
  
* github.com/mattn/go-sqlite3 - Provided a Sqlite3 driver so my code could interact with the sqlite db.
  
* github.com/umahmood/haversine - I did not feel the need to reinvent the wheel with the haversine formula. This package provides a function that takes coordinates of two desinations and reports back the distance.