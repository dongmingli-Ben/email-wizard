# API Definitions

## `email` service

### GetEmails

Request:

```
config: json string
n_mails: int
```

Response:

```
message: json string
```

#### Example

Request:

config:

```json
{
    "username": "xxx@outlook.com",
    "protocol": "outlook"
}
```

n_mails:

```
2
```

Response:

```json
{
    "items": [
        {
            "email_id": "AQMkADAwATNiZmYAZC01YTBlLTVlMjMtMDACLTAwCgBGAAADWf5YPRULfUyDmeRFNdjmKgcAnOT21DDXNEqBQujyKJX_pAAAAgEMAAAAnOT21DDXNEqBQujyKJX_pAAAAHzUVt8AAAA=",
            "item": {
                "subject": "API Updates: GPT-3.5 Turbo, Function Calling, Longer Context, Lower Prices",
                "sender": "noreply@email.openai.com",
                "date": "2023-06-13T23:13:32Z",
                "recipient": [
                    "xxx@outlook.com"
                ],
                "content": "More capable, more steerable OpenAI models We're excited to announce a few updates to the OpenAI developer platform. GPT-3.5 Turbo This model has been updated with a new version: gpt-3.5-turbo-0613 which is more steerable with the system message and includes a new capability: function calling. By describing functions in your prompts, the model can intelligently output a JSON object containing arguments to call these functions based on user input — perfect for integrating with other tools or APIs. Learn more in our function calling documentation . Plus enjoy a 25% cost reduction for input tokens on GPT-3.5 Turbo (now $0.0015 per 1K input tokens), effective immediately. Longer Context We're also introducing gpt-3.5-turbo-16k . This model offers four times the context length of the 4k base model and is priced at $0.003 per 1K input tokens and $0.004 per 1K output tokens. Model Transitioning You can begin using the new gpt-3.5-turbo-0613 model today. On June 27th, the stable gpt-3.5-turbo will be automatically upgraded to this new version. If you need more time to transition, you can specify gpt-3.5-turbo-0301 to keep using the older version, which will remain available until September 13th as part of our upgrade and deprecation process . —The OpenAI team OpenAI 548 Market Street, PMB 97273 , San Francisco , CA 94104-5401 Unsubscribe - Unsubscribe Preferences"
            }
        },
        {
            "email_id": "AQMkADAwATNiZmYAZC01YTBlLTVlMjMtMDACLTAwCgBGAAADWf5YPRULfUyDmeRFNdjmKgcAnOT21DDXNEqBQujyKJX_pAAAAgEMAAAAnOT21DDXNEqBQujyKJX_pAAAAHwqDUwAAAA=",
            "item": {
                "subject": "Microsoft account unusual sign-in activity",
                "sender": "account-security-noreply@accountprotection.microsoft.com",
                "date": "2023-06-12T10:10:00Z",
                "recipient": [
                    "xxx@outlook.com"
                ],
                "content": "Microsoft account Unusual sign-in activity We detected something unusual about a recent sign-in to the Microsoft account gu**4@outlook.com . Sign-in details Country/region: China IP address: 218.81.240.163 Date: 6/12/2023 10:10 AM (GMT) Platform: Android Browser: Android Please go to your recent activity page to let us know whether or not this was you. If this wasn't you, we'll help you secure your account. If this was you, we'll trust similar activity in the future. Review recent activity To opt out or change where you receive security notifications, click here . Thanks, The Microsoft account team"
            }
        }
    ]
}
```

## `parse` service

### ParseEmail

Request:

```
email: json string
additional_info: json string
```

Response:

```
message: json string
```

### Example

Request:

email:

```json
{
    "subject": "重要通知 | 2023大学杰出毕业生奖提名者自荐材料征集",
    "sender": "someone@example.com",
    "date": "2023-06-07T13:58:41Z",
    "recipient": ["whoami@outlook.com"],
    "content": "亲爱的 同学： 你好！ 经由学院、书院推荐，你现已进入今年“大学杰出毕业生奖（ Presidential Award for Outstanding Students ）”的提名名单。为帮助学院、书院及遴选委员会更全面、深入了解你的情况，我们诚挚地邀请你提供自荐信以及证明等资料作为补充资料供相关部门及委员会 参 考。 提交方式： 请点击下方链接或扫描二维码，提交相关材料。 链接： https://wj.cuhk.edu.cn/vm/YVgulbu.aspx 二维码： 截止日期： 4 月 6 日（星期四）中午 12:00 * 注意 * ： 1. 逾期未提交或未按要求提交的同学，将视作放弃自荐机会，请知悉。 2. 如不接受提名，请 在 4 月 5 日 之前回复此邮件 告知，我们将通知学院、书院尽快提供后补人选。 3 ．如果你已提交，则无需再次提交。 奖项介绍 大学杰出毕业生奖 （ Presidential Award for Outstanding Students ）是大学毕业生的最高荣誉，旨在表彰在大学期间具有杰出的学术表现、卓越的领导力以及对大学发展有贡献的同学，今年提名人数（含学院及书院提名）约为毕业生总人数的 3.5% ，获奖人数约为毕业生总人数的 1% 。 关于该奖项的评选，如有疑问，可联系学生事务处杨老师 yangxiaoyu@cuhk.edu.cn 。 顺祝， 平安顺遂 Best Regards, Office of Student Affairs The Chinese University of Hong Kong, Shenzhen Hotline: (86) 0755 — 8427 3671 E-mail: osa@cuhk.edu.cn Office Hours: 8:30-12:00 13:00-17:30 (Mon-Fri) 学生事务处 香港中文大学（深圳） 热线电话：（ 86 ） 0755 — 8427 3671 电子邮件： osa@cuhk.edu.cn 办公时间： 8:30-12:00 13:00-17:30 ( 周一至周五 ) 扫码关注 OSA 公众号 ：",
}
```

additional_info:

```json
{
    "timezone": "Asia/Shanghai"
}
```

Response:

```json
{
    "events": [
        {
            "event_type": "registration",
            "end_time": "2023-04-06T00:00:00 Asia/Shanghai",
            "summary": "2023大学杰出毕业生奖提名者自荐材料征集",
            "venue": "https://wj.cuhk.edu.cn/vm/YVgulbu.aspx"
        }
    ]
}
```

## Backend API Endpoint

### `events` endpoint

This API is meant for all-in-one communication between clients and the backend server. It will automatically authenticate users' secret, update users' events and return their latest events.

Endpoint: `http://public-ip:8080/events`

Params: 

```yaml
user_id: string
secret: string
```

Response: 

A json string of events (something like the below)

```json
[
    {
        "event_type": "registration",
        "end_time": "2023-04-06T12:00:00 Asia/Shanghai",
        "summary": "2023大学杰出毕业生奖提名者自荐材料征集",
        "venue": "https://wj.cuhk.edu.cn/vm/YVgulbu.aspx"
    }
]
```