from dotenv import load_dotenv
load_dotenv()

from utils import tda

import discord
import json
import os

refresh_token = os.getenv("REFRESH_TOKEN")
consumer_key = os.getenv("CONSUMER_KEY")
tdctx = tda.TDClient(refresh_token, consumer_key)

async def fundamentals(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    ticker = message.content.split(" ")[1]
    
    try:
        data = tdctx.fundamentals(ticker)
    except Exception as e:
        await message.channel.send(premsg + "Something went wrong with the request, check the ticker or try again! If the issue persists, please contact aditya#0001.")
        return

    titles = ["Price-Earnings Ratio", "Price-Book Ratio", "Debt-Equity Ratio", 
              "PEG Ratio", "Beta", "Dividend Yield", "Equity Return", "Market Cap"]
    keys = ["peRatio", "pbRatio", "totalDebtToEquity", "pegRatio", "beta", "dividendYield",
            "returnOnEquity", "marketCap"]
    
    fundamentals = data[ticker.upper()]["fundamental"]
    wl_data = {k:fundamentals[v] for k,v in zip(titles,keys) if fundamentals[v] != 0}
    company_name = data[ticker.upper()]["description"]

    if "Dividend Yield" in wl_data:
        wl_data["Dividend Yield"] = str(wl_data['Dividend Yield']) + "%"
    if "Market Cap" in wl_data:
        wl_data["Market Cap"] = "$" + str(round(wl_data["Market Cap"]/1000, 2)) + "B"

    embed = discord.Embed(title=f"[{ticker.upper()}] {company_name} Fundamentals")
    for name,value in wl_data.items():
        embed.add_field(name=name, value=str(value))
    await message.channel.send(premsg, embed=embed)

async def dividends(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    ticker = message.content.split(" ")[1]

    try:
        data = tdctx.fundamentals(ticker)
    except Exception as e:
        await message.channel.send(premsg + "Something went wrong with the request, check the ticker or try again! If the issue persists, please contact aditya#0001.")
        return

    titles = ["Ex-Div Date", "Dividend ($)", "Yield", "Pay Date"]
    keys = ["dividendDate", "dividendPayAmount", "dividendYield", "dividendPayDate"]

   


    fundamentals = data[ticker.upper()]["fundamental"]
    wl_data = {k:fundamentals[v] for k,v in zip(titles,keys)}
    company_name = data[ticker.upper()]["description"]
    
    wl_data["Yield"] = str(wl_data['Yield']) + "%"
    wl_data["Ex-Div Date"] = wl_data["Ex-Div Date"].split(" ")[0]
    wl_data["Pay Date"] = wl_data["Pay Date"].split(" ")[0]
    wl_data["Dividend ($)"] = "$" + str(wl_data["Dividend ($)"])

    embed = discord.Embed(title=f"[{ticker.upper()}] {company_name} Fundamentals")
    for name,value in wl_data.items():
        embed.add_field(name=name, value=str(value))
    await message.channel.send(premsg, embed=embed)
