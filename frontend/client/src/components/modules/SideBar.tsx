import React, { useEffect } from "react";
import UserAccountInfo from "./UserAccountInfo";
import { useNavigate } from "@reach/router";

import "./SideBar.css";
import { get } from "../../utilities/requestUtility";
import { backendConfig } from "../../utilities/requestUtility";

type userInfoType = {
  username: string;
  useraccounts: string[];
};

type SideBarProps = {
  userId: number;
  userSecret: string;
  userInfo: userInfoType | undefined;
  setUserInfo: (info: userInfoType) => void;
  setAddAccount: (status: boolean) => void;
  setUserId: (userId: number) => void;
  setUserSecret: (userSecret: string) => void;
};

const getUserInfoAPI = async (
  userId: number,
  userSecret: string
): Promise<{ userName: string; userAccounts: string[]; errMsg: string }> => {
  return get(backendConfig.user_profile, {
    userId: userId,
    userSecret: userSecret,
  })
    .then((resp) => {
      console.log("mailboxes: ", resp.mailboxes);
      let mailboxes = resp.mailboxes.length > 0 ? resp.mailboxes : [];
      return {
        userName: resp.user_name,
        userAccounts: mailboxes.map((ele) => ele.username),
        errMsg: "",
      };
    })
    .catch((e) => {
      console.log("fail to get user profile:", e);
      return {
        userName: "",
        userAccounts: [],
        errMsg: "fail to get user profile",
      };
    });
};

const SideBar = (props: SideBarProps) => {
  const navigate = useNavigate();

  useEffect(() => {
    getUserInfoAPI(props.userId, props.userSecret)
      .then(({ userName, userAccounts, errMsg }) => {
        console.log(userName);
        console.log(userAccounts);
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
          props.userInfo && props.userInfo.useraccounts.length > 0
            ? props.userInfo.useraccounts
            : ["No accounts"]
        }
        setAddAccount={props.setAddAccount}
      />
      <div className="logout-btn-container">
        <button
          className="u-button logout-btn u-link"
          type="button"
          onClick={(e) => {
            props.setUserId(-1);
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
