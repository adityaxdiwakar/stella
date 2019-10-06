from utils import single_report as sr
from utils import anom_all as aa
from datetime import datetime
import time

async def main(message, canary=False):
    premsgs = ["**[Canary]** ", " "]
    premsg = premsgs[not canary]

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, please enter a report type!")
        return

    try:
        report_type = {"gfs": "gfs", "gefs": "gfs-ensemble", "emcwf": "emwf", "eps": "ecmwf-ensemble"}[message_split[1]]
    except:
        await message.channel.send(premsg + "Sorry, your report type could not be interpreted")
        return
    msg = await message.channel.send(premsg + "Generating chart, stand by.")
    chart_link = sr.main(report_type)
    await msg.edit(content=premsg + f"Here you go: {chart_link}")


async def all_anom(message, canary=False):
    premsgs = ["**[Canary]** ", " "]
    premsg = premsgs[not canary]

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, please enter a report type!")
        return

    try:
        report_type = {"gfs": "gfs", "gefs": "gfs-ensemble", "emcwf": "emwf", "eps": "ecmwf-ensemble"}[message_split[1]]
    except:
        await message.channel.send(premsg + "Sorry, your report type could not be interpreted")
        return
    msg = await message.channel.send(premsg + "Generating chart, stand by.")
    chart_link = aa.main(report_type)
    await msg.edit(content=premsg + f"Here you go: {chart_link}")

