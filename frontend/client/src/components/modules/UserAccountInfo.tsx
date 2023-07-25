import React from "react";

import "./UserAccountInfo.css";

type UserAccountInfoProps = {
  userName: string;
  userAccounts: string[];
  setAddAccount: (status: boolean) => void;
};

const UserNameBar = (props: { userName: string }) => {
  return <div className="username-container">{props.userName}</div>;
};

const UserAccountBars = (props: { userAccounts: string[] }) => {
  console.log(props.userAccounts.map((ele) => ele));
  return (
    <div className="useraccounts-container">
      {props.userAccounts.map((account: string, index: number) => {
        return (
          <div className="useraccount-cell-container" key={index}>
            {account}
          </div>
        );
      })}
    </div>
  );
};

const AddAccountButton = (props: {
  setAddAccount: (status: boolean) => void;
}) => {
  return (
    <div
      className="u-link add-account-btn"
      onClick={() => {
        props.setAddAccount(true);
      }}
    >
      Add Account
    </div>
  );
};

const UserAccountInfo = (props: UserAccountInfoProps) => {
  return (
    <div className="userinfo-container">
      <UserNameBar userName={props.userName} />
      <UserAccountBars userAccounts={props.userAccounts} />
      <AddAccountButton setAddAccount={props.setAddAccount} />
    </div>
  );
};

export default UserAccountInfo;
