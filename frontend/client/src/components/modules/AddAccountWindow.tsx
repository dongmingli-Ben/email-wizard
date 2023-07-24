import React, { useState } from "react";
import { userInfoType } from "./SideBar";
import {
  verifyOutlook,
  verifyIMAP,
  verifyPOP3,
} from "../../utilities/verifyEmail";

import "./AddAccountWindow.css";
import { backendConfig, post } from "../../utilities/requestUtility";

type AddAccountWindowProps = {
  userId: string;
  userSecret: string;
  userInfo: userInfoType | undefined;
  setUserInfo: (info: userInfoType) => void;
  setAddAccount: (status: boolean) => void;
};

const verifyEmailAccount = async (req): Promise<string> => {
  let errMsg: string;
  if (req.emailtype === "outlook") {
    errMsg = await verifyOutlook(req.emailaddress);
  } else if (req.emailtype === "IMAP") {
    errMsg = await verifyIMAP(req.emailaddress, req.password, req.imapServer);
  } else if (req.emailaddress === "POP3") {
    errMsg = await verifyPOP3(req.emailaddress, req.password, req.pop3Server);
  } else {
    console.log(`Un-recognized account type: ${req.emailtype}`);
    errMsg = `Un-recognized account type: ${req.emailtype}`;
  }
  console.log(errMsg);
  return errMsg;
};

const addEmailAccountDBAPI = async (req): Promise<string> => {
  let add_req: { [key: string]: string };
  add_req = {
    userId: req.userId,
    userSecret: req.userSecret,
    type: req.type,
    address: req.emailaddress,
  };
  if (req.type === "IMAP") {
    add_req.password = req.password;
    add_req.imap_server = req.imapServer;
  } else if (req.type === "POP3") {
    add_req.password = req.password;
    add_req.pop3_server = req.pop3Server;
  }
  let errMsg = await post(backendConfig.add_mailbox, add_req)
    .then((resp) => {
      if ("errMsg" in resp) {
        console.log(resp);
        return resp.errMsg;
      }
      return "";
    })
    .catch((e) => {
      console.log("caught error when adding mailbox:", e);
      return "fail to add mailbox.";
    });
  return errMsg;
};

const newEmailAccount = async (
  req
): Promise<{ userInfo: userInfoType; errMsg: string }> => {
  let errMsg = await verifyEmailAccount(req);
  if (errMsg !== "") {
    return {
      userInfo: { username: "", useraccounts: [] },
      errMsg: errMsg,
    };
  }
  errMsg = await addEmailAccountDBAPI(req);
  if (errMsg !== "") {
    return {
      userInfo: { username: "", useraccounts: [] },
      errMsg: errMsg,
    };
  }
  return {
    userInfo: {
      username: "jake",
      useraccounts: ["jake@outlook.com", "jake@gmail.com", req.emailaddress],
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
