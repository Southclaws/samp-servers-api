# announce-backend

[![Travis](https://img.shields.io/travis/Southclaws/announce-backend.svg)](https://travis-ci.org/Southclaws/announce-backend)[![Coverage](http://gocover.io/_badge/github.com/Southclaws/announce-backend)](http://gocover.io/github.com/Southclaws/announce-backend)

A SA:MP server listing API service. Anyone can POST a game server address which is added to a periodically queried queue and up-to-date information is provided as a JSON API.

## API - Draft `v2`

### GET `http://samp.southcla.ws/v2/servers`

Returns a JSON array of `server` objects. The reason the field names are so short is to minimise unnecessary network traffic. There are also only a handful of fields required for a list, so properties such as the rules list and the players are omitted from this response.

#### Query Parameters

- `sort`
  Either:
  - `asc`
  - `desc`
- `by`
  Either:
  - `players`
  - `rank`
- `filter`
  Comma separated set of filter parameters that reduce the amount of results:
  - `password` - removes passworded servers
  - `empty` - removes empty servers
  - `full` - removes full servers

#### Result

The result is a JSON array of `ServerCore` objects. A `ServerCore` includes just the information required to render a row on a list of servers.

```json
[
    {
        "ip": "sa-arp.net:7777",
        "hn": "[1994] SA Advanced Role-Play   (First Person)",
        "pc": 5,
        "pm": 160,
        "gm": "SA:ARP v3.2.1 r14 (Roleplay)",
        "la": "All (English)",
        "pa": false
    },
    {
        "ip": "server.redcountyrp.com:7777",
        "hn": "Red County Roleplay",
        "pc": 104,
        "pm": 150,
        "gm": "RC-RP 2.5.2 R1",
        "la": "English",
        "pa": false
    }
]
```

### GET `http://samp.southcla.ws/v2/server/{ip}`

Returns a `Server` object for a particular server with more fields filled in. This can be used to show a user a more detailed overview of a server they may be interested in while only requesting the information when it's needed.

#### Result

The result includes the `core` object used for rendering the server row on a list as well as additional information such as the players, the rules and some new fields made possible by this project offering more customisation for server owners.

Owners can define a `description` to help sell their server with more information about what's offered and a `banner` can be specified which is a URL to an image which can be rendered as part of a server browser implementation.

```json
{
    "core": {
        "ip": "151.80.108.109:8660",
        "hn": "\ufffd\ufffd\ufffd\ufffd | LOS SANTOS GANG WARS | \ufffd\ufffd\ufffd\ufffd",
        "pc": 29,
        "pm": 100,
        "gm": "Gang Wars/TDM/Turfs/mini",
        "la": "English/Espa\ufffdol",
        "pa": false
    },
    "ru": {
        "lagcomp": "On",
        "mapname": "San Andreas",
        "version": "0.3.7-R2",
        "weather": "12",
        "weburl": "samp-lsgw.com",
        "worldtime": "01:00"
    },
    "pl": [
        "Serega_Kgerth",
        "Eduardo_Sanchez",
        "SuperGamerxD"
    ],
    "description": "",
    "banner": ""
}
```

### POST `http://samp.southcla.ws/v2/server/{ip}`

This is how server owners provide information about their server. This is simply the reverse of the `GET` method on the same endpoint. The response body must be a `Server` object with the only required fields being:
- `ip` - IP or domain address
- `hn` - current hostname
- `pm` - maximum amount of players
- `gm` - current gamemode name

Aside from that, owners can specify any details they want inside the `ru` (Rules) field such as Discord URL, forum link, donation link, youtube videos, etc.

Note: when performing a POST to this endpoint, the `{ip}` in the URL and the `ip` field in the payload **must** match the IP of the request. For example, if you try to POST to `/v2/server/100.0.0.1` from the network `100.0.0.2` the request will fail with `400 BAD REQUEST`. This means only the server owner can update the records of their own server from the same physical machine/network as the actual server. If you use a more complex networking setup that this causes problems with, please open an issue and we can work something out.