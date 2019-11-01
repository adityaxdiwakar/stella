from utils import single_report as sr
from utils import anom_all as aa
from utils import custom_rep as cr
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
        report_type = {"gfs": "gfs", "gefs": "gfs-ensemble", "ecmwf": "ecmwf", "eps": "ecmwf-ensemble"}[message_split[1]]
    except:
        await message.channel.send(premsg + "Sorry, your report type could not be interpreted")
        return
    msg = await message.channel.send(premsg + "Generating chart, stand by.")
    chart_link = sr.main(report_type)
    if chart_link[1] == None:
        await msg.edit(content=premsg + f"An error occured: {chart_link[0]}")
        return
    await msg.edit(content=premsg + f"{chart_link[0]} with an addition of **{chart_link[1]}** FoF and **{chart_link[2]}** over 30yr avg.")


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

async def custom(message, canary=False):
    premsgs = ["**[Canary]** ", " "]
    premsg = premsgs[not canary]

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, I need atleast one line to plot!")
        return

    feed_forward = message_split[1:]

    for x in range(len(feed_forward)):
        item = feed_forward[x]
        if item.split("@")[0] not in ["gfs", "gefs", "eps", "ecmwf"]:
            await message.channel.send(premsg + "Sorry, you entered an invalid report type")
            return
        beginning = {"gfs": "gfs", "ecmwf": "ecmwf", "eps": "ecmwf-ensemble", "gefs": "gfs-ensemble"}[item.split("@")[0]]
        feed_forward[x] = f"{beginning}@{item.split('@')[1]}"

    msg = await message.channel.send(premsg + "Generating chart, stand by.")
    chart_link = cr.main(feed_forward)
    await msg.edit(content=premsg + f"Here you go: {chart_link}")

