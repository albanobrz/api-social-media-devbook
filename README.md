# Social Media Devbook API

This API was developed to practice, learn, and improve backend development skills. Its a simple social media app, inspired on Twitter. The architecture was inspired by my team's project, the hexagonal pattern.

## Documentation

You can find the API documentation [here](https://documenter.getpostman.com/view/27691165/2s93mATewe). It provides detailed information about the API endpoints, request/response examples, and usage instructions.

## How to Run

You have two options to run the API: using Docker Compose or running it locally. Follow the steps below:

**Using Docker Compose:**

1. Make sure you have Docker installed on your machine.
2. Run the following commands:
   - `sudo docker compose build`
   - `sudo docker compose up`
3. Obtain the container IP by running: `sudo docker inspect mongo`. Set the obtained IP as the value for the `mongo_URI` environment variable in the `.env` file.

**Running Locally:**

1. Set up the `.env` file as per the provided example.
2. Run the API using your preferred Go development tool or by running the main package.

## Running Tests

To run the tests, use the following command:

`go test ./`
