import React from "react";
import ReactDOM from "react-dom";

// renders React Component "Root" into the DOM element with ID "root"
ReactDOM.render(<h1>This is a test. If you see this, it is working</h1>, document.getElementById("root"));

// allows for live updating
module.hot.accept();