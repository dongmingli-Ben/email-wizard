import React, { useEffect, useState } from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";

type calendarProps = {
  userId: string;
  userSecret: string;
};

type EventType = {
  title: string;
  date: string;
};

const EVENTS = [
  { title: "event 1", date: "2023-07-01" },
  { title: "event 2", date: "2023-04-02" },
];

const getAPIEvents = async (userId: string, userSecret: string) => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return EVENTS;
};

const Calendar = (props: calendarProps) => {
  const [events, setEvents] = useState<EventType[]>([]);

  useEffect(() => {
    getAPIEvents(props.userId, props.userSecret).then(
      (_events: EventType[]) => {
        setEvents(_events);
      }
    );
  }, []);
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
