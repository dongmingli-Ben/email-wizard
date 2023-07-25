import React, { useEffect, useState } from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";
import { backendConfig, get } from "../../utilities/requestUtility";

type calendarProps = {
  userId: string;
  userSecret: string;
  query: string;
};

type EventType = {
  title: string;
  date: string;
};

const getEventsAPI = async (
  userId: string,
  userSecret: string
): Promise<EventType[]> => {
  return get(backendConfig.events, {
    userId: userId,
    userSecret: userSecret,
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
    getEventsAPI(props.userId, props.userSecret).then(
      (_events: EventType[]) => {
        setEvents(_events);
      }
    );
  }, []);

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
