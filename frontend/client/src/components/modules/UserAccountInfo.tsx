import React from "react";

import "./UserAccountInfo.css";

type UserAccountInfoProps = {
  userName: string;
  userAccounts: { address: string; protocol: string }[];
  setAddAccount: (status: boolean) => void;
};

const UserNameBar = (props: { userName: string }) => {
  return <div className="username-container">{props.userName}</div>;
};

const UserAccountBars = (props: {
  userAccounts: { address: string; protocol: string }[];
}) => {
  console.log(props.userAccounts.map((ele) => ele));
  return (
    <div className="useraccounts-container">
      {props.userAccounts.map(
        (account: { address: string; protocol: string }, index: number) => {
          return (
            <div className="useraccount-cell-container" key={index}>
              {account.address}
            </div>
          );
        }
      )}
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
