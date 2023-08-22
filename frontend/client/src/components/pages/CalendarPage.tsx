import React, { useState, useEffect } from "react";
import Feed from "../modules/Feed";
import SideBar, { userInfoType } from "../modules/SideBar";
import AddAccountWindow from "../modules/AddAccountWindow";
import { useNavigate } from "@reach/router";

// to use styles, import the necessary CSS files
import "./CalendarPage.css";
import "../../utility.css";

type CalendarPageProps = {
  userId: number;
  userSecret: string;
  setUserId: (userId: number) => void;
  setUserSecret: (userSecret: string) => void;
  path: string;
};

/**
 * Define the "CalendarPage" component as a function.
 */
const CalendarPage = (props: CalendarPageProps) => {
  const [addAccount, setAddAccount] = useState(false);
  const [userInfo, setUserInfo] = useState<userInfoType>();

  const navigate = useNavigate();
  useEffect(() => {
    if (props.userId <= 0) {
      navigate("/login");
    }
  });

  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div
        className={`
            ${
              addAccount ? "app-inactive-container" : "app-active-container"
            } u-flex
          `}
      >
        <SideBar
          userId={props.userId}
          userSecret={props.userSecret}
          setUserId={props.setUserId}
          setUserSecret={props.setUserSecret}
          userInfo={userInfo}
          setUserInfo={setUserInfo}
          setAddAccount={setAddAccount}
        />
        <Feed userId={props.userId} userSecret={props.userSecret} />
      </div>
      {addAccount ? (
        <AddAccountWindow
          userId={props.userId}
          userSecret={props.userSecret}
          userInfo={userInfo}
          setUserInfo={setUserInfo}
          setAddAccount={setAddAccount}
        />
      ) : (
        <></>
      )}
    </>
  );
};

export default CalendarPage;
