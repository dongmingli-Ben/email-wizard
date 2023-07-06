import React from "react";
import ToolBar from "./ToolBar";
import Calendar from "./Calendar";

type FeedProps = {
  userId: string;
  userSecret: string;
};

const Feed = (props: FeedProps) => {
  return (
    <div className="feed-container">
      <ToolBar />
      <Calendar />
    </div>
  );
};

export default Feed;
