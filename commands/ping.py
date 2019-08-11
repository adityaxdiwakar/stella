from datetime import datetime
import time

async def main(message, canary=False):
    if canary == True:
        premsg = "**[Canary]** "
    msg = await message.channel.send(premsg + f":ping_pong: Calculating ping...")
    roundtrip = round((msg.created_at.timestamp() - message.created_at.timestamp()) * 1000)
    await msg.edit(content=premsg + f":ping_pong: WS roundtrip complete in {roundtrip}ms!")