# gowatch

a command tool written in go to watch scoeps of HackerOne
- features:
    - get new programs
    - get new added scopes
    - notify in discord

## configuration
you should add a `.env` file, including your discord `Bot_Token`, your channelID `ChannelId_general` and your `webhook_URL` .

```json
Bot_Token="OTM***********YBC5k"
ChannelId_general="9********5"
WebHook_URL="https://discord.com/api/webhooks/{webhook.id}/{webhook.token}"
```

### schedule execution

- #### linux:
after you installed it, add this line in `/etc/crontab/` file.

example:
```
00 7 * * * <username> gowatch
```

- #### windows:
add a new task in `Task Scheduler`.

## installation

you need **go** and **mongodb** as a requirement on your machine.
```
go install github.com/blumid/gowatch@latest
```