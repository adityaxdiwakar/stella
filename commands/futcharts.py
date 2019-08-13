from urllib.parse import urlencode
import discord
import requests
import io
import os

prefix = os.getenv("BOT_PREFIX")

async def main(message, canary=False):
    premsgs = ["**[Canary]** ", " "]
    premsg = premsgs[not canary]
    try:
        chart_type = message.content[len(prefix) + 1]
        if chart_type == " ":
            chart_type = 1
        chart_type = int(chart_type)
    except ValueError:
        await message.channel.send(premsg + "Something went wrong with your request, check the command try again!")
        return

    if chart_type > 5 or chart_type < 1:
        await message.channel.send(premsg + "You asked for a chart type that we don't have, check the bot command channel for help!")
        return

    chart_type -= 1

    timeframes = ["m5", "h1", "d1", "w1", "m1"]
    timeframe_names = ["5 minute", "hourly", "daily", "weekly", "monthly"]

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, couldn't identify your ticker! Try again!")
        return

    ticker = message_split[1]

    query = {
        "t": ticker,
        "p": timeframes[chart_type],
        "s": "l"
    }

    root_url = "https://finviz.com/fut_chart.ashx"

    qstr = urlencode(query)

    file = requests.get(f"{root_url}?{qstr}")

    if len(file.content) == 0:
        await message.channel.send(premsg + f"Chart not found! An error occured, try again. If you need cash market equities, use ``?c``.")
        return

    await message.channel.send(premsg + f"Alright, here's your {timeframe_names[chart_type]} chart:", file=discord.File(io.BytesIO(file.content), "chart.png"))

