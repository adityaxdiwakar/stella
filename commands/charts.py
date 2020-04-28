from urllib.parse import urlencode
from io import BytesIO
import discord
import requests
import io
import os
import time

prefix = os.getenv("BOT_PREFIX")

timeframes = ["i1", "i3", "i5", "i15", "i30", "d", "w", "m"]
timeframe_names = ["1 minute intraday", "3 minute intraday", "5 minute intraday", "15 minute intraday", "30 minute intraday", "daily", "weekly", "monthly"]


def create_chart(ticker, chart_type, is_content=False):
    query = {
        "t": ticker,
        "ty": "c",
        "ta": str(int(chart_type < 6)),
        "p": timeframes[chart_type],
        "s": "l"
    }

    root_url = "https://elite.finviz.com/chart.ashx"
    if chart_type > 4 and chart_type != 5:
        root_url = "https://finviz.com/chart.ashx"

    if chart_type == 5:
        query["ta"] = "st_c,sch_200p,sma_50,sma_200,sma_20,sma_100,bb_20_2,rsi_b_14,macd_b_12_26_9,stofu_b_14_3_3"

    qstr = urlencode(query)



    file = requests.get(f"{root_url}?{qstr}", headers={"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.162 Safari/537.36"})

    rn = round(time.time())

    if len(file.content) < 7500 or len(ticker) > 8:
       return (None, "Chart not found! An error occured, try again. If you need futures, use ``?f``.")

    if is_content == True:
        return (file.content, None)

    with open(f"/var/www/html/u/ca/{rn}.png", "wb") as f:
        f.write(file.content)

    return (rn, None)


async def main(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    try:
        chart_type = message.content[len(prefix) + 1] #the third id
        if chart_type == " ":
            chart_type = 5
        chart_type = int(chart_type)
    except ValueError:
        await message.channel.send(premsg + "Something went wrong with your request, check the command try again!")
        return

    if chart_type > 7 or chart_type < 0:
        await message.channel.send(premsg + "You asked for a chart type that we don't have, check the bot command channel for help!")
        return

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, couldn't identify your ticker! Try again!")
        return

    ticker = message_split[1]

    msg = await message.channel.send(premsg + "Grabbing chart, stand by.")

    try:
        rn, error = create_chart(ticker, chart_type)

        if error != None:
            await msg.edit(content=premsg + error)
            return

        await msg.edit(content=premsg + f"Alright, here's your {timeframe_names[chart_type]} chart: https://img.adi.wtf/ca/{rn}.png")

    except Exception as e:
        await msg.edit(content=premsg + f"Something went wrong, contact <@192696739981950976> with ```{e}```")


async def multi(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    try:
        chart_type = message.content[len(prefix) + 2] #the third id
        if chart_type == " ":
            chart_type = 5
        chart_type = int(chart_type)
    except ValueError:
        await message.channel.send(premsg + "Something went wrong with your request, check the command try again!")
        return

    if chart_type > 7 or chart_type < 0:
        await message.channel.send(premsg + "You asked for a chart type that we don't have, check the bot command channel for help!")
        return

    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, couldn't identify your ticker! Try again!")
        return

    msg = await message.channel.send(premsg + "Grabbing chart, stand by.")

    joined_split = " ".join(message_split[1:])
    ticker_split = joined_split.split(",")
    if len(ticker_split) == 1:
        ticker_split = joined_split.split(" ")
    if len(ticker_split) == 1:
        await msg.edit(content=premsg + " Sorry, to use the ``?mc`` command, enter more than one ticker!")
        return

    ticker_split = [x.strip() for x in ticker_split]

    if len(ticker_split) > 8:
        await msg.edit(content=premsg + "Sorry, too many tickers entered. There is a maximum of 8!")
        return

    link_times = []
    errored_tickers = []
    t_error = None
    for ticker in ticker_split:
        rn, error = create_chart(ticker, chart_type, is_content=True)

        if error != None:
            t_error = error
            errored_tickers.append(ticker)
        else:
            link_times.append(rn)

    files = [discord.File(BytesIO(x), filename=f"{time.time()}.png") for x in link_times]
    omitted_tickers = " ".join(errored_tickers)
    omitted_tickers = omitted_tickers.upper()

    text = f"{premsg}Alright, here are your {timeframe_names[chart_type]} charts."
    if len(errored_tickers) > 0:
        text += f" ``{omitted_tickers}`` had to be ommitted due to errors during fetching."
    await message.channel.send(text, files=files)
    await msg.delete()

