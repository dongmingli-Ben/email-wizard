import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "@reach/router";

// to use styles, import the necessary CSS files
import "./LoginPage.css";
import "../../utility.css";
import { backendConfig, get } from "../../utilities/requestUtility";

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
): Promise<{ userId: string; userSecret: string; errMsg: string }> => {
  return get(backendConfig.verify_user, {
    username: username,
    password: password,
  })
    .then((resp) => {
      return {
        userId: resp.user_id,
        userSecret: resp.user_secret,
        errMsg: "",
      };
    })
    .catch((e) => {
      console.log(`error in user verification: ${e}`);
      return {
        userId: "",
        userSecret: "",
        errMsg: "Cannot verify your user name and password",
      };
    });
};

/**
 * Define the "CalendarPage" component as a function.
 */
const LoginPage = (props: LoginPageProps) => {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg("");
    authenticateUserPassword(username, password).then(
      (resp: { userId: string; userSecret: string; errMsg: string }) => {
        console.log(resp);
        if (resp.userId.length > 0) {
          props.setUserId(resp.userId);
          props.setUserSecret(resp.userSecret);
          navigate("/calendar");
        } else {
          setErrorMsg(resp.errMsg);
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
              {errorMsg === "" ? (
                <></>
              ) : (
                <div className="u-error-msg">{errorMsg}</div>
              )}
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
                  "Log In"
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
          <div className="u-textCenter u-prompt">
            Do not have an account? <Link to="/register">Register</Link> Now!
          </div>
        </div>
      </div>
    </>
  );
};

export default LoginPage;
