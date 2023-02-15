# Doppler Example

## Setup:

A Doppler account is needed and a project with the following config values set:
 - 
 - HTTP_ADDR - Address (port to listen on - :5000, 8080, etc.)
 - ROACH_CONN - The connction string to connect to CockroachDB over `postgres://` wire)
 - ROACH_DB - Database name to show connection to.
 - DOPPLER_TOKEN - the service or personal token with access to the project end environment secrets you want

 1. Install the Doppler CLI:

     [Install instructions](https://docs.doppler.com/docs/install-cli)

 2. Login to Doppler CLI:
     ```sh
     doppler login
     ```
 3. Setup the project:
    ```sh
    doppler setup
    ```
 4. Install and setup Ngrok

     [Install instructions](https://ngrok.com/docs/getting-started)
    

## Running the demo app

 To run and test webhooks locally, a tool like Ngrok is needed to open a secure tunnel to the outside world without exposing your local machine

 ```sh
 doppler run --commend="ngrok http -region=us $HTTP_ADDR"
 ```

Once everything is setup(ngrok and doppler CLI downloaded, logged in and setup):
 ```sh
 doppler run -- go run main.go
 ```

## Presentation Slides:

Slides are made with Slidev.js and can be run by `cd`ing into `slides/` 

Running `npm install` and then `npm run dev`

The slides should then show in your browser


### ToDo's: 
- [ ] Add a Makefile to set everything up (CLI and Ngrok download)
- [ ] Make table of data on "/" dynamically change as DB is changed.
