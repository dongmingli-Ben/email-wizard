import React from "react";

import "./UserAccountInfo.css";

type UserAccountInfoProps = {
  userName: string;
  userAccounts: string[];
};

const UserNameBar = (props: { userName: string }) => {
  return <div className="username-container">{props.userName}</div>;
};

const UserAccountBars = (props: { userAccounts: string[] }) => {
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

const AddAccountButton = () => {
  return <div className="u-link add-account-btn">Add Account</div>;
};

const UserAccountInfo = (props: UserAccountInfoProps) => {
  return (
    <div className="userinfo-container">
      <UserNameBar userName={props.userName} />
      <UserAccountBars userAccounts={props.userAccounts} />
      <AddAccountButton />
    </div>
  );
};

export default UserAccountInfo;
