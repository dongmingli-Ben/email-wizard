import React from "react";
import SearchBar from "./SearchBar";

type ToolBarProps = {
  setQuery: (query: string) => void;
};

const ToolBar = (props: ToolBarProps) => {
  return (
    <div className="toolbar-container">
      <SearchBar setQuery={props.setQuery} />
    </div>
  );
};

export default ToolBar;
