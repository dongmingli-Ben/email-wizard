import React, { useState } from "react";
import ToolBar from "./ToolBar";
import Calendar from "./Calendar";

import "./Feed.css";
import { userInfoType } from "./SideBar";
import { Box } from "@mui/material";

type FeedProps = {
  userId: number;
  userSecret: string;
  userInfo: userInfoType | undefined;
};

const Feed = (props: FeedProps) => {
  const [query, setQuery] = useState("");
  return (
    <Box
      sx={{
        width: "80%",
        height: "100%",
        boxSizing: "border-box",
        overflow: "auto",
      }}
      m={0}
      p="3%"
    >
      <ToolBar setQuery={setQuery} />
      <Calendar
        userId={props.userId}
        userSecret={props.userSecret}
        query={query}
        setQuery={setQuery}
        userInfo={props.userInfo}
      />
    </Box>
  );
};

export default Feed;
