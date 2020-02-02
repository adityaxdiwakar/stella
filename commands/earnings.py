from bs4 import BeautifulSoup
import requests
import discord

async def company(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    msg = await message.channel.send(premsg + "Loading...")

    if len(message.content.split(" ")) < 2:
        await msg.edit(content=premsg + "Sorry, please enter the ticker you would like to lookup earnings for!")
        return

    ticker = message.content.split(" ")[1]

    r = requests.get(f"https://earningswhispers.com/stocks/{ticker}")

    soup = BeautifulSoup(r.content, 'html.parser')
    whisper_box = soup.find(id="whisperbox")
    date_box = soup.find(id="datebox")
    try:
        company_name = soup.find(id="compname").text
    except AttributeError:
        await msg.edit(content=premsg + "Sorry, that company does not exist in the earnings database. Try again?")
        return

    date_box_children = [child.text for child in date_box.children if child.text != ""]

    whisper_box_children = [child.text for child in whisper_box.children if child.text != ""]

    embed = discord.Embed(title=ticker.upper(), description=f"{company_name} Earnings ({whisper_box_children[2]})")
    
    try:
        embed.add_field(name="Reporting Date", value=f"{date_box_children[1]} ({date_box_children[0]})", inline=True)
    except:
        embed.add_field(name="Reporting Date", value="N/A")
    try:
        embed.add_field(name="Reporting Time", value=date_box_children[2], inline=True)
    except:
        embed.add_field(name="Reporting Time", value="N/A")
    try:
        if whisper_box_children[1] == "$0.00":
            whisper_box_children[1] == "N/A"
        embed.add_field(name="EW EPS", value=whisper_box_children[1], inline=True)
    except:
        embed.add_field(name="EW EPS", value="N/A")
    try:
        value = whisper_box_children[3].split(" ")[2]
        if value == "$0.00":
            value = "N/A"
        embed.add_field(name="WS EPS", value=value, inline=True)
    except:
        embed.add_field(name="WS EPS", value="N/A")
    try:
        value = " ".join(whisper_box_children[4].split(" ")[1:3])
        if value == "$0.00 Mil":
            value = "N/A"
        embed.add_field(name="WS Revenue", value=value, inline=True)
    except:
        embed.add_field(name="WS Revenue", value="N/A")

    await msg.edit(content=premsg, embed=embed)
