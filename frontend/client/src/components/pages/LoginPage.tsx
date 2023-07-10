import React, { useState, useEffect } from "react";
import { useNavigate } from "@reach/router";

// to use styles, import the necessary CSS files
import "./LoginPage.css";
import "../../utility.css";

type LoginPageProps = {
  userId: string;
  userSecret: string;
  setUserId: (userId: string) => void;
  setUserSecret: (userSecret: string) => void;
  path: string;
};

/**
 * Define the "CalendarPage" component as a function.
 */
const LoginPage = (props: LoginPageProps) => {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = () => {};
  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div className="page-container">
        <div className="login-container">
          <h3 className="u-textCenter">Log In</h3>
          <form onSubmit={handleSubmit} className="login-form u-flexColumn">
            <div className="u-form-group u-flexColumn">
              <label htmlFor="User Name" className="u-form-lable">
                User Name
              </label>
              <input
                type="text"
                className="u-input"
                value={username}
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
                required
              />
            </div>
            <div className="u-form-group u-flexColumn">
              <label htmlFor="User Password" className="u-form-lable">
                User Password
              </label>
              <input
                type="text"
                className="u-input"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
                required
              />
            </div>
            <div className="u-form-group u-flex u-flex-justifyCenter">
              <button type="submit" className="u-submit-btn u-link u-button">
                Submit
              </button>
              <button
                type="button"
                className="u-cancel-btn u-link u-button"
                onClick={(e) => {
                  navigate("/");
                }}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
};

export default LoginPage;
