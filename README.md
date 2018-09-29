# samp-servers-api

[![Travis](https://img.shields.io/travis/Southclaws/samp-servers-api.svg)](https://travis-ci.org/Southclaws/samp-servers-api)

A SA:MP server listing API service. Anyone can POST a game server address which
is added to a periodically queried queue and up-to-date information is provided
as a JSON API.

---

# v2

This is an automatically generated documentation page for the v2 API endpoints.

## serverAdd

`POST`: `/v2/server/{address}`

Add a server to the index using just the IP address. This endpoint requires no
body and no additional information. The IP address is added to an internal queue
and will be queried periodically for information via the legacy server API. This
allows any server to be added with the basic information provided by SA:MP
itself.

## serverPost

`POST`: `/v2/server`

Provide additional information for a server such as a description and a banner
image. This requires a body to be posted which contains information for the
server.

### Accepts

```json
{
  "core": {
    "ip": "127.0.0.1:7777",
    "hn": "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
    "pc": 32,
    "pm": 128,
    "gm": "Grand Larceny",
    "la": "English",
    "pa": false,
    "vn": "0.3.7-R2"
  },
  "ru": {
    "lagcomp": "On",
    "mapname": "San Andreas",
    "version": "0.3.7-R2",
    "weather": "10",
    "weburl": "www.sa-mp.com",
    "worldtime": "10:00"
  },
  "description": "An awesome server! Come and play with us.",
  "banner": "https://i.imgur.com/Juaezhv.jpg",
  "active": true
}
```

## serverGet

`GET`: `/v2/server/{address}`

Returns a full server object using the specified address.

### Returns

```json
{
  "core": {
    "ip": "127.0.0.1:7777",
    "hn": "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
    "pc": 32,
    "pm": 128,
    "gm": "Grand Larceny",
    "la": "English",
    "pa": false,
    "vn": "0.3.7-R2"
  },
  "ru": {
    "lagcomp": "On",
    "mapname": "San Andreas",
    "version": "0.3.7-R2",
    "weather": "10",
    "weburl": "www.sa-mp.com",
    "worldtime": "10:00"
  },
  "description": "An awesome server! Come and play with us.",
  "banner": "https://i.imgur.com/Juaezhv.jpg",
  "active": true
}
```

## serverList

`GET`: `/v2/servers`

Returns a list of servers based on the specified query parameters. Supported
query parameters are: `page` `sort` `by` `filters`.

### Query parameters

Example: `by=player&filters=full&filters=password&page=2&sort=asc`

### Returns

```json
[
  {
    "ip": "127.0.0.1:7777",
    "hn": "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
    "pc": 32,
    "pm": 128,
    "gm": "Grand Larceny",
    "la": "English",
    "pa": false,
    "vn": "0.3.7-R2"
  },
  {
    "ip": "127.0.0.1:7777",
    "hn": "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
    "pc": 32,
    "pm": 128,
    "gm": "Grand Larceny",
    "la": "English",
    "pa": false,
    "vn": "0.3.7-R2"
  },
  {
    "ip": "127.0.0.1:7777",
    "hn": "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
    "pc": 32,
    "pm": 128,
    "gm": "Grand Larceny",
    "la": "English",
    "pa": false,
    "vn": "0.3.7-R2"
  }
]
```

## serverStats

`GET`: `/v2/stats`

Returns a some statistics of the server index.

### Returns

```json
{
  "servers": 1000,
  "players": 10000,
  "players_per_server": 10
}
```
