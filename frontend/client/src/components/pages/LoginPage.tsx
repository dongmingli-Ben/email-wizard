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

const authenticateUserPassword = async (
  username: string,
  password: string
): Promise<{ userId: string; userSecret: string }> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return {
    userId: "Jake",
    userSecret: "secret",
  };
};

/**
 * Define the "CalendarPage" component as a function.
 */
const LoginPage = (props: LoginPageProps) => {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setLoading(true);
    authenticateUserPassword(username, password).then(
      (resp: { userId: string; userSecret: string }) => {
        console.log(resp);
        if (resp.userId.length > 0) {
          props.setUserId(resp.userId);
          props.setUserSecret(resp.userSecret);
          navigate("/");
        }
        setLoading(false);
      }
    );
  };
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
                type="password"
                className="u-input"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                }}
                required
              />
            </div>
            <div className="u-form-group u-flex u-flex-justifyCenter">
              <button
                type="submit"
                className="u-submit-btn u-link u-button"
                disabled={loading}
              >
                {loading ? (
                  <div className="u-spin-btn u-flex u-flex-justifyCenter">
                    <img src="./static/refresh.svg" className="u-btn-image" />
                  </div>
                ) : (
                  "Submit"
                )}
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
