import React from "react";
import { useState } from "react";

import "./SearchBar.css";

const SearchBar = () => {
  const [searchQuery, setSearchQuery] = useState("");
  return (
    <div className="search-container u-flex u-flex-justifyCenter">
      <div className="search-inner-container u-flex ">
        <div className="search-input-container">
          <input
            type="text"
            value={searchQuery}
            className="search-input-text-container"
            onChange={(event) => {
              setSearchQuery(event.target.value);
            }}
          />
        </div>
        <div className="search-btn-container">
          <button
            onClick={() => {
              console.log(searchQuery);
              // props.setQuery(inputQuery);
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
