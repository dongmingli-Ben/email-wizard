import React from "react";
import UserAccountInfo from "./UserAccountInfo";

type SideBarProps = {
  userId: string;
  userSecret: string;
};

const USERNAME = "jake";
const USERACCOUNTS = ["jake@outlook.com", "jake@gmail.com"];

const getUserInfo = (
  userId: string,
  userSecret: string
): [string, string[]] => {
  return [USERNAME, USERACCOUNTS];
};

const SideBar = (props: SideBarProps) => {
  const [userName, userAccounts] = getUserInfo(props.userId, props.userSecret);

  return (
    <div className="sidebar-container">
      <UserAccountInfo userName={userName} userAccounts={userAccounts} />
    </div>
  );
};

export default SideBar;
