import React, { useState, useEffect } from "react";
import { Router } from "@reach/router";
import CalendarPage from "./pages/CalendarPage";
import LoginPage from "./pages/LoginPage";

// to use styles, import the necessary CSS files
import "./App.css";
import "../utility.css";

/**
 * Define the "App" component as a function.
 */
const App = () => {
  const [userId, setUserId] = useState<string>("");
  const [userSecret, setUserSecret] = useState("");
  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div className="app-container body">
        <Router>
          <CalendarPage path="/" userId={userId} userSecret={userSecret} />
          <LoginPage
            path="/login"
            userId={userId}
            userSecret={userSecret}
            setUserId={setUserId}
            setUserSecret={setUserSecret}
          />
        </Router>
      </div>
    </>
  );
};

export default App;
