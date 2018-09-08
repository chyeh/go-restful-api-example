## Run

### Docker Compose

To run the application by `docker-compose` on the `localhost`, simply run the following command:

```shell
docker-compose up -d
```

### Binary Executable

Make sure the PostgreSQL service for the application is up. In the following example, the PostgreSQL service is running on `localhost:5432`. To run the application, build the application by `go build` and run the application by the following command:

```shell
go build -o app
./app --dsn "postgres://hellofresh:hellofresh@localhost:5432/hellofresh?sslmode=disable"
```

The description of the flags is as follows:

|   Flag   | Type       | Description                                                  |
| :------: | ---------- | ------------------------------------------------------------ |
| `--dsn`  | **string** | PostgreSQL database connection string. It **must be set** or the application occurs panic. |
| `--host` | **string** | Host that the http service binds to.                         |
| `--port` | **string** | Port that the http service listens to. The default value is `8080`. |



## Test

### Unit Test

The unit test requires the connection to the database. The PostgreSQL host is specify by the flag `-test.db.host`. If the flag is not set, all the data layer tests are skipped. The data layer tests creates and drops a database `test_hellofresh`. The following command is an example of connecting to the PostgreSQL host running on `localhost:5432`:

```shell
go test -v -ginkgo.v -test.db.host postgres://hellofresh:hellofresh@localhost:5432/
```

### Integration Test

The integration test runs the application on the `localhost` with `docker-compose up -d`. Run `docker rm -f $(docker ps -aqf "name=chyeh-api-test")` in advance to make sure the current `localhost` is no running instance started by previous `docker-compose up -d` command . To run the integration test, run the script at the root directory of the project:

```
scripts/integration-test.sh
```

### Manual Test

To test the application manually, run scripts that set up the environment on `localhost` in advance:

```shell
docker-compose up -d
cd scripts
./init-db-schema.sh postgres://hellofresh:hellofresh@localhost:5432/hellofresh
./init-db-user-data.sh postgres://hellofresh:hellofresh@localhost:5432/hellofresh
```

The `init-db-schema.sh` script creates the schema in the database. The `init-db-user-data.sh` inserts several users and their access tokens for testing as follows:

| User       | Access Token                   |
| ---------- | ------------------------------ |
| hellofresh | `aGVsbG9mcmVzaDpoZWxsb2ZyZXNo` |
| chyeh      | `Y2h5ZWg6Y2h5ZWg=`             |
| foo        | `Zm9vOmJhcg==`                 |
| user       | `dXNlcjpwYXNzd29yZA==`         |

For testing manually by `curl` command, examples are in the `scripts/integration-test.sh` file.



## API Specifications

There are several terms used in the following. The description are as follow:

* `Protected`: For the API endpoints that are marked as `protected`, the access token must be set with the key `Authorization` in the **HTTP request header**.

* `Mandatory`: The following request arguments that are marked as `Mandatory` causes `500 internal server error` response if not set.

* **boolean**: The following request arguments that are marked as type **boolean** accept `1`, `t`, `T`, `TRUE`, `true`, `True` as **true** value and `0`, `f`, `F`, `FALSE`, `false`, `False` as **false** value.

* `RECIPE JSON` & `RECIPE JSON ARRAY`:

  The following JSON data is an example of a HTTP response body from the API endpoints that marked with `RECIPE JSON`.

  ```json
  {
      "id":1,
      "name":"name1",
      "prepare_time":null,
      "difficulty":null,
      "is_vegetarian":false,
      "rating": 0,
      "rated_num": 0
  }
  ```

  The following JSON data is an example of a HTTP response body from the endpoint of searching recipes that marked with `RECIPE JSON ARRAY`:

  ```json
  [
      {
          "id":1,
          "name":"name1",
          "prepare_time":null,
          "difficulty":null,
          "is_vegetarian":false,
          "rating": 0,
          "rated_num": 0
      },
      {
          "id":11,
          "name":"name11",
          "prepare_time":1,
          "difficulty":2,
          "is_vegetarian":true,
          "rating": 0,
          "rated_num": 0
      }
  ]
  ```

  The description of the fields in the JSON data are as follows:

  * `id`: The ID of the recipe.
  * `name`: The name of the recipe.
  * `prepare_time`: The preparation time of the recipe. The unit of the time is minute.
  * `difficulty`: The difficulty of the recipe.
  * `is_vegetarian`: Specify if the recipe is vegetarian or not.
  * `rating`: The current rating of the recipe.
  * `rated_num`: The number of times the recipe is being rated.

### `GET /recipes`: Search Recipes

#### Request

Pagination and filtering arguments are supported. 

##### Paging

