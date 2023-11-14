# wikimedia-pageviews-api

- list available tasks: `task --list`
- startup server: `task run` or `go run internal/main.go`
- generate executable binary: `task build` or `go build -o bin/wikimedia-pageviews-api internal/main.go`
- run tests: `task test` or `go test`
- see tests coverage: `task test.coverage` or `go test -coverprofile=c.out ./..`
- generate updated docs: `task swagger.doc` or `docker run -i yousan/swagger-yaml-to-html < pkg/swagger/swagger.yml > docs/index.html`

## Swagger API Documentation

See [docs](docs/index.html)

## Known Issues and Pending TODO

- If the article has special characters the URL validation for `ViewsPerArticleWeeklyHandler` fails

## Assumptions

- The endpoints that return the top articles for a week and a month return only the top 10. I chose to do this because it was easier to test with smaller result sets.
- The endpoints that return the top articles for a week and a month might return different results because they call different wikipedia endpoints (TEST THIS!!!)
- When I make calls to the Wikipedia API I consider anything different than HTTP 200 a failure, even though all 2xx codes are considered successful according to the IETF HTTP spec. This is done for simplicity reasons in the exercise context and it wouldn't be done in a real world scenario.
- The start of the week is assumed to be Monday (see Future Improvements for more).

## Future Improvements and Next Steps

- Improve test coverage, especially for handlers (currently the error cases are not covered)
- Documentation: move documentation in the code by adding swagger comments and be able to generate updated documentation. Use [go-swagger](https://github.com/go-swagger/go-swagger) to do that.
- Parallelize calls: the endpoint that returns the top articles for a given week has a loop that makes 7 calls to the wikipedia API, one for each day of the week.
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
- Note that there are more small improvements noted throughtout the codebase using `// Enhancement:` or `// TODO:`
