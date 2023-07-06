import React from "react";

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
      {props.userAccounts.map((account: string) => {
        return <div className="useraccount-cell-container">{account}</div>;
      })}
    </div>
  );
};

const AddAccountButton = () => {
  return <div className="u-button">Add Account</div>;
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
