import React, { useEffect } from "react";
import UserAccountInfo from "./UserAccountInfo";
import { useNavigate } from "@reach/router";

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
  setUserId: (userId: string) => void;
  setUserSecret: (userSecret: string) => void;
};

const USERNAME = "jake";
const USERACCOUNTS = ["jake@outlook.com", "jake@gmail.com"];

const getUserInfoAPI = async (
  userId: string,
  userSecret: string
): Promise<{ userName: string; userAccounts: string[] }> => {
  return {
    userName: USERNAME,
    userAccounts: USERACCOUNTS,
  };
};

const SideBar = (props: SideBarProps) => {
  const navigate = useNavigate();

  useEffect(() => {
    getUserInfoAPI(props.userId, props.userSecret)
      .then(({ userName, userAccounts }) => {
        props.setUserInfo({
          username: userName,
          useraccounts: userAccounts,
        });
      })
      .catch((e) => {
        console.log("fail to fetch user profile:", e);
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
      <div className="logout-btn-container">
        <button
          className="u-button logout-btn u-link"
          type="button"
          onClick={(e) => {
            props.setUserId("");
            props.setUserSecret("");
            navigate("/");
          }}
        >
          Log Out
        </button>
      </div>
    </div>
  );
};

export default SideBar;
export type { userInfoType };
