# announce-backend

Backend RESTful API for the "announce" SA:MP server plugin offering a consumable JSON API for listing servers

*In development - subject to change*

### Important terms of agreement

- Every host that announces to the backend RESTful API is a free San Andreas Multiplayer service for San Andreas Multiplayer users.
- Any user is free to choose what service they want to use.
- The backend RESTful API needs to provide non-redundant*, complete and correct datasets regarding the requests it recieves.
- Every announcement has a non-biased lifetime that is acceptible by the backend RESTful API.
- Servers and the backend RESTful API can additionally agree on a custom announcement lifetime.
- Multiple requests in a short period of time can be rate-limited by the backend RESTful API to prevent heavy workload.
- Servers can't announce other servers.
- The backend RESTful API can obtain IP addresses from the San Andreas Multiplayer API, and obtain server data accordingly.

## API - Draft `v0`

### `https://samp.southcla.ws/v0/servers`

Returns a JSON array of minimal `server` objects - may be paginated if requests get too large. The reason the field names are so short is to minimise unnecessary network traffic. There are also only a handful of fields required for a list, so properties such as the rules list and the players are omitted from this response.

Example:

```json
[
    {
        "ip": "0.0.0.0:7272",
        "hn": "My awesome server",
        "pc": 2,
        "mp": 10,
        "gm": "tdm v1.0",
        "la": "French",
        "pa": false
    },
    {
        "ip": "1.1.1.1:7272",
        "hn": "My awesomer server",
        "pc": 2,
        "mp": 10,
        "gm": "RPG",
        "la": "English",
        "pa": true
    }
]
```

### `https://samp.southcla.ws/v0/server/{ip}`

Returns a `server` object for a particular server with more fields filled in. This can be used to show a user a more detailed overview of a server they may be interested in while only requesting the information when it's needed.

Example:

```json
{
    "ip": "0.0.0.0:7272",
    "hn": "My awesome server",
    "pc": 2,
    "mp": 10,
    "gm": "tdm v1.0",
    "la": "French",
    "pa": false,
    "ru": {
        "description": "My awesome server is really awesome and you should come play here because I am offering refunds for everyone, free stuff forever!",
        "website": "https://myawesomesampserver.com",
        "discord": "discord.me/myawesomeserver",
        "ts3": "ts3.myawesomeserver.com",
        "irc": "irc://freenode.net/myawesomeserver"
    },
    "pl": [
        "Southclaws",
        "Dogmeat"
    ]
}
```

### `https://samp.southcla.ws/v0/players/{ip}`

Returns an array of `string` objects representing player names for a particular server.

Example:

```json
[
    "Southclaws",
    "Sheriffic",
    "Avariam"
]
```


## Dev/Testing

You can test the endpoint by submitting a POST:

```bash
curl -XPOST localhost:7790/server/samp.southcla.ws -d '{
    "ip": "samp.southcla.ws",
    "hn": "My cool server",
    "pc": 12,
    "pm": 32,
    "gm": "SS",
    "la": "English",
    "pa": false,
    "ru": {
        "weburl": "http://southcla.ws"
    },
    "pl": [
        "steve",
        "bob",
        "laura"
    ]
}'
```

Then grab the data back with a GET:

```bash
curl localhost:7790/server/samp.southcla.ws
{"ip":"samp.southcla.ws","hn":"My cool server","pc":12,"pm":32,"gm":"SS","la":"English","pa":false,"ru":{"weburl":"http://southcla.ws"},"pl":["steve","bob","laura"]}
```