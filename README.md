# jazzApi

Example project from go.dev that I have expanded upon adding JWT authorization and user management.

Application reads api data from the file "data/example.json". 
Application reads user data from "data/userExample.json". 

I plan on refactoring this to use a database instead of JSON files.

The api has 6 endpoints.

## /login

Users authenticate through the `/login` endpoint with a POST request. The API accepts URL/JSON formats
There are 2 user roles "admin" and "user". Some endpoints are restricted to admin only.

Example Requests:
```bash    
    curl -X POST localhost/login --data 'username=user&password=pass'

    curl -X POST localhost/login -H 'Content-Type: application/json' --data '{"username":"user", "password":"pass"}'
```

Presenting valid credentials will return a JWT in the JSON response. These tokens remain valid for 1 hour.

```json
{
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ3NDA1MTksInJ..."
}
```

This token must be included in the "Authorization" header to requests to other endpoints.

## User Endpoints

### /albums
A GET request to the `/albums` endpoint returns all the albums in the file 'data/example.json' and any newly created albums in json format.

Example Request:
```bash
    curl localhost/albums -H 'Authorization: <JWT>'
```

Response:
```json
[
    {
        "id": 1,
        "title": "Blue Train",
        "artist": "John Coltrane",
        "price": 56.99
    },
    {
        "id": 2,
        "title": "Jeru",
        "artist": "Gerry Mulligan",
        "price": 12.99
    },
    {
        "id": 3,
        "title": "Sarah Vaughan and Clifford Brown",
        "artist": "Sarah Vaughan",
        "price": 39.99
    }
]
```

### /albums/:id
A GET request to the `/albums/:id` endpoint with a valid id number will return a json response containing that album.

Example Request:
```bash
    curl localhost/albums/1 -H 'Authorization: <JWT>
```

Response:
```json
{
    "id": 1,
    "title": "Blue Train",
    "artist": "John Coltrane",
    "price": 56.99
}
```

## Admin Endpoints

### /albums/create

A POST request to the `/albums/create` endpoint with valid JSOM data will create a new album, the Response will contain the newly created object.

Example Request:

```bash
    curl -X POST localhost/albums/create -H 'Authorization: <Admin JWT>' -H 'Content-Type: application/json' --data '{"id":4, "title": "example", "artist": "example", "price": 1.99}'
```

Response:

```json
{
    "id": 4,
    "title": "example",
    "artist": "example",
    "price": 1.99
}
```


### /users

A GET request to the `/users` endpoint with a valid JWT with the claim of role = admin will return a json response containing all of the users in the `data/userExample.json` file and any newly created users.

Example Request:

```bash
    curl localhost/users -H 'Authorization: <Admin JWT>'
```

Response:

```json
[
    {
        "id": 1,
        "username": "admin",
        "password": "21232f297a57a5a743894a0e4a801fc3",
        "role": "admin"
    },
    {
        "id": 2,
        "username": "testuser",
        "password": "5d9c68c6c50ed3d02a2fcf54f63993b6",
        "role": "user"
    }
]
```

### /users/create

A POST request to the `/users/create` endpoint with a valid admin token and valid url/json data will result in the creation of a new user. You must specify the username, password and role for the new user. The application will automatically generate the users ID and save the users password as an MD5 hash.

Example Requests:
```bash
    curl -X POST localhost/users/create -H 'Authorization: <Admin JWT>' --data 'username=test&password=test&role=user'

    curl -X POST localhost/users/create -H 'Authorization: <Admin JWT>' -H 'Content-Type: application/json' --data '{"username":"test", "password":"test", "role":"user"}'
```

Example Response:
```json
{
    "id": 3,
    "username": "test",
    "password": "test",
    "role": "user"
}
```
