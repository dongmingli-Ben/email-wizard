import React, { useState, useEffect } from "react";
import Feed from "./modules/Feed";
import SideBar, { userInfoType } from "./modules/SideBar";
import AddAccountWindow from "./modules/AddAccountWindow";

// to use styles, import the necessary CSS files
import "./App.css";
import "../utility.css";

const userId = "jake";
const userSecret = "secret";

/**
 * Define the "App" component as a function.
 */
const App = () => {
  const [addAccount, setAddAccount] = useState(false);
  const [userInfo, setUserInfo] = useState<userInfoType>();
  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div className="app-container body">
        <div
          className={`
            ${
              addAccount ? "app-inactive-container" : "app-active-container"
            } u-flex
          `}
        >
          <SideBar
            userId={userId}
            userSecret={userSecret}
            userInfo={userInfo}
            setUserInfo={setUserInfo}
            setAddAccount={setAddAccount}
          />
          <Feed userId={userId} userSecret={userSecret} />
        </div>
        {addAccount ? <AddAccountWindow /> : <></>}
      </div>
    </>
  );
};

export default App;
