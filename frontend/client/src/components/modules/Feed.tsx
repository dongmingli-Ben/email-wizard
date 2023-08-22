import React, { useState } from "react";
import ToolBar from "./ToolBar";
import Calendar from "./Calendar";

import "./Feed.css";

type FeedProps = {
  userId: number;
  userSecret: string;
};

const Feed = (props: FeedProps) => {
  const [query, setQuery] = useState("");
  return (
    <div className="feed-container u-relative">
      <div className="feed-inner-container">
        <ToolBar setQuery={setQuery} />
        <Calendar
          userId={props.userId}
          userSecret={props.userSecret}
          query={query}
        />
      </div>
    </div>
  );
};

export default Feed;
