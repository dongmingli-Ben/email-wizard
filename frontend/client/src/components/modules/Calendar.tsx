import React, { useEffect, useState } from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";
import { backendConfig, get, post } from "../../utilities/requestUtility";
import { getAccessToken } from "../../utilities/verifyEmail";
import { userInfoType } from "./SideBar";

type calendarProps = {
  userId: number;
  userSecret: string;
  query: string;
  userInfo: userInfoType | undefined;
};

type EventType = {
  title: string;
  date: string;
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
      throw error;
    }
  }
};

const updateAccountEventsAPI = async (
  userId: number,
  userSecret: string,
  address: string,
  protocol: string
): Promise<void> => {
  if (protocol === "IMAP" || protocol == "POP3") {
    return post(backendConfig.events, {
      user_id: userId,
      user_secret: userSecret,
      address: address,
    });
  } else if (protocol == "outlook") {
    let access_token = await getAccessToken(address);
    if (access_token.length === 0) {
      console.log("fail to get access token, got: ", access_token);
      return;
    }
    return post(backendConfig.events, {
      user_id: userId,
      user_secret: userSecret,
      address: address,
      kwargs: {
        auth_token: access_token,
      },
    });
  }
};

const getEventsAPI = async (
  userId: number,
  userSecret: string
): Promise<EventType[]> => {
  return get(backendConfig.events, {
    user_id: userId,
    user_secret: userSecret,
  })
    .then((resp) => {
      console.log(`events returned`);
      console.log(resp);
      let events: EventType[] = [];
      for (const e of resp) {
        if ("end_time" in e) {
          events = [
            ...events,
            {
              title: e.summary,
              date: e.end_time.split("T")[0],
            },
          ];
        }
      }
      console.log(events);
      return events;
    })
    .catch((e) => {
      console.log("fail to get user events:", e);
      return [];
    });
};

const match = (event: EventType, query: string): boolean => {
  // use a naive match for now
  let event_words = event.title.toLowerCase().split(" ");
  let matched_words = event_words.filter(
    (word) => query.toLowerCase().search(word) !== -1
  );
  return matched_words.length > 0;
};

const search = (events: EventType[], query: string): EventType[] => {
  console.log("query:", query);
  if (query === "") {
    return events;
  }
  let matched_events: EventType[] = [];
  matched_events = events.filter((event) => match(event, query));
  return matched_events;
};

const Calendar = (props: calendarProps) => {
  const [events, setEvents] = useState<EventType[]>([]);
  const [displayEvents, setDisplayEvents] = useState<EventType[]>([]);

  useEffect(() => {
    console.log("updating events for:", props.userInfo);
    if (props.userInfo !== undefined) {
      updateEvents(props.userId, props.userSecret, props.userInfo).then(() => {
        getEventsAPI(props.userId, props.userSecret).then(
          (_events: EventType[]) => {
            setEvents(_events);
          }
        );
      });
    }
  }, [props.userInfo]);

  useEffect(() => {
    let _events = search(events, props.query);
    setDisplayEvents(_events);
  }, [events, props.query]);

  return (
    <div className="calendar-container u-block">
      <FullCalendar
        plugins={[dayGridPlugin]}
        initialView="dayGridMonth"
        weekends={true}
        events={displayEvents}
      />
    </div>
  );
};

export default Calendar;