Paging arguments are defined in the **HTTP request header**.

| Argument      | Type        | Description                                                  |
| ------------- | ----------- | ------------------------------------------------------------ |
| `page-number` | **integer** | Specify the page number. Negative value causes `500 internal server error` response. An empty string value is consider not set. The default value is `1`. |
| `page-size`   | **integer** | Specify the number of recipes in each page. Negative value causes `500 internal server error` response. An empty string value is consider not set. The default value is `20`. |

##### Filtering

Filtering arguments are defined in the **URL query string**.

| Argument            | Type        | Description                                                  |
| ------------------- | ----------- | ------------------------------------------------------------ |
| `name`              | **string**  | Find recipes whose names **contains** the value. An empty string value is consider not set. |
| `prepare_time_from` | **integer** | Find recipes whose preparation time is **greater than or equal to** the value |
| `prepare_time_to`   | **integer** | Find recipes whose preparation time is **less than or equal to** the value |
| `difficulty_from`   | **integer** | Find recipes whose difficulty time is **greater than or equal to** the value |
| `difficulty_to`     | **integer** | Find recipes whose preparation time is **less than or equal to** the value |
| `is_vegetarian`     | **boolean** | Find recipes which are **vegetarian** or **not vegetarian**. An empty value is consider not set. An invalid **boolean** value causes `400 bad request` response. |

#### Response `RECIPE JSON ARRAY`

The HTTP response body contains the result of the search according to the paging and filtering arguments.

### `POST /recipes`: Add a New Recipe `Protected`

#### Request

The aruments are defined by **JSON data** in the HTTP request.

| Field           | Type        | Description                                                  | Description |
| --------------- | ----------- | ------------------------------------------------------------ | ----------- |
| `name`          | **string**  | `Mandatory` An empty string value is consider not set.       |             |
| `prepare_time`  | **integer** | The value must be **greater than or equal to** `1` or it causes `500 internal server error` response |             |
| `difficulty`    | **integer** | The value must be **greater than or equal to** `1` and **less than or equal to** `3` or it causes `500 internal server error` response |             |
| `is_vegetarian` | **boolean** | `Mandatory` An invalid **boolean** value causes `400 bad request` response. |             |

#### Response `RECIPE JSON`

The HTTP response body contains the data of the recipe that is just added.

### `GET /recipes/{id}`: Get an Existent Recipe

#### Request

The argument of the recipe ID is defined by the **URL parameter**.

| Type        | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| **integer** | If there is no recipe that has an ID matching the value of the argument, it responses with `404 not found`. |

#### Response `RECIPE JSON`

The HTTP response body contains the data of the specified recipe.

### `PUT /recipes/{id}`: Modify an Existent Recipe `Protected`

#### Request

The argument of the recipe ID is defined by the **URL parameter**.

| Type        | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| **integer** | If there is no recipe that has an ID matching the value of the argument, it responses with `404 not found`. |

The aruments for updating the recipe are defined by the **JSON data** in the HTTP request body.

| Field           | Type        | Description                                                  |
| --------------- | ----------- | ------------------------------------------------------------ |
| `name`          | **string**  | An empty string value causes `500 internal server error` response. |
| `prepare_time`  | **integer** | The value must be **greater than or equal to** `1` or it causes `500 internal server error` response. |
| `difficulty`    | **integer** | The value must be **greater than or equal to** `1` and **less than or equal to** 3 or it causes `500 internal server error` response. |
| `is_vegetarian` | **boolean** | An invalid **boolean** value causes `400 bad request` response. |

#### Response `RECIPE JSON`

The HTTP response body contains the data of the recipe that is just modified.

### `DELETE /recipes/{id}`: Delete an Existent Recipe `Protected`

#### Request

The argument of the recipe ID is defined by the **URL parameter**.

| Type        | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| **integer** | If there is no recipe that has an ID matching the value of the argument, it responses with `404 not found`. |

#### Response `RECIPE JSON`

The HTTP response body contains the data of the recipe that is just deleted.

### `POST /recipes/{id}/rating`: Rate an Existent Recipe

#### Request

The argument of the recipe ID is defined by the **URL parameter**.

| Type        | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| **integer** | If there is no recipe that has an ID matching the value of the argument, it responses with `404 not found`. |

The arument of rating the recipe is defined by **JSON data** in the HTTP request.

| Field    | Type        | Description                                                  |
| -------- | ----------- | ------------------------------------------------------------ |
| `rating` | **integer** | `Mandatory` The value must be **greater than or equal to** `1` and **less than or equal to** `5`. |

#### Response `RECIPE JSON`

The HTTP response body contains the data of the recipe that is just rated.