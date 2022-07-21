# pointsystems

To run this application, you will need a couple dependencies:

* Install [Encore](https://encore.dev/docs/install)
* Have a recent version of [Golang](https://go.dev)
* Docker will need to be running

## Running the application

```shell
encore run
```

should start a local environment with a Postgres database and server.  You should see a banner after everything starts up
with an admin dashboard and URI for the API:

```
  ✔ Building Encore application graph... Done!
  ✔ Analyzing service topology... Done!
  ✔ Creating PostgreSQL database cluster... Done!
  ✔ Generating boilerplate code... Done!
  ✔ Compiling application source code... Done!
  ✔ Running database migrations... Done!
  ✔ Starting Encore application... Done!

  Encore development server running!

  Your API is running at:     http://localhost:4000
  Development Dashboard URL:  http://localhost:50568/pointsystems-zat2
```

Your URL may vary.

## Testing API calls

Encore builds with it an admin dashboard that you can use to make API calls.  The right side will show Schema and have a selector
for making a Call directly.  Simply replace `:id` with `1` if you're running locally.

![img.png](img.png)
