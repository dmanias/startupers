// import necessary modules
import { check } from 'k6';
import http from 'k6/http';

export default function () {
  // Get the connection string from the command line arguments or use the default value
  const connectionString = __ENV.CONNECTION_STRING || 'http://localhost:9000';

  // Define the URL using the provided connection string and path
  const url = `${connectionString}/users/1`;

  // Send a GET request and save response as a variable
  const res = http.get(url);

  // Check that response is 200
  check(res, {
    'response code was 200': (res) => res.status == 200,
  });

  // Parse the response body to JSON format
  const responseBody = JSON.parse(res.body);

  // Check the properties of the response
  check(responseBody, {
    'id is correct': (body) => body.id === '1',
    'name is correct': (body) => body.name === 'John',
  });
}

// Define test options
export let options = {
  vus: 100, // Number of virtual users to simulate
  duration: '10s', // Duration of the test in seconds
};
