# announce-backend

Backend RESTful API for the "announce" SA:MP server plugin offering a consumable JSON API for listing servers

*In development - subject to change*

## API - Draft `v0`

### `https://samp.southcla.ws/v0/servers`

Returns a JSON array of `server` objects - may be paginated if requests get too large

Example:

```json
[
    {
        "ip": "0.0.0.0:7272",
        "hn": "My awesome server",
        "gm": "tdm v1.0",
        "cn": "France",
        "ma": "San Fierro",
        "la": "French"
    },
    {
        "ip": "1.1.1.1:7272",
        "hn": "My awesomer server",
        "gm": "RPG",
        "cn": "UK",
        "ma": "San Fierro",
        "la": "English"
    }
]
```

### `https://samp.southcla.ws/v0/server/{ip}`

Returns a `server_detail` object for a particular server.

Example:

```json
{
    "server": {
        "ip": "0.0.0.0:7272",
        "hn": "My awesome server",
        "gm": "tdm v1.0",
        "cn": "France",
        "ma": "San Fierro",
        "la": "French"
    },
    "description": "My awesome server is really awesome and you should come play here because I am offering refunds for everyone, free stuff forever!",
    "website": "https://myawesomesampserver.com",
    "discord": "discord.me/myawesomeserver",
    "ts3": "ts3.myawesomeserver.com",
    "irc": "irc://freenode.net/myawesomeserver",
    "meta": {
        "arbitrary": "data",
        "that_is": "displayed",
        "on_the": "browser ui"
    }
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
