import React, { useState, useEffect } from "react";
import Feed from "./modules/Feed";
import SideBar from "./modules/SideBar";

// to use styles, import the necessary CSS files
import "./App.css";
import "../utility.css";

const userId = "jake";
const userSecret = "secret";

/**
 * Define the "App" component as a function.
 */
const App = () => {
  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div className="app-container body u-flex">
        <SideBar userId={userId} userSecret={userSecret} />
        <Feed userId={userId} userSecret={userSecret} />
      </div>
    </>
  );
};

export default App;
