# Superman-Detector
This is an api that identifies logins by a user that occur from locations that are farther apart than a normal person can reasonably travel on an airplane. These locations are determined by looking up the source IP addresses of successful logins in a GeoIP database.

The api takes a POST request to the `/loginRequest` endpoint with a json payload that contains details of a users login attempt.

A response is returned in json format that contains location details of the current login attempt along with previous and subsequent location details of login attempts if they exist.

## Getting Started
1. First step is to clone this repo into a directory of your choosing:
   
   `$ git clone https://github.com/mysteryboy73/Superman-Detector.git`

2. Once you've cloned the repo navigate to the src directory inside this repo:

    `$ cd */Superman-Detector/src`

3. This directory contains a *Dockerfile* and *docker-compose.yml* that we can use issuing the following command:

    `$ docker-compose up -d`
    > This command will build a docker image named *superman-detector* via the commands listed in the *Dockerfile*. It will then start a container using that image which will start a server listening on port *1210*. It will also volume map the location */db/geolite2.db* to the container to utilize the *sqlite* db.

4. When this container is running we can start to make requests to *http://localhost:1210/loginRequest*

## Making Request
