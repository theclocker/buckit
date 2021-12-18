# Api manager

### An application that manages apis and provides a single end-point for developers to use instead of calling it directly, or handling api calls limits

---
* Allow users to register apis, and use the software as a proxy
* Allow users to configure an api’s limitations
* Allow users to send bulk requests to the proxy and let it manage itself
* Log API requests timing, and network requests
* Allow management of users for each API, allowing users to disable or enable access to certain endpoints (mapped or not)
* Allow users to map API endpoints, and abstract variables
* Allow users to define an endpoint to either return the api results or an indicator that the api results are available, with an endpoint to download (or stream) the results
* A user interface with api management and statistics

---
## Todo
* Move configurations to the database
* Keep requests and responses in a database
* Add the forwarding of API responses to defined endpoints using a messaging engine
* Add headers to API requests
* Add configuration overriding to requests
* Add post requests (bulk / singles)