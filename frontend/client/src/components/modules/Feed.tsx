import React from "react";
import ToolBar from "./ToolBar";
import Calendar from "./Calendar";

import "./Feed.css";

type FeedProps = {
  userId: string;
  userSecret: string;
};

const Feed = (props: FeedProps) => {
  return (
    <div className="feed-container u-relative">
      <div className="feed-inner-container">
        <ToolBar />
        <Calendar userId={props.userId} userSecret={props.userSecret} />
      </div>
    </div>
  );
};

export default Feed;
