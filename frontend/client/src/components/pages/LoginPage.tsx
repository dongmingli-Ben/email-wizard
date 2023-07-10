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
          <h3 className="u-textCenter">New Account</h3>
          <form onSubmit={handleSubmit} className="login-form u-flexColumn">
            <div className="form-group u-flexColumn">
              <label htmlFor="User Name" className="form-lable">
                User Name
              </label>
              <input
                type="text"
                className="form-input-container"
                value={username}
                onChange={(e) => {
                  setUsername(e.target.value);
                }}
                required
              />
            </div>
            <div className="form-group u-flexColumn">
              <label htmlFor="User Password" className="form-lable">
                User Password
              </label>
              <input
                type="text"
                className="form-input-container"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
                required
              />
            </div>
            <div className="form-group u-flex u-flex-justifyCenter">
              <button type="submit" className="login-submit-btn u-link">
                Submit
              </button>
              <button
                type="button"
                className="login-cancel-btn u-link"
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
