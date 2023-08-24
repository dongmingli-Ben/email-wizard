import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "@reach/router";

// to use styles, import the necessary CSS files
import "./IntroPage.css";
import "../../utility.css";

type IntroPageProps = {
  path: string;
};

/**
 * Define the "IntroPage" component as a function.
 */
const IntroPage = (props: IntroPageProps) => {
  const navigate = useNavigate();

  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div className="intro-page-container">
        <div className="intro-container">
          <div className="intro-title-container u-flex u-flex-alignCenter">
            <span>Email Wizard: View Your Emails in Calendar with AI</span>
            <a href="https://github.com/dongmingli-Ben/email-wizard">
              <img
                src="./static/github-mark.svg"
                alt="GitHub Logo"
                className="github-logo"
              />
            </a>
          </div>
          <div className="intro-feature-container">
            Streamline event management with Email Wizard. Our web app reads
            emails, extracts key details, and populates your calendar
            automatically. Say goodbye to manual entry and inbox clutter. With
            Email Wizard, invitations and activities are effortlessly organized.
            Enjoy seamless integration, as dates, times, locations, and
            participants are recognized instantly. Simplify event planning and
            unlock efficiency with Email Wizard.
          </div>
          <div className="intro-feature-container intro-message-container">
            Preview! Still in progress.
          </div>
          <div className="button-container u-flex u-flex-justifyCenter">
            <button
              type="button"
              className="u-link u-button intro-btn"
              onClick={(e) => {
                navigate("/login");
              }}
            >
              Log In
            </button>
            <button
              type="button"
              className="u-link u-button intro-btn"
              onClick={(e) => {
                navigate("/register");
              }}
            >
              Register
            </button>
          </div>
        </div>
        <div className="intro-image-container u-flex u-flex-justifyCenter">
          <img src="./static/intro.png" className="intro-image" />
        </div>
      </div>
    </>
  );
};

export default IntroPage;
