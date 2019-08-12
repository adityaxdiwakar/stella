from dotenv import load_dotenv
import os

load_dotenv() #grab env variables from config

prefix = os.getenv("BOT_PREFIX")
is_dev = bool(int(os.getenv("IS_DEV")))

from commands import ping
from commands import charts
from commands import futcharts
from commands import fxcharts

import copy
import discord
import time

class Stella(discord.Client):
    async def on_ready(self):
        print('Logged on as', self.user)

    async def on_message(self, message):
        message.content = message.content.lower()
        if message.content.startswith(f"{prefix}"):
            if message.content[len(prefix)] == " ":
                char_array = list(message.content)
                del char_array[len(prefix)]
                message.content = ""
                for char in char_array:
                    message.content += char

            if message.content.startswith(f"{prefix}ping"):
                await ping.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}c"):
                await charts.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}f"):
                await futcharts.main(message, canary=is_dev)

            elif message.content.startswith(f"{prefix}x"):
                await fxcharts.main(message, canary=is_dev)

ctx = Stella()
ctx.run(os.getenv("BOT_TOKEN"))
