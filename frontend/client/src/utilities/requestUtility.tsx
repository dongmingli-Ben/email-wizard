/**
 * Utility functions to make API requests.
 * By importing this file, you can use the provided get and post functions.
 * You shouldn't need to modify this file, but if you want to learn more
 * about how these functions work, google search "Fetch API"
 *
 * These functions return promises, which means you should use ".then" on them.
 * e.g. get('/api/foo', { bar: 0 }).then(res => console.log(res))
 */

// ex: formatParams({ some_key: "some_value", a: "b"}) => "some_key=some_value&a=b"
function formatParams(params: Object): string {
  // iterate of all the keys of params as an array,
  // map it to a new array of URL string encoded key,value pairs
  // join all the url params using an ampersand (&).
  return Object.keys(params)
    .map((key) => key + "=" + encodeURIComponent(params[key]))
    .join("&");
}

// convert a fetch result to a JSON object with error handling for fetch and json errors
async function convertToJSON(res: Response) {
  if (!res.ok) {
    console.log(
      `API request failed with response status ${res.status} and text: ${res.statusText}`
    );
    let errData = await res.json();
    throw new Error(errData.errMsg);
  }

  return res
    .clone() // clone so that the original is still readable for debugging
    .json() // start converting to JSON object
    .catch((error) => {
      // throw an error containing the text that couldn't be converted to JSON
      return res.text().then((text) => {
        throw `API request's result could not be converted to a JSON object: \n${text}`;
      });
    });
}

// Helper code to make a get request. Default parameter of empty JSON Object for params.
// Returns a Promise to a JSON Object.
export async function get(endpoint: string, params = {}) {
  const fullPath = endpoint + "?" + formatParams(params);
  try {
    const res = await fetch(fullPath, {
      method: "GET",
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET,POST,OPTIONS,DELETE,PUT",
      },
    });
    return convertToJSON(res);
  } catch (error) {
    // give a useful error message
    throw `GET request to ${fullPath} failed with error:\n${error}`;
  }
}

// Helper code to make a post request. Default parameter of empty JSON Object for params.
// Returns a Promise to a JSON Object.
export async function post(endpoint: string, params = {}) {
  try {
    const res = await fetch(endpoint, {
      method: "post",
      headers: { "Content-type": "application/json" },
      body: JSON.stringify(params),
    });
    return convertToJSON(res);
  } catch (error) {
    // give a useful error message
    throw `POST request to ${endpoint} failed with error:\n${error}`;
  }
}

export const backendConfig = {
  verify_email: "https://47.243.42.37:8080/verify_email",
  events: "https://47.243.42.37:8080/events",
  add_user: "https://47.243.42.37:8080/add_user",
  verify_user: "https://47.243.42.37:8080/verify_user",
  add_mailbox: "https://47.243.42.37:8080/add_mailbox",
  user_profile: "https://47.243.42.37:8080/user_profile",
};
