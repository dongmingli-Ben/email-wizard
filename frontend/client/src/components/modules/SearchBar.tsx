import React from "react";
import { useState } from "react";

const SearchBar = () => {
  const [searchQuery, setSearchQuery] = useState("");
  return (
    <div className="search-container">
      <div className="search-input-container u-inlineBlock">
        <input
          type="text"
          value={searchQuery}
          onChange={(event) => {
            setSearchQuery(event.target.value);
          }}
        />
      </div>
      <div className="search-btn-container u-inlineBlock">
        <button
          onClick={() => {
            console.log(searchQuery);
            // props.setQuery(inputQuery);
          }}
          value=""
          type="submit"
          className="submit-btn"
        ></button>
      </div>
    </div>
  );
};

export default SearchBar;
