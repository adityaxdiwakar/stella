import random

# Suggested by Periklean

quotes = [
    "As I see it, yes.",
    "Ask again later.",
    "Better not tell you now.",
    "Cannot predict now.",
    "Concentrate and ask again.",
    "Don’t count on it.",
    "It is certain.",
    "It is decidedly so.",
    "Most likely.",
    "My reply is no.",
    "My sources say no.",
    "Outlook not so good.",
    "Outlook good.",
    "Reply hazy, try again.",
    "Signs point to yes.",
    "Very doubtful.",
    "Without a doubt.",
    "Yes.",
    "Yes – definitely.",
    "You may rely on it.",
    "Dwyer has no doubt, trust the Dwyer.",
    "Ask yourself, what would Cramer say?"
]

async def main(message, canary=False):
    try:
        premsg = ["**[Canary]** ", " "][not canary]
        index = random.randint(0, len(quotes) - 1)
        await message.channel.send(premsg + quotes[index])
    except:
        await message.channel.send(premsg + "Unfortunately, an error has occured. Developers have been notified of the issues.")