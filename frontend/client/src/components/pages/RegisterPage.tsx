import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";

// to use styles, import the necessary CSS files
import "./RegisterPage.css";
import "../../utility.css";
import { backendConfig, post } from "../../utilities/requestUtility";

const registerUserPassword = async (
  username: string,
  password: string
): Promise<{ errMsg: string }> => {
  return post(backendConfig.add_user, {
    username: username,
    password: password,
  })
    .then((resp) => {
      return { errMsg: "" };
    })
    .catch((e) => {
      console.log("fail to add new user:", e);
      return {
        errMsg: "fail to add: Please change your user name or password",
      };
    });
};

/**
 * Define the "CalendarPage" component as a function.
 */
const RegisterPage = (props: {}) => {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [registerSuccess, setRegisterSuccess] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setLoading(true);
    setErrorMsg("");
    registerUserPassword(username, password).then(
      (resp: { errMsg: string }) => {
        console.log(resp);
        if (resp.errMsg.length > 0) {
          setErrorMsg(resp.errMsg);
        } else {
          setRegisterSuccess(true);
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
        <div className="register-container">
          {registerSuccess ? (
            <div>
              You have successfully created a new account.{" "}
              <Link to="/login">Sign in</Link> Now!
            </div>
          ) : (
            <>
              <h3 className="u-textCenter">Register</h3>
              <form
                onSubmit={handleSubmit}
                className="register-form u-flexColumn"
              >
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
                        <img
                          src="./static/refresh.svg"
                          className="u-btn-image"
                        />
                      </div>
                    ) : (
                      "Register"
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
            </>
          )}
          <div className="u-textCenter  u-prompt">
            Already have an account? <Link to="/login">Log In</Link> Now!
          </div>
        </div>
      </div>
    </>
  );
};

export default RegisterPage;
