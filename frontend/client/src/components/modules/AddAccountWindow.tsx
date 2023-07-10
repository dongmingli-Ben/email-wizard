import React, { useState } from "react";
import { userInfoType } from "./SideBar";

import "./AddAccountWindow.css";

type AddAccountWindowProps = {
  userId: string;
  userSecret: string;
  userInfo: userInfoType | undefined;
  setUserInfo: (info: userInfoType) => void;
  setAddAccount: (status: boolean) => void;
};

const newEmailAccountAPI = async (req): Promise<string> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return req.emailaddress;
};

const AddAccountWindow = (props: AddAccountWindowProps) => {
  const [emailType, setEmailType] = useState("");
  const [emailAddress, setEmailAddress] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);

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
    };
    console.log(req);
    newEmailAccountAPI(req)
      .then((address: string) => {
        props.setUserInfo({
          username: props.userInfo ? props.userInfo.username : "No User Name",
          useraccounts: props.userInfo
            ? [...props.userInfo.useraccounts, address]
            : [address],
        });
        setLoading(false);
        props.setAddAccount(false);
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
