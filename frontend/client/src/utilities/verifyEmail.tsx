import { msalInstance } from "../index";
import { graphConfig, tokenRequest, gmailConfig } from "./authConfig";
import { backendConfig, get } from "./requestUtility";

type VerifyResposne = {
  errMsg: string;
  credentials: { [key: string]: string };
};

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
    // if (error instanceof InteractionRequiredAuthError) {
    // fallback to interaction when silent call fails
    request.loginHint = username;
    console.log(request);
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
    // } else {
    //   console.warn(error);
    //   return String(error);
    // }
  });
}

async function readMail(username: string, token: string) {
  await callMSGraph(graphConfig.graphMailEndpoint, token);
}

const getAccessToken = async (address: string): Promise<string> => {
  let token = await getTokenPopup(tokenRequest, address)
    .then(async (response) => {
      if (typeof response !== "string") {
        return response.accessToken;
      }
      console.log("err: fail to get access token");
      console.log(response);
      return "";
    })
    .catch((error) => {
      console.error(error);
      console.log(error);
      return "";
    });
  return token;
};

const verifyOutlook = async (address: string): Promise<VerifyResposne> => {
  let errMsg = "";
  await getAccessToken(address).then(async (token) => {
    if (token.length === 0) {
      console.log("fail to get access token, got:", token);
      return;
    }
    await readMail(address, token);
  });

  return { errMsg: errMsg, credentials: {} };
};

const verifyIMAP = async (
  address: string,
  password: string,
  imapServer: string
): Promise<VerifyResposne> => {
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
    return { errMsg: resp.errMsg, credentials: {} };
  }
  return {
    errMsg: "",
    credentials: {
      imap_server: imapServer,
      password: password,
    },
  };
};
const verifyPOP3 = async (
  address: string,
  password: string,
  pop3Server: string
): Promise<VerifyResposne> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return { errMsg: "", credentials: {} };
};

const verifyGmail = async (address: string): Promise<VerifyResposne> => {
  let oauth_url = `https://accounts.google.com/o/oauth2/v2/auth?scope=${gmailConfig.scope.join(
    " "
  )}&access_type=${gmailConfig.access_type}&include_granted_scopes=${
    gmailConfig.include_granted_scopes
  }&response_type=${gmailConfig.response_type}&redirect_uri=${
    gmailConfig.redirect_uri
  }&client_id=${gmailConfig.client_id}&login_hint=${address}`;
  // redirect to oauth url
  window.location.href = oauth_url;
  // after redirecting back
  let url = window.location.search;
  let urlParams = new URLSearchParams(url);
  let error = urlParams.get("error");
  let authCode = urlParams.get("code");
  if (error !== null) {
    return { errMsg: error, credentials: {} };
  } else if (authCode === null) {
    return {
      errMsg: `No error in Gmail Oauth2 response but cannot find auth code`,
      credentials: {},
    };
  }
  console.log("Auth code for gmail:", authCode);
  return { errMsg: "", credentials: { auth_code: authCode } };
};

export {
  VerifyResposne,
  verifyOutlook,
  verifyIMAP,
  verifyPOP3,
  verifyGmail,
  getAccessToken,
};
