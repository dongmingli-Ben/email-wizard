from __future__ import print_function

import logging
import json

import grpc
import service.parse_pb2 as pb2
import service.parse_pb2_grpc as pb2_grpc


def run(email, **kwargs):
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    print("Will try to request the server ...")
    with grpc.insecure_channel('localhost:50052') as channel:
        stub = pb2_grpc.ParserStub(channel)
        response = stub.ParseEmail(pb2.EmailContentRequest(
            email=json.dumps(email), additional_info=json.dumps(kwargs)))
        print("GetEmails client received: ")
        emails = json.loads(response.message)
        print(json.dumps(emails, ensure_ascii=False, indent=4))


if __name__ == '__main__':
    logging.basicConfig()

    email = {
        "subject": "重要通知 | 2023大学杰出毕业生奖提名者自荐材料征集",
        "sender": "someone@example.com",
        "date": "2023-06-07T13:58:41Z",
        "recipient": ["whoami@outlook.com"],
        "content": "亲爱的 同学： 你好！ 经由学院、书院推荐，你现已进入今年“大学杰出毕业生奖（ Presidential Award for Outstanding Students ）”的提名名单。为帮助学院、书院及遴选委员会更全面、深入了解你的情况，我们诚挚地邀请你提供自荐信以及证明等资料作为补充资料供相关部门及委员会 参 考。 提交方式： 请点击下方链接或扫描二维码，提交相关材料。 链接： https://wj.cuhk.edu.cn/vm/YVgulbu.aspx 二维码： 截止日期： 4 月 6 日（星期四）中午 12:00 * 注意 * ： 1. 逾期未提交或未按要求提交的同学，将视作放弃自荐机会，请知悉。 2. 如不接受提名，请 在 4 月 5 日 之前回复此邮件 告知，我们将通知学院、书院尽快提供后补人选。 3 ．如果你已提交，则无需再次提交。 奖项介绍 大学杰出毕业生奖 （ Presidential Award for Outstanding Students ）是大学毕业生的最高荣誉，旨在表彰在大学期间具有杰出的学术表现、卓越的领导力以及对大学发展有贡献的同学，今年提名人数（含学院及书院提名）约为毕业生总人数的 3.5% ，获奖人数约为毕业生总人数的 1% 。 关于该奖项的评选，如有疑问，可联系学生事务处杨老师 yangxiaoyu@cuhk.edu.cn 。 顺祝， 平安顺遂 Best Regards, Office of Student Affairs The Chinese University of Hong Kong, Shenzhen Hotline: (86) 0755 — 8427 3671 E-mail: osa@cuhk.edu.cn Office Hours: 8:30-12:00 13:00-17:30 (Mon-Fri) 学生事务处 香港中文大学（深圳） 热线电话：（ 86 ） 0755 — 8427 3671 电子邮件： osa@cuhk.edu.cn 办公时间： 8:30-12:00 13:00-17:30 ( 周一至周五 ) 扫码关注 OSA 公众号 ：",
    }
    run(email, timezone='Asia/Shanghai')
