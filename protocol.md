# Protocol description

## Master commands

| id  | name                          |
| --- | ----------------------------- |
| 1   | [AuthRequired](#authrequired) |
| 2   | [AuthSuccess](#authsuccess)   |
| 3   | [RoomCreate](#roomcreate)     |
| 4   | [RoomCancel](#roomcancel)     |

### AuthRequired

ID: 1

Events:

- on connection
- on Auth, error

Payload: error text, string

Occurs on connection/reconnection to notify session server for authorization. Session server should keep existing rooms.

### AuthSuccess

ID: 2

Events:

- on Auth, successfull

Payload: none

After this command master could create rooms on session server.

### RoomCreate

ID: 3

Event:

- on external request

Payload: room id, client ids

Request for new room with specified id and clients.

### RoomCancel

ID: 4

Event:

- on external request
- on room errors

Payload: room id

Request for new room with specified id and clients.

## Session server commands

| id  | name                          |
| --- | ----------------------------- |
| 1   | [Auth](#auth)                 |
| 2   | [RoomCreated](#roomcreated)   |
| 3   | [RoomError](#roomerror)       |
| 4   | [RoomFinished](#roomfinished) |
| 5   | [Stats](#stats)               |

### Auth

ID: 1

Events:

- on AuthRequired

Payload: auth token, byte array

Authorization or error handling.

### RoomCreated

ID: 2

Event:

- on RoomCreate, successfull

Payload: room id; endpoint; clients id, token

After successfull room creation.

### RoomError

ID: 3

Event:

- on RoomCreate, error

Payload: room id; error text

After room creation with error.

### RoomFinished

ID: 4

Event:

- on room closed by room processor

Payload: room id; clients id, room result

### Stats

ID: 5

Event:

- on room count change
- on server stop requested
- periodic

Payload: capacity (how many rooms can be created)

Periodic report to the master. If required shutdown, then session server should report zero capacity and process all existing rooms until finish.
