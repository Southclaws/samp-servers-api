# Notes

## Requirements

- Compact
- Portable and human readable
- Non-redundant*
- Complete


## Important terms of agreement

- Every host that announces to the backend RESTful API is a free San Andreas Multiplayer service for San Andreas Multiplayer users.
- Any user is free to choose what service they want to use.
- The backend RESTful API needs to provide non-redundant*, complete and correct datasets regarding the requests it recieves.
- Every announcement has a non-biased lifetime that is acceptible by the backend RESTful API.
- Servers and the backend RESTful API can additionally agree on a custom announcement lifetime.
- Multiple requests in a short period of time can be rate-limited by the backend RESTful API to prevent heavy workload.
- Servers can't announce other servers.
- The backend RESTful API can obtain IP addresses from the San Andreas Multiplayer API, and obtain server data accordingly.


## Overview of interfaces

- Servers can send data on an agreed format to the backend RESTful API.
- Users can send requests on an agreed format to the backend RESTful API to recieve a list of servers on an agreed format of San Andreas Multiplayer servers.

## API

### Server to backend RESTful API (Announce)

The communication between the server and the backend RESTful API uses JSON (JavaScript Object Notation) that consists of these attributes:

- Host: Domain name or IPv4 address
- Hostname: Hostname of the server
- Gamemode: Gamemode of the server
- Language: Language of the server
- MaxPlayers: Maximum amount of players that the server supports
- Lifetime: Number of seconds to life as integer, -1 for default lifetime
- Extended: This attribute is an object that cotains multiple attributes like for example website URL, contact email and etc.

The message is contained in a HTTP POST request "/servers/{domain name or IPv4 address and port}" for example "/servers/127.0.0.1:7777"


### Backend RESTful API to server (Response)

The communication between the backend RESTful API and the server uses JSON (JavaScript Object Notation) that consists of these attributes:

- Message: The message related to the current status


### User to backend RESTful API (Request server list)

The communication between the user and the backend RESTful API has no body.

The message is contained in a HTTP GET request "/servers"


### Backend RESTful API to user (Server list response)

The communication between the backend RESTful API and the user uses JSON (JavaScript Object Notation) that consists of these attributes for each server:

- Host: Domain name or IPv4 address
- Hostname: Hostname of the server
- Gamemode: Gamemode of the server
- Language: Language of the server
- MaxPlayers: Maximum amount of players that the server supports

This object is an element structure of an array.


### User to backend RESTful API (Request server information)

The communication between the user and the backend RESTful API has no body.

The message is contained in a HTTP GET request "/servers/{domain name or IPv4 address and port}" for example "/servers/127.0.0.1:7777"


### Backend RESTful API to user (Server information response)

The communication between the backend RESTful API and the user uses JSON (JavaScript Object Notation) that consists of these attributes:

- Host: Domain name or IPv4 address
- Hostname: Hostname of the server
- Gamemode: Gamemode of the server
- Language: Language of the server
- MaxPlayers: Maximum amount of players that the server supports
- Extended: This attribute is an object that cotains multiple attributes like for example website URL, contact email and etc.


## Side notes

- \* Non-redundant means there are no duplicated entries
