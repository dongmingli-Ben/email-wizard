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

const timeoutMillis = 100 * 1000;

// Helper code to make a get request. Default parameter of empty JSON Object for params.
// Returns a Promise to a JSON Object.
export async function get(endpoint: string, params = {}, extraHeaders = {}) {
  return axios
    .get(endpoint, {
      params: params,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET,POST,OPTIONS,DELETE,PUT",
        ...extraHeaders,
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
export async function post(endpoint: string, params = {}, extraHeaders = {}) {
  return axios
    .post(endpoint, params, {
      headers: { "Content-type": "application/json", ...extraHeaders },
      timeout: timeoutMillis,
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

export async function appGet(
  endpoint: string,
  user_id: number,
  user_secret: string,
  params = {}
) {
  return get(endpoint.replace("{user_id}", user_id.toString()), params, {
    "X-User-Secret": user_secret,
  });
}

export async function appPost(
  endpoint: string,
  user_id: number,
  user_secret: string,
  params = {}
) {
  return post(endpoint.replace("{user_id}", user_id.toString()), params, {
    "X-User-Secret": user_secret,
  });
}

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export const backendConfig = {
  verify_email: "https://www.toymaker-ben.online/api/verify_email",
  events: "https://www.toymaker-ben.online/api/users/{user_id}/events",
  add_user: "https://www.toymaker-ben.online/api/users",
  verify_user: "https://www.toymaker-ben.online/api/authenticate",
  add_mailbox: "https://www.toymaker-ben.online/api/users/{user_id}/mailboxes",
  user_profile: "https://www.toymaker-ben.online/api/users/{user_id}/profile",
};
