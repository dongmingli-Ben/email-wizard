import React, { useEffect } from "react";
import UserAccountInfo from "./UserAccountInfo";

import "./SideBar.css";

type userInfoType = {
  username: string;
  useraccounts: string[];
};

type SideBarProps = {
  userId: string;
  userSecret: string;
  userInfo: userInfoType | undefined;
  setUserInfo: (info: userInfoType) => void;
  setAddAccount: (status: boolean) => void;
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
  useEffect(() => {
    const [userName, userAccounts] = getUserInfo(
      props.userId,
      props.userSecret
    );
    props.setUserInfo({
      username: userName,
      useraccounts: userAccounts,
    });
  }, []);

  return (
    <div className="sidebar-container">
      <UserAccountInfo
        userName={props.userInfo ? props.userInfo.username : ""}
        userAccounts={
          props.userInfo ? props.userInfo.useraccounts : ["No accounts"]
        }
        setAddAccount={props.setAddAccount}
      />
    </div>
  );
};

export default SideBar;
export type { userInfoType };
