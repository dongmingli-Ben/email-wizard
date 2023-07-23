import { msalInstance } from "../index";
import {
  InteractionRequiredAuthError,
  AuthenticationResult,
} from "@azure/msal-browser";
import { graphConfig, tokenRequest, loginRequest } from "./msalAuthConfig";
import { backendConfig, get } from "./requestUtility";

async function callMSGraph(endpoint: string, token: string) {
  const headers = new Headers();
  const bearer = `Bearer ${token}`;

  headers.append("Authorization", bearer);

  const options = {
    method: "GET",
    headers: headers,
  };

  console.log("request made to Graph API at: " + new Date().toString());

  await fetch(endpoint, options)
    .then((response) => response.json())
    .then((response) => console.log(response))
    .catch((error) => console.log(error));
}

function getTokenPopup(request, username: string) {
  /**
   * See here for more info on account retrieval:
   * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-common/docs/Accounts.md
   */
  request.account = msalInstance.getAccountByUsername(username);

  return msalInstance.acquireTokenSilent(request).catch((error) => {
    console.warn("silent token acquisition fails. acquiring token using popup");
    if (error instanceof InteractionRequiredAuthError) {
      // fallback to interaction when silent call fails
      return msalInstance
        .acquireTokenPopup(request)
        .then((tokenResponse) => {
          console.log(tokenResponse);
          return tokenResponse;
        })
        .catch((error) => {
          console.error(error);
          return String(error);
        });
    } else {
      console.warn(error);
      return String(error);
    }
  });
}

function seeProfile(username: string) {
  getTokenPopup(loginRequest, username)
    .then((response) => {
      if (typeof response !== "string") {
        callMSGraph(graphConfig.graphMeEndpoint, response.accessToken);
        return;
      }
      console.log("err: auth token not obtained");
      console.log(response);
    })
    .catch((error) => {
      console.error(error);
    });
}

async function readMail(username: string) {
  await getTokenPopup(tokenRequest, username)
    .then(async (response) => {
      if (typeof response !== "string") {
        await callMSGraph(graphConfig.graphMailEndpoint, response.accessToken);
        return;
      }
      console.log("err: fail to read emails");
      console.log(response);
    })
    .catch((error) => {
      console.error(error);
      console.log(error);
    });
}

const verifyOutlook = async (address: string): Promise<string> => {
  let req = loginRequest;
  req.loginHint = address;
  let errMsg = "";
  let username: string;
  console.log(`msal login:`);
  console.log(req);
  await msalInstance.loginPopup(req).then(async (response) => {
    console.log("logged user in");
    if (response !== null) {
      if (response.account === null) {
        console.log("empty username from msal login response");
        errMsg = "error: empty username from msal login response";
        return;
      }
      username = response.account.username;
      await readMail(username);
    } else {
      console.log("null response from msal");
      errMsg = "error: null response from msal";
    }
  });

  return errMsg;
};

const verifyIMAP = async (
  address: string,
  password: string,
  imapServer: string
): Promise<string> => {
  let resp = await get(backendConfig.verify_email, {
    username: address,
    password: password,
    imap_server: imapServer,
    type: "IMAP",
  }).catch((e) => {
    console.log("error in verifyIMAP:", e);
    return { errMsg: "Fail to verify your email. Please check your inputs." };
  });
  console.log(resp);
  if ("errMsg" in resp) {
    return resp.errMsg;
  }
  return "";
};
const verifyPOP3 = async (
  address: string,
  password: string,
  pop3Server: string
): Promise<string> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return "";
};

export { verifyOutlook, verifyIMAP, verifyPOP3 };
