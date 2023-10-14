import React, { useState } from "react";
import { userInfoType } from "./SideBar";
import {
  verifyOutlook,
  verifyIMAP,
  verifyPOP3,
  verifyGmail,
  VerifyResposne,
} from "../../utilities/verifyEmail";

import "./AddAccountWindow.css";
import { backendConfig, post } from "../../utilities/requestUtility";

type AddAccountWindowProps = {
  userId: number;
  userSecret: string;
  userInfo: userInfoType | undefined;
  setUserInfo: (info: userInfoType) => void;
  setAddAccount: (status: boolean) => void;
};

const verifyEmailAccount = async (req): Promise<VerifyResposne> => {
  let resp: VerifyResposne;
  if (req.emailtype === "outlook") {
    resp = await verifyOutlook(req.emailaddress);
  } else if (req.emailtype === "IMAP") {
    resp = await verifyIMAP(req.emailaddress, req.password, req.imapServer);
  } else if (req.emailtype === "POP3") {
    resp = await verifyPOP3(req.emailaddress, req.password, req.pop3Server);
  } else if (req.emailtype === "gmail") {
    resp = await verifyGmail(req.emailaddress);
  } else {
    console.log(`Un-recognized account type: ${req.emailtype}`);
    let errMsg = `Un-recognized account type: ${req.emailtype}`;
    resp = {
      errMsg: errMsg,
      credentials: {},
    };
  }
  console.log(resp);
  return resp;
};

const addEmailAccountDBAPI = async (
  req,
  credentials: { [key: string]: string }
): Promise<string> => {
  let add_req = {
    user_id: req.userId,
    user_secret: req.userSecret,
    type: req.emailtype,
    address: req.emailaddress,
    credentials: credentials,
  };
  let errMsg = await post(backendConfig.add_mailbox, add_req)
    .then((resp) => {
      return "";
    })
    .catch((e) => {
      console.log("caught error when adding mailbox:", e);
      console.log(add_req);
      return "fail to add mailbox.";
    });
  return errMsg;
};

const newEmailAccount = async (
  req
): Promise<{ userInfo: userInfoType; errMsg: string }> => {
  let resp = await verifyEmailAccount(req);
  if (resp.errMsg !== "") {
    return {
      userInfo: { username: "", useraccounts: [] },
      errMsg: resp.errMsg,
    };
  }
  let errMsg = await addEmailAccountDBAPI(req, resp.credentials);
  if (errMsg !== "") {
    return {
      userInfo: { username: "", useraccounts: [] },
      errMsg: errMsg,
    };
  }
  return {
    userInfo: {
      username: "",
      useraccounts: [{ address: req.emailaddress, protocol: req.emailtype }],
    },
    errMsg: "",
  };
};

const AddAccountWindow = (props: AddAccountWindowProps) => {
  const [emailType, setEmailType] = useState("");
  const [emailAddress, setEmailAddress] = useState("");
  const [password, setPassword] = useState("");
  const [imapServer, setIMAPServer] = useState("");
  const [pop3Server, setPOP3Server] = useState("");

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const requirePassword = (emailType: string): boolean => {
    let needPasswordEmails = ["IMAP", "POP3"];
    return needPasswordEmails.includes(emailType);
  };

  const handleSubmit = (e) => {
    setLoading(true);
    e.preventDefault();
    let req = {
      emailtype: emailType,
      emailaddress: emailAddress,
      password: password,
      userId: props.userId,
      userSecret: props.userSecret,
      imapServer: imapServer,
      pop3Server: pop3Server,
    };
    console.log(req);
    newEmailAccount(req)
      .then((resp: { userInfo: userInfoType; errMsg: string }) => {
        setLoading(false);
        if (resp.errMsg === "") {
          console.log("adding new mailbox to user:", resp);
          props.setUserInfo({
            username: props.userInfo ? props.userInfo.username : "No User Name",
            useraccounts: props.userInfo
              ? [
                  ...props.userInfo.useraccounts,
                  resp.userInfo.useraccounts[
                    resp.userInfo.useraccounts.length - 1
                  ],
                ]
              : [
                  resp.userInfo.useraccounts[
                    resp.userInfo.useraccounts.length - 1
                  ],
                ],
          });
          props.setAddAccount(false);
        } else {
          setErrorMsg(resp.errMsg);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  };

  return (
    <div className="add-account-overlay-container">
      <div className="add-account-container">
        <h3 className="u-textCenter">New Email Account</h3>
        <form onSubmit={handleSubmit} className="add-account-form u-flexColumn">
          <div className="u-form-group u-flexColumn">
            <label htmlFor="Email Type" className="u-form-lable">
              Select Mailbox Type
            </label>
            <select
              className="dropdown-cell-container u-input"
              value={emailType}
              onChange={(e) => {
                setEmailType(e.target.value);
              }}
              required
            >
              <option value="">-- Select --</option>
              <option value="outlook">Outlook</option>
              <option value="gmail">Gmail</option>
              <option value="IMAP">IMAP</option>
              <option value="POP3">POP3</option>
            </select>
          </div>
          <div className="u-form-group u-flexColumn">
            <label htmlFor="Email Address" className="u-form-lable">
              Email Address
            </label>
            <input
              type="email"
              className="form-input-container u-input"
              value={emailAddress}
              onChange={(e) => {
                setEmailAddress(e.target.value);
              }}
              required
            />
            {errorMsg === "" ? (
              <></>
            ) : (
              <div className="u-error-msg">{errorMsg}</div>
            )}
          </div>
          {requirePassword(emailType) ? (
            <div className="u-form-group u-flexColumn">
              <label htmlFor="Email Password" className="u-form-lable">
                Password
              </label>
              <input
                type="password"
                className="form-input-container u-input"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
            </div>
          ) : (
            <></>
          )}
          {emailType === "IMAP" ? (
            <div className="u-form-group u-flexColumn">
              <label htmlFor="IMAP server" className="u-form-lable">
                IMAP server
              </label>
              <input
                type="text"
                className="form-input-container u-input"
                value={imapServer}
                onChange={(e) => {
                  setIMAPServer(e.target.value);
                }}
              />
            </div>
          ) : (
            <></>
          )}
          {emailType === "POP3" ? (
            <div className="u-form-group u-flexColumn">
              <label htmlFor="POP3 server" className="u-form-lable">
                POP3 server
              </label>
              <input
                type="text"
                className="form-input-container u-input"
                value={pop3Server}
                onChange={(e) => {
                  setPOP3Server(e.target.value);
                }}
              />
            </div>
          ) : (
            <></>
          )}
          <div className="u-form-group u-flex u-flex-justifyCenter">
            <button
              type="submit"
              className="u-submit-btn u-link u-button"
              disabled={loading}
            >
              {loading ? (
                <div className="u-spin-btn u-flex u-flex-justifyCenter">
                  <img src="./static/refresh.svg" className="u-btn-image" />
                </div>
              ) : (
                "Submit"
              )}
            </button>
            <button
              type="button"
              className="u-cancel-btn u-link u-button"
              onClick={(e) => {
                props.setAddAccount(false);
                console.log("setting addAccount to false");
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddAccountWindow;
