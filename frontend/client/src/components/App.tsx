import React, { useState, useEffect } from "react";
import { Router } from "@reach/router";
import CalendarPage from "./pages/CalendarPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";

// to use styles, import the necessary CSS files
import "./App.css";
import "../utility.css";
import IntroPage from "./pages/IntroPage";

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
          <CalendarPage
            path="/calendar"
            userId={userId}
            userSecret={userSecret}
            setUserId={setUserId}
            setUserSecret={setUserSecret}
          />
          <LoginPage
            path="/login"
            userId={userId}
            userSecret={userSecret}
            setUserId={setUserId}
            setUserSecret={setUserSecret}
          />
          <RegisterPage path="/register" />
          <IntroPage path="/" />
        </Router>
      </div>
    </>
  );
};

export default App;
