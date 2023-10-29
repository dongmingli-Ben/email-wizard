import React, { useState } from "react";
import ToolBar from "./ToolBar";
import Calendar from "./Calendar";

import "./Feed.css";
import { userInfoType } from "./SideBar";

type FeedProps = {
  userId: number;
  userSecret: string;
  userInfo: userInfoType | undefined;
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
          setQuery={setQuery}
          userInfo={props.userInfo}
        />
      </div>
    </div>
  );
};

export default Feed;
