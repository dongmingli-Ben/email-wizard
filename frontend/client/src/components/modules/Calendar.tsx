import React from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";

type calendarProps = {
  userId: string;
  userSecret: string;
};

const Calendar = (props: calendarProps) => {
  return (
    <div className="calendar-container">
      <FullCalendar
        plugins={[dayGridPlugin]}
        initialView="dayGridMonth"
        weekends={true}
        events={[
          { title: "event 1", date: "2023-07-01" },
          { title: "event 2", date: "2023-04-02" },
        ]}
      />
    </div>
  );
};

export default Calendar;
