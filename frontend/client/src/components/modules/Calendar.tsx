import React, { useEffect, useState } from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";
import { appGet, appPost, backendConfig } from "../../utilities/requestUtility";
import { getAccessToken } from "../../utilities/verifyEmail";
import { userInfoType } from "./SideBar";
import { Box, Link, Tooltip, Typography } from "@mui/material";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";
import CelebrationIcon from "@mui/icons-material/Celebration";
import HowToRegIcon from "@mui/icons-material/HowToReg";
import ScheduleIcon from "@mui/icons-material/Schedule";
import PlaceIcon from "@mui/icons-material/Place";

type calendarProps = {
  userId: number;
  userSecret: string;
  query: string;
  setQuery: (query: string) => void;
  userInfo: userInfoType | undefined;
};

const updateEvents = async (
  userId: number,
  userSecret: string,
  userInfo: userInfoType
): Promise<void> => {
  for (const mailbox of userInfo.useraccounts) {
    try {
      await updateAccountEventsAPI(
        userId,
        userSecret,
        mailbox.address,
        mailbox.protocol
      );
    } catch (error) {
      console.log(error);
      console.log(`fail to update events for mailbox: ${mailbox}`);
    }
  }
};

const updateAccountEventsAPI = async (
  userId: number,
  userSecret: string,
  address: string,
  protocol: string
): Promise<void> => {
  if (protocol === "IMAP" || protocol == "POP3" || protocol == "gmail") {
    return appPost(backendConfig.events, userId, userSecret, {
      address: address,
      kwargs: {},
    });
  } else if (protocol == "outlook") {
    let access_token = await getAccessToken(address);
    if (access_token.length === 0) {
      console.log("fail to get access token, got: ", access_token);
      return;
    }
    return appPost(backendConfig.events, userId, userSecret, {
      address: address,
      kwargs: {
        auth_token: access_token,
      },
    });
  } else {
    throw `un-recognized mailbox type: ${protocol}`;
  }
};

const getEventsAPI = async (
  userId: number,
  userSecret: string,
  query: string
): Promise<{ [key: string]: any }[]> => {
  return appGet(backendConfig.events, userId, userSecret, {
    query: query,
  })
    .then((resp) => {
      console.log(`events returned`);
      console.log(resp);
      let events: { [key: string]: any }[] = [];
      for (const e of resp) {
        if ("end_time" in e && e.end_time != "unspecified") {
          let startTime = "start_time" in e ? e.start_time : e.end_time;
          events = [
            ...events,
            {
              title: e.summary,
              start: startTime.split(" ")[0],
              end: e.end_time.split(" ")[0],
              extendedProps: {
                event: e,
              },
            },
          ];
        }
      }
      console.log("parsed events:");
      console.log(events);
      return events;
    })
    .catch((e) => {
      console.log("fail to get user events:", e);
      return [];
    });
};

const EventPopupDisplay = ({ event }: { event: { [key: string]: string } }) => {
  const getLocalTime = (time: string) => {
    let localtime = time.split(" ")[0].split("T")[1];
    return localtime.split(":").slice(0, 2).join(":");
  };
  console.log(event);
  return (
    <Box>
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
        }}
      >
        <Box pr={1}>
          {event.event_type === "registration" ? (
            <HowToRegIcon></HowToRegIcon>
          ) : (
            <CelebrationIcon></CelebrationIcon>
          )}
        </Box>
        <Typography>{event.summary}</Typography>
      </Box>

      <Box
        sx={{
          display: "flex",
          alignItems: "center",
        }}
      >
        <Box pr={1}>
          <ScheduleIcon></ScheduleIcon>
        </Box>
        {"start_time" in event ? (
          <Typography>
            {getLocalTime(event.start_time)} - {getLocalTime(event.end_time)}
          </Typography>
        ) : (
          <Typography>{getLocalTime(event.end_time)}</Typography>
        )}
      </Box>
      {event.venue === "" || event.venue === "unspecified" ? (
        <></>
      ) : (
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
          }}
        >
          <Box pr={1}>
            <PlaceIcon></PlaceIcon>
          </Box>
          <Typography>
            {event.venue.startsWith("http") ? (
              <Link href={event.venue} color="inherit">
                URL Link
              </Link>
            ) : (
              event.venue
            )}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

const CustomEvent = ({ event }) => {
  const e = event.extendedProps.event;
  console.log(e);
  return (
    <Tooltip
      title={<EventPopupDisplay event={e}></EventPopupDisplay>}
      placement="right"
      sx={{
        width: "inherit",
      }}
    >
      <Box
        sx={{
          fontSize: "inherit",
          color: "primary.main",
          display: "flex",
          alignItems: "center",
          width: "100%",
          boxSizing: "border-box",
          "&:hover": {
            cursor: "default",
          },
        }}
      >
        <FiberManualRecordIcon fontSize="inherit"></FiberManualRecordIcon>
        <Typography
          fontSize="inherit"
          noWrap
          sx={{
            width: "100%",
          }}
        >
          {e.summary}
        </Typography>
      </Box>
    </Tooltip>
  );
};

const Calendar = (props: calendarProps) => {
  const [events, setEvents] = useState<{ [key: string]: any }[]>([]);

  useEffect(() => {
    console.log("updating events for:", props.userInfo);
    if (props.userInfo !== undefined) {
      getEventsAPI(props.userId, props.userSecret, props.query)
        .then((_events: { [key: string]: any }[]) => {
          setEvents(_events);
        })
        .then(() => {
          if (props.userInfo !== undefined) {
            updateEvents(props.userId, props.userSecret, props.userInfo).then(
              () => {
                getEventsAPI(props.userId, props.userSecret, props.query).then(
                  (_events: { [key: string]: any }[]) => {
                    setEvents(_events);
                  }
                );
              }
            );
          }
        });
    }
  }, [props.userInfo]);

  useEffect(() => {
    if (props.query.length > 0) {
      alert(
        "Elastic search is temperarily disabled due to limited resources. Please try later!"
      );
    }
    return;
    getEventsAPI(props.userId, props.userSecret, props.query).then(
      (resp: { [key: string]: string }[]) => {
        console.log("query result:");
        console.log(resp);
        setEvents(resp);
      }
    );
  }, [props.query]);

  return (
    <div className="calendar-container u-block">
      <FullCalendar
        plugins={[dayGridPlugin]}
        initialView="dayGridMonth"
        weekends={true}
        events={events}
        headerToolbar={{
          left: "title",
          center: "",
          right: "today prev,next",
        }}
        eventBackgroundColor="white"
        // eventDidMount={(info) => {
        //   const tooltip = (
        //     <Tooltip title={info.event.title}>
        //       <div dangerouslySetInnerHTML={{ __html: info.el.outerHTML }} />
        //     </Tooltip>
        //   );
        //   info.el.innerHTML = '';
        //   info.el.appendChild(tooltip);
        //   console.log(info);
        // }}
        eventContent={(arg) => <CustomEvent event={arg.event} />}
      />
    </div>
  );
};

export default Calendar;
