from urllib.parse import urlencode
import discord
import requests
import io
import os
import time

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

    timeframes = ["m5", "h1", "d1", "w1", "mo"]
    timeframe_names = ["5 minute", "hourly", "daily", "weekly", "monthly"]

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, couldn't identify your ticker! Try again!")
        return

    ticker = message_split[1]


    currencies = ["EURUSD", "GBPUSD", "USDJPY", "USDCAD", "USDCHF", "AUDUSD", "NZDUSD", "EURGBP", "GBPJPY", "BTCUSD"]
    if ticker.upper() not in currencies:
        reply = "Hey, I'm sorry but I can only do the following tickers:"
        for tick in currencies:
            reply += f", {tick}"
        await message.channel.send(premsg + reply)
        return

    msg = await message.channel.send(premsg + "Grabbing chart, stand by.")

    try:
        root_url = "https://finviz.com/fx_image.ashx"

        timeframe = timeframes[chart_type]
        file = requests.get(f"{root_url}?{ticker}_{timeframe}_l.png")

        rn = round(time.time())
        with open(f"/var/www/html/u/fx/{rn}.png", "wb") as f:
            f.write(file.content)

        await message.channel.send(premsg + f"Alright, here's your {timeframe_names[chart_type]} chart: https://img.adi.wtf/fx/{rn}.png")

    except Exception as e:
        await msg.edit(content=f"Something went wrong, contact <@192696739981950976> with ```{e}```")