# react-pairing-challenge
a React pairing challenge for prospective candidates

## Intro
This is a simple API that we require a new front end for. The API and the challenge will be explained in detail at the start of the meeting.

## Build
You will need Docker installed on your computer for this challenge. Details can be found here... https://docs.docker.com/engine/install/

To build the API, open your CLI in the root of this React-Pairing-Challenge repo and run the command...
```
make build
```

## Run
To run the API while still in the root of this repo run the command...
```
make start
```
If you want to see start up logs run...
```
make start_with_logs
```
To stop the API at any point run the command...
```
make stop
```

Once the program is running, Swagger documentation reagarding the API will be available at http://localhost:8083/swagger/index.html