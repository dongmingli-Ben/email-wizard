import React, { useState, useEffect } from "react";
import Feed from "../modules/Feed";
import SideBar, { userInfoType } from "../modules/SideBar";
import AddAccountWindow from "../modules/AddAccountWindow";
import { useNavigate } from "react-router-dom";
import { createTheme, ThemeProvider } from "@mui/material/styles";

// to use styles, import the necessary CSS files
import "./CalendarPage.css";
import "../../utility.css";

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

type CalendarPageProps = {
  userId: number;
  userSecret: string;
  setUserId: (userId: number) => void;
  setUserSecret: (userSecret: string) => void;
};

/**
 * Define the "CalendarPage" component as a function.
 */
const CalendarPage = (props: CalendarPageProps) => {
  const [addAccount, setAddAccount] = useState(false);
  const [userInfo, setUserInfo] = useState<userInfoType>();

  const navigate = useNavigate();
  useEffect(() => {
    const delay = 100;

    // Introduce a delay using setTimeout
    const timerId = setTimeout(() => {
      // Check user id here
      if (props.userId <= 0 || props.userSecret.length === 0) {
        // Navigate the user to the login page
        console.log(
          `current user id: ${props.userId}, secret: ${props.userSecret}`
        );
        navigate("/login");
      }
    }, delay);

    // Clear the timeout if the component unmounts or the effect is re-executed
    return () => clearTimeout(timerId);
  }, [props.userId, props.userSecret]);

  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <ThemeProvider theme={defaultTheme}>
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
        <Feed
          userId={props.userId}
          userSecret={props.userSecret}
          userInfo={userInfo}
        />
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
    </ThemeProvider>
  );
};

export default CalendarPage;
