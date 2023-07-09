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
  return req.emailAddress;
};

const AddAccountWindow = (props: AddAccountWindowProps) => {
  const [emailType, setEmailType] = useState("");
  const [emailAddress, setEmailAddress] = useState("");
  const [password, setPassword] = useState("");

  const requirePassword = (emailType: string): boolean => {
    let needPasswordEmails = ["IMAP", "POP3"];
    return needPasswordEmails.includes(emailType);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    let req = {
      emailtype: emailType,
      emailaddress: emailAddress,
      password: password,
    };
    newEmailAccountAPI(req)
      .then((address: string) => {
        props.setUserInfo({
          username: props.userInfo ? props.userInfo.username : "No User Name",
          useraccounts: props.userInfo
            ? [...props.userInfo.useraccounts, address]
            : [address],
        });
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
        <form onSubmit={handleSubmit} className="add-account-form">
          <label htmlFor="Email Type">Select Mailbox Type</label>
          <select
            className="dropdown-cell-container"
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
          <label htmlFor="Email Address">Email Address</label>
          <input
            type="email"
            className="form-input-container"
            value={emailAddress}
            onChange={(e) => {
              setEmailAddress(e.target.value);
            }}
          />
          {requirePassword(emailType) ? (
            <>
              <label htmlFor="Email Password">Password</label>
              <input
                type="email"
                className="form-input-container"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
              />
            </>
          ) : (
            <></>
          )}
          <button type="submit">Submit</button>
          <button
            type="button"
            onClick={(e) => {
              props.setAddAccount(false);
              console.log("setting addAccount to false");
            }}
          >
            Cancel
          </button>
        </form>
      </div>
    </div>
  );
};

export default AddAccountWindow;
