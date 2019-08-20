async def bearish(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    await message.channel.send(premsg + "Here you are: https://i.imgur.com/Fiua3bN.png")

async def tradewar(message, canary=False):
    premsg = ["**[Canary]**", " "][not canary]
    await message.channel.send(premsg + "Here you are: https://i.imgur.com/JXNlKXj.png")
