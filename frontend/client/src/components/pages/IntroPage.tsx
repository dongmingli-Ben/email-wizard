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
      <div className="page-container">
        <div className="intro-container">
          <div className="intro-title-container">
            Email Wizard -- View Your Emails in Calendar with AI
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
          <div className="button-container">
            <button
              type="button"
              className="u-link u-button"
              onClick={(e) => {
                navigate("/login");
              }}
            >
              Log In
            </button>
            <button
              type="button"
              className="u-link u-button"
              onClick={(e) => {
                navigate("/register");
              }}
            >
              Register
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default IntroPage;
