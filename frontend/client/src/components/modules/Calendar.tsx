import React, { useEffect, useState } from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";

type calendarProps = {
  userId: string;
  userSecret: string;
  query: string;
};

type EventType = {
  title: string;
  date: string;
};

const EVENTS = [
  { title: "event 1", date: "2023-07-01" },
  { title: "event 2", date: "2023-04-02" },
];

const getEventsAPI = async (userId: string, userSecret: string) => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return EVENTS;
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
