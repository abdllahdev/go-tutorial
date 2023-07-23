## Running the application

To run the application use the following command

```
make docker-up
```

Folder structure

```
  project
    | cmd/app
      | main.go // starts the apps
    | internal
      | entity // where all the models, requests, and responses structs defined
      | guest_list // This is the guest list service
    | pkg/database // This package is responsible for connecting with and querying the DB
```
