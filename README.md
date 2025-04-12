# jazzApi

Example project from go.dev.

Application reads json from the file "example.json" then creates an api using gin web framework.

The api has three endpoints.

### GET
/albums - returns all albums

/albums/:id - returns the album with the corresponding id

### POST
/create - creates a new album with json POST data
