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
    return post(backendConfig.events, {
      user_id: userId,
      user_secret: userSecret,
      address: address,
      kwargs: {},
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
  } else {
    throw `un-recognized mailbox type: ${protocol}`;
  }
};

const getEventsAPI = async (
  userId: number,
  userSecret: string,
  query: string
): Promise<EventType[]> => {
  return get(backendConfig.events, {
    user_id: userId,
    user_secret: userSecret,
    query: query,
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
      console.log("parsed events:");
      console.log(events);
      return events;
    })
    .catch((e) => {
      console.log("fail to get user events:", e);
      return [];
    });
};

const Calendar = (props: calendarProps) => {
  const [events, setEvents] = useState<EventType[]>([]);

  useEffect(() => {
    console.log("updating events for:", props.userInfo);
    if (props.userInfo !== undefined) {
      getEventsAPI(props.userId, props.userSecret, props.query)
        .then((_events: EventType[]) => {
          setEvents(_events);
        })
        .then(() => {
          if (props.userInfo !== undefined) {
            updateEvents(props.userId, props.userSecret, props.userInfo).then(
              () => {
                getEventsAPI(props.userId, props.userSecret, props.query).then(
                  (_events: EventType[]) => {
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
      (resp: EventType[]) => {
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
      />
    </div>
  );
};

export default Calendar;
