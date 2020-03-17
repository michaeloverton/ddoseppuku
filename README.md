# ddoseppuku

## Setup

You're gonna need Go installed and a proper `$GOPATH` and all that. Then:

```
go get -u github.com/michaeloverton/ddoseppuku
```

## Architecture

There are three layers, which operate as separate applications.

The `sentinel` application receives requests for attack. It publishes these attack messages to Redis. The endpoint accepts either `GET` or `POST` attacks.

The `laser` application subscribes to the Redis attack topic. Multiple lasers can be scaled out horizontally (via `--scale laser=3` in the `build-all` target in the Makefile). This will spin up multiple containers with however many laser images you want. The lasers will make a maximum number of requests and then chill. The maximum number of requests each laser will make can be modified via `LSR_MAX_REQUESTS` in `/build/docker/docker-compose.yml`. Be careful making the max requests too high: at higher levels, the lasers can max out their own CPU. We are trying to destroy ourselves, but not in this layer. Also be careful scaling out too many lasers, or have fun cleaning up your mess. Hint: Restart Docker, then:

```
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
```

The `target` application has three endpoints: `GET /health`, `GET /thrash`, and `POST /login`. The health endpoint always responds with 200 if the server is okay. The thrash endpoint mocks a task. Currently it reverses the text of Infinite Jest. The intensity of the task can be set via `TGT_TASK_INTENSITY` in the docker-compose. A single reversal takes about 30ms. A single increment in the task intensity basically doubles that value. The health endpoint normally responds in about 1ms, so if you slow that down, you're doing a good job. The login endpoint accepts `POST` requests of the form:

```
{
  "username":"uGot",
  "password":"ddosedBro"
}
```

![](/diagram.jpg?raw=true)

## Run

Have Docker running.

```
make build run-all
```

## Take A Hostage

`POST` to the sentinel at `http://localhost:3000/attack` to trigger an attack pulse:

```
{
	"url":"http://localhost:3001/thrash",
	"method": "GET"
}
```

or

```
{
	"url":"http://localhost:3001/login",
	"method": "POST",
	"body": {
		"username": "uGot",
		"password": "ddosedBro"
	}
}
```

To perform a "continuous" attack, use the `suppress.sh` script. It takes the time between attack pulses as the first argument, and the URL to attack as the second. The rest of the attack request is hard-coded but can be easily modified in the script.

```
sh suppress.sh 30 http://target:3001/thrash
```

Target endpoints:

```
GET http://localhost:3001/health
GET http://localhost:3001/thrash
POST http://localhost:3001/login
```

You can monitor these while an attack is going on. They will be fucked.
