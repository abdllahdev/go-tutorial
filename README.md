# GetGround Tech Task 2023 - Abdullah Elsayed

> Note: It's my first time using Go so there might be better ways to implement things that I'm not aware of. Thanks :)

## Running the application

To run the application use the following command

```md
make docker-up
```

## Folder structure

```md
project
| cmd/app
| | main.go
| internal
| | entity // This is where all the models, requests, and responses structs defined
| | guest_list // This is the guest list service and it includes comprehensive unit tests
| test
| api.go // Implements a helper function that is used to test API endpoints
| pkg/database // This package is responsible for connecting with and querying the DB
```

## Possible improvements in my code

- Move secrets to .env files
- Better error handling. The API returns error messages as strings without detailed information about the error
- Implement more tests in the postman collection

## Possible improvements to the API

The some of the endpoints don't follow the conventional patterns used to implement REST APIs for example the `POST /guest_list/{name}` endpoint in the conventional way should be implemented as `POST /guests` and the `name` of the guest should be submitted as part of the request body. Also, the use of the guest `name` as an identifier is problematic because there can be more than one guests called the same `name`.

So, if I was to implement the endpoints in the API guide, I would implement the following endpoints

- Create a new table `POST /tables`
  - ReqBody: `{ ...same_fields }`
- Count seats `GET /tables/empty-seats`
- Add new guest `POST /guests`
  - ReqBody: `{ "name": string, ...other_fields }`
- Get guests: `GET /guests`, to filter the list of guests we can use string query instead of creating two endpoints to retrieve many guests from the DB
- Check in guest `PUT /guests/{id}`
  - ReqBody: `{ ...same_fields }`
- Check out guest `DELETE /guests/{id}`
