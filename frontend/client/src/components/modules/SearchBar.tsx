import React from "react";
import { useState } from "react";

import "./SearchBar.css";

type SearchBarProps = {
  setQuery: (query: string) => void;
};

const SearchBar = (props: SearchBarProps) => {
  const [inputQuery, setInputQuery] = useState("");
  return (
    <div className="search-container u-flex u-flex-justifyCenter">
      <div className="search-inner-container u-flex ">
        <div className="search-input-container">
          <input
            type="text"
            value={inputQuery}
            className="search-input-text-container"
            onChange={(event) => {
              setInputQuery(event.target.value);
            }}
          />
        </div>
        <div className="search-btn-container">
          <button
            onClick={() => {
              console.log(inputQuery);
              props.setQuery(inputQuery);
            }}
            value=""
            type="submit"
            className="submit-btn u-pointer"
          ></button>
        </div>
      </div>
    </div>
  );
};

export default SearchBar;
