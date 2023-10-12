# Smart Email Parser

## Goals

1. Json Representation

Given an email, parse the content into json lists of events. Each event is represented as a json object, with the following properties:

```json
{
  "event_type": "notification",
  "summary": "what is the notification about?"
}
```

or

```json
{
  "event_type": "activity",
  "start_time": "2023-01-01T00:00:00Z",
  "end_time": "2023-01-01T00:30:00Z",
  "summary": "what is the activity about?",
  "venue": "where will the activity take place?"
}
```

or

```json
{
  "event_type": "registration",
  "end_time": "2023-01-01T00:00:00Z",
  "summary": "what is the registration about?",
  "venue": "Is there a registration link? If there is one, include it. If not, just leave it empty."
}
```

## Run as a service

The parser can be setup as a microservice with gRPC.

First setup OPENAI_API_KEY:

```bash
export OPENAI_API_KEY=xxx
```

Second (optional) (re)generate the gRPC codes with proto file:

```bash
cd service
bash gen_grpc_code.sh
```

Third run the server:

```bash
bash run_server.sh
```

Fourth (optional) test whether the server is healthy and working:

```bash
bash test_server.sh
```

Upon success, you should see something like this:

```json
GetEmails client received:
{
    "events": [
        {
            "event_type": "registration",
            "end_time": "2023-04-06T12:00:00 Asia/Shanghai",
            "summary": "2023大学杰出毕业生奖提名者自荐材料征集",
            "venue": "https://wj.cuhk.edu.cn/vm/YVgulbu.aspx"
        }
    ]
}
```

## Prompts

```
    You will be given an email. Please try to summarize the email.
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
        "start_time": "2023-01-01T00:00:00Z",
        "end_time": "2023-01-01T00:30:00Z",
        "summary": "what is the activity about?"
    }

    A registration event is a registration to an activity. You should be able to
    see a registration deadline for the activity. Make sure your mark that time down.
    You should parse a registration into the following format:

    {
        "event_type": "registration",
        "end_time": "2023-01-01T00:00:00Z",
        "summary": "what is the registration about?"
    }

    Note that an email may contains multiple events. For example, an invitation to
    some activity may contains an activity event and a registration event at the
    same time. In that case, return a list of events.

    Here is the user's email:

    <email-content>
```

## Performance

Evaluation is pretty robust with current evidence. Seems like the randomness in ChatGPT is not an issue.

Baseline prompt performance:

```
           notification  registration  activity
precision      0.598333      0.254167  0.469136
recall         0.800000      0.270833  0.500000
f1             0.644386      0.260417  0.480247
```

Another 2 evalutations on the same data:

```
           notification  registration  activity
precision      0.606667      0.254167  0.469136
recall         0.800000      0.270833  0.500000
f1             0.651053      0.260417  0.480247

           notification  registration  activity
precision      0.615000      0.295833  0.432099
recall         0.800000      0.312500  0.462963
f1             0.656507      0.302083  0.443210
```

Analysis: most errors are due to incorrect year (since year is not provided to the parser).

### Add Received Dates and User Timezone

```
           notification  registration  activity
precision      0.578571      0.354167  0.711538
recall         0.600000      0.354167  0.698718
f1             0.586364      0.354167  0.701465
```

### OpenAI 06-13 Model Update

Note that the model for evaluation also changes.

```
           notification  registration  activity
precision      0.438406      0.333333  0.760417
recall         0.478261      0.347826  0.770833
f1             0.450311      0.339130  0.757937
```
