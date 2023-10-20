import React, { useState, useEffect } from "react";
import { Routes, Route } from "react-router-dom";
import CalendarPage from "./pages/CalendarPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";

// MSAL imports
import { MsalProvider } from "@azure/msal-react";
import { IPublicClientApplication } from "@azure/msal-browser";

// to use styles, import the necessary CSS files
import "./App.css";
import "../utility.css";
import IntroPage from "./pages/IntroPage";

type AppProps = {
  pca: IPublicClientApplication;
};

/**
 * Define the "App" component as a function.
 */
const App = ({ pca }: AppProps) => {
  const [userId, setUserId] = useState(-1);
  const [userSecret, setUserSecret] = useState("");

  useEffect(() => {
    if (userId > 0 && userSecret.length > 0) {
      console.log("setting userId and userSecret to", userId, userSecret);
      sessionStorage.setItem("userId", userId.toString());
      sessionStorage.setItem("userSecret", userSecret);
    }
  }, [userId, userSecret]);

  useEffect(() => {
    const storedId = sessionStorage.getItem("userId");
    const storedSecret = sessionStorage.getItem("userSecret");
    if (storedId !== null && storedSecret !== null) {
      console.log("loading id and secret from session storage ...");
      setUserId(parseInt(storedId));
      setUserSecret(storedSecret);
      console.log(`current id: ${userId}, secret: ${userSecret}`);
    }
  }, []);

  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <MsalProvider instance={pca}>
      <div className="app-container body">
        <Routes>
          <Route
            path="/calendar"
            element={
              <CalendarPage
                userId={userId}
                userSecret={userSecret}
                setUserId={setUserId}
                setUserSecret={setUserSecret}
              />
            }
          />
          <Route
            path="/login"
            element={
              <LoginPage
                userId={userId}
                userSecret={userSecret}
                setUserId={setUserId}
                setUserSecret={setUserSecret}
              />
            }
          />
          <Route
            path="/login"
            element={
              <LoginPage
                userId={userId}
                userSecret={userSecret}
                setUserId={setUserId}
                setUserSecret={setUserSecret}
              />
            }
          />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/" element={<IntroPage />} />
        </Routes>
      </div>
    </MsalProvider>
  );
};

export default App;
