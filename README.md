# wikimedia-pageviews-api

## Run the API

```shell
docker build -t wikimedia-pageviews-api .
docker run --publish 8080:8080 wikimedia-pageviews-api
```

Send requests to `http://localhost:8080`.

## Other Commands

- list available tasks: `task --list`
- startup server: `task run` or `go run internal/main.go`
- generate executable binary: `task build` or `go build -o bin/wikimedia-pageviews-api internal/main.go`
- run tests: `task test` or `go test`
- see tests coverage: `task test.coverage` or `go test -coverprofile=c.out ./..`
- generate updated docs: `task swagger.doc` or `docker run -i yousan/swagger-yaml-to-html < pkg/swagger/swagger.yml > docs/index.html`

## API Documentation

- [Swagger docs](docs/index.html)
- [Postman collection](docs/wikipedia-pageviews-api.postman_collection.json)

## Known Issues and Pending TODO

- If the article has special characters the URL validation for `ViewsPerArticleWeeklyHandler` fails
- Add health endpoint?

## Assumptions

- The endpoints that return the top articles for a week and a month return only the top 10 (instead of the 1000 that Wikipedia gives us). This was done for convinience since it's easier to do the manual tests with smaller result sets.
- The endpoints that return the top articles for a week and a month will return different results because they call different wikipedia endpoints and often include different dates. For example, if we retrieve the top 10 articles for January 2023 using the `/articles/top/monthly/2023/01` endpoint the results will be for the days 1 to 31 of January
- When I make calls to the Wikipedia API I consider anything different than HTTP 200 a failure, even though all 2xx codes are considered successful according to the IETF HTTP spec. This is done for simplicity reasons in the exercise context and it wouldn't be done in a real world scenario.
- The start of the week is assumed to be Monday (see Future Improvements for more).

## Future Improvements and Next Steps

- Parallelize calls: the endpoint that returns the top articles for a given week has a loop that makes 7 calls to the wikipedia API, one for each day of the week.
- Use HTTPS. Currently the API uses HTTP but in a real world scenario we would encrypt the communication using SSL/TLS.
- Add healthcheck endpoint.
- Improve test coverage, especially for handlers (currently the error cases are not covered).
- Improve input validation. For example, the API does validate that the input for year, month, or week is numeric but it does not validate if that is a valid number (for example the week is not greater than 52).
- The regular expressions that match the URL called with the route could use some more work. The one validating the article name works for most wikipedia articles, but fails for some cases with weird characters. For example if you call try to get the pageviews for the article `https://en.wikipedia.org/wiki/Æthelred_the_Unready` you will get a 404 as the refexp cannot match the request to a route. However, if you URL-encode the input it will work: `http://localhost:8080/article/%25C3%2586thelred_the_Unready/weekly/2023/03`.
- Documentation: move documentation in the code by adding swagger comments and be able to generate updated documentation. Use [go-swagger](https://github.com/go-swagger/go-swagger) to do that.
- Improve error responses: currently when an error occurs the API returns a JSON doc containing one `Error` string. Refactor this to include more info properly structured. For example:
  ```json
  {
    "status": 400,
    "message": "detailed error",
    "more_info": "http://www.mydomain.com/link/to/error/docs"
  }
  ```
- Logging: currently the API simply logs in the console whenever an error occurs.
- Mocking: currently the API tests are making calls to the wikipedia API. We can moke the API and test against the expected results.
- Configuration file: move values like the wikipedia base URL, port number, etc in a configuration file
- Make the start of the week part of the API input. It could be Monday, Sunday, or Saturday (if the API was available to the Middle East or North Africa).
- The Wikipedia API has some rules that were not taken into consideration.
- API versioning.
- Note that there are more small improvements noted throughtout the codebase using `// Enhancement:` or `// TODO:`
