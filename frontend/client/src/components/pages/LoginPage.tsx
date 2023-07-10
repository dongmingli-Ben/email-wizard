import React, { useState, useEffect } from "react";

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
  return (
    // <> is like a <div>, but won't show
    // up in the DOM tree
    <>
      <div>Log in page</div>
    </>
  );
};

export default LoginPage;
