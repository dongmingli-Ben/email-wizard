/**
 * Utility functions to make API requests.
 * By importing this file, you can use the provided get and post functions.
 * You shouldn't need to modify this file, but if you want to learn more
 * about how these functions work, google search "Fetch API"
 *
 * These functions return promises, which means you should use ".then" on them.
 * e.g. get('/api/foo', { bar: 0 }).then(res => console.log(res))
 */
import axios, { AxiosResponse, AxiosError } from "axios";

// Helper code to make a get request. Default parameter of empty JSON Object for params.
// Returns a Promise to a JSON Object.
export async function get(endpoint: string, params = {}) {
  return axios
    .get(endpoint, {
      params: params,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET,POST,OPTIONS,DELETE,PUT",
      },
    })
    .then((resp: AxiosResponse) => {
      // console.log(resp);
      return resp.data;
    })
    .catch((error: AxiosError) => {
      console.log(error.response?.data);
      throw `GET request to ${endpoint} failed with error:\n${error}`;
    });
}

// Helper code to make a post request. Default parameter of empty JSON Object for params.
// Returns a Promise to a JSON Object.
export async function post(endpoint: string, params = {}) {
  return axios
    .post(endpoint, params, {
      headers: { "Content-type": "application/json" },
    })
    .then((resp: AxiosResponse) => {
      console.log(resp);
      return resp.data;
    })
    .catch((error) => {
      console.log(error.response?.data);
      throw `POST request to ${endpoint} failed with error:\n${error}`;
    });
}

export const backendConfig = {
  verify_email: "https://47.243.42.37:8080/verify_email",
  events: "https://47.243.42.37:8080/events",
  add_user: "https://47.243.42.37:8080/add_user",
  verify_user: "https://47.243.42.37:8080/verify_user",
  add_mailbox: "https://47.243.42.37:8080/add_mailbox",
  user_profile: "https://47.243.42.37:8080/user_profile",
};
