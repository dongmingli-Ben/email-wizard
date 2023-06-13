def get_prompt(email: dict, *kwargs) -> str:
    prompt = """You will be given an email. Please try to summarize the email. 
    Make sure you includes the important dates because the user wants to add 
    information to his calendar if the email is an activity invitation or 
    contains a registration ddl.

    Here is more information on the your return format. Most importantly, give
    your response in JSON format. For each email, there are three different
    types of events, e.g. notification, activity, and registration. Below is 
    about how each type of events should be represented in JSON.

    A notification is a piece of information that should not appear on user's 
    calender, because it contains no dates (such as system generated no-reply
    message). You should parse a notification into the following format:

    {
        "event_type": "notification",
        "summary": "what is the notification about?"
    }

    An activity is an invitation to an activity, which CONTAINS the start time
    and the end time of the activity. You should be precise on the time, since
    the activity's time will be displayed to user's calendar. Here is what you
    should return if there is an activity:

    {
        "event_type": "activity",
        "start_time": "2023-01-01T00:00:00",
        "end_time": "2023-01-01T00:30:00",
        "summary": "what is the activity about?",
        "venue": "where will the activity take place?"
    }

    A registration event is a registration to an activity. You should be able to 
    see a registration deadline for the activity. Make sure your mark that time down.
    You should parse a registration into the following format:

    {
        "event_type": "registration",
        "end_time": "2023-01-01T00:00:00",
        "summary": "what is the registration about?",
        "venue": "Is there a registration link? If there is one, include it. If not, just leave it empty."
    }

    Note that an email may contains multiple events. For example, an invitation to 
    some activity may contains an activity event and a registration event at the 
    same time. Therefore, to make the return format uniform, please always return a list of events. 
    Use the following format:

    {
        "events": [
            {...}, // event 1
            {...}, // event 2
            ...,
            {...} // event n
        ]
    }

    Here is the user's email:


    """
    return prompt + email['content']