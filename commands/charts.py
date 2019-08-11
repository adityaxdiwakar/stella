from urllib.parse import urlencode
import discord
import requests
import io
import os

prefix = os.getenv("BOT_PREFIX")

async def main(message, canary=False):
    if canary == True:
        premsg = "**[Canary]** "
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

    timeframes = ["i1", "i3", "i5", "i15", "i30", "d", "w", "m"]
    timeframe_names = ["1 minute intraday", "3 minute intraday", "5 minute intraday", "15 minute intraday", "30 minute intraday", "daily", "weekly", "monthly"]
    
    message_split = message.content.split(" ")
    if len(message_split) < 2:
        await message.channel.send(premsg + "Sorry, couldn't identify your ticker! Try again!")
        return

    ticker = message_split[1]

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

    file = requests.get(f"{root_url}?{qstr}")

    await message.channel.send(premsg + f"Alright, here's your {timeframe_names[chart_type]} chart:", file=discord.File(io.BytesIO(file.content), "chart.png"))

