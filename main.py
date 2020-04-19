from bs4 import BeautifulSoup
from dotenv import load_dotenv
import datetime
import requests
import os

load_dotenv() #grab env variables from config

prefix = os.getenv("BOT_PREFIX")
is_dev = bool(int(os.getenv("IS_DEV")))

from commands import ping
from commands import charts
from commands import futcharts
from commands import fxcharts
from commands import refs
from commands import ng_all
from commands import ng_rep
from commands import earnings
from commands import eightball
from commands import evalmod
from commands import position_size
from commands import tdcommands

from utils import reactions

import copy
import discord
import time
import asyncio
import ast

# {prefix: component}
module_links = {
    "ping": ping.main,
    "c": charts.main,
    "mc": charts.multi,
    "fun": tdcommands.fundamentals,
    "div": tdcommands.dividends,
    "f": futcharts.main,
    "x": fxcharts.main,
    "ngall": ng_all.main,
    "rc": ng_rep.custom,
    "anom": ng_rep.all_anom,
    "r": ng_rep.main,
    "earnings": earnings.company,
    "addtag": refs.add_ref,
    "8ball": eightball.main,
    "showtags": refs.show_tags,
    "tag": refs.use_tag,
    "eval": evalmod.main,
    "pos": position_size.calculator
}



class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)

    async def on_message(self, message):
        message.content.split(" ")[0] = message.content.split(" ")[0].lower()
        if message.content.startswith(f"{prefix}"):
            if message.content[len(prefix)] == " ":
                char_array = list(message.content)
                del char_array[len(prefix)]
                message.content = ""
                for char in char_array:
                    message.content += char

        for mod in module_links:
            if message.content.startswith(prefix + mod):
                await module_links[mod](message, canary=is_dev)
                break

    async def on_raw_reaction_add(self, payload):
        await reactions.handler(self, payload, "add")
    
    async def on_raw_reaction_remove(self, payload):
        await reactions.handler(self, payload, "remove")

        

ctx = Stella()
ctx.run(os.getenv("BOT_TOKEN"))
