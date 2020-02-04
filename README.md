# ddoseppuku

## Setup

You're gonna need Go installed and a proper `$GOPATH` and all that. Then:

```
go get -u github.com/michaeloverton/ddoseppuku
```

## Architecture

There are three layers, which operate as separate applications. Ideally, within the Docker network, the attack layer is separated from the defense layer. I'm not sure if this is working properly yet.

The `sentinel` application receives requests for attack or cease-fire. It then queues up messages accordingly in Redis. 

The `laser` application consumes from the Redis queues and acts accordingly. Multiple lasers can be scaled out horizontally (via `--scale laser=3` in the `build-all` target in the Makefile). This will spin up multiple containers with however many laser images you want. The lasers will make a maximum number of requests and then chill. The maximum number of requests each laser will make can be modified via `LSR_MAX_REQUESTS` in `/build/docker/docker-compose.yml`. Be careful making the max requests too high - at higher levels, the lasers can max out their own CPU. We are trying to destroy ourselves, but not in this layer. Also be careful scaling out too many lasers, or have fun cleaning up your mess.

The `target` application has two endpoints - a `/health` endpoint and a `/thrash` endpoint. The health endpoint always responds with 200 if the server is okay. The thrash endpoint mocks a task. Currently it reverses the text of Infinite Jest. The intensity of the task can be set via `TGT_TASK_INTENSITY` in the docker-compose. A single reversal takes about 30ms. A single increment in the task intensity basically doubles that value. The health endpoint normally responds in about 1ms, so if you slow that down, you're doing a good job.

![](/diagram.jpg?raw=true)

## Run

Have Docker running.

```
make build-all run-all
```

## Take A Hostage

`POST` to the sentinel at `http://localhost:3000/attack` to trigger an attack:

```
{
	"url":"http://target:3001/thrash"
}
```

curl:

```
curl --request POST \
  --url http://localhost:3000/attack \
  --header 'content-type: application/json' \
  --data '{
	"url":"http://target:3001/thrash"
}'
```

Stop the attack by sending a `GET` to `http://localhost:3000/ceasefire`.

Target `GET` endpoints:

```
http://localhost:3001/health
http://localhost:3001/thrash
```

You can monitor these while an attack is going on. They will be fucked.

## The Golden Goose

Destroy the target by attacking only the health endpoint.
