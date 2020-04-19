import json

async def add_ref(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]
    
    if message.author.id not in [577714296599871488, 119247462996115456, 192696739981950976]:
        await message.channel.send(premsg + "🔑 Sorry! You are not allowed to create tags.")
        return

    if len(message.content.split(" ")) < 3:
        await message.channel.send(premsg + "Please enter an actual tag to be added!")
        return

    tag_title = message.content.split(" ")[1]

    tag_content = " ".join(message.content.split(" ")[2:])
    
    with open("bin/tags.json", "r") as f:
        tags = json.load(f)

    tags.update({tag_title: tag_content})

    with open("bin/tags.json", "w") as f:
        json.dump(tags, f, indent=4)

    await message.channel.send(premsg + f"Awesome! The **{tag_title}** tag was added!")

async def show_tags(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    with open("bin/tags.json", "r") as f:
        tags = json.load(f)

    tags_str = ", ".join(list(tags.keys()))
    
    await message.channel.send(premsg + f"The current tags available for use are: ``{tags_str}``. Please only test in #bot-spam")

async def use_tag(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    with open("bin/tags.json", "r") as f:
        tags = json.load(f)


    if len(message.content.split(" ")) < 2:
        await message.channel.send(premsg + "Please enter a tag name to run!")
        return

    tag_title = message.content.split(" ")[1]

    if tag_title not in tags:
        await message.channel.send(premsg + "Tag not found in the system, sorry! Run ``?showtags`` to see what tags are available.")
        return

    await message.channel.send(premsg + tags[tag_title])

async def rm_tag(message, canary=False):
    premsg = ["**[Canary]** ", " "][not canary]

    if message.author.id not in [577714296599871488, 119247462996115456, 192696739981950976]:
        await message.channel.send(premsg + "🔑 Sorry! You are not allowed to delete tags.")
        return

    if len(message.content.split(" ")) < 2:
        await message.channel.send(premsg + "Please enter a tag name to delete!")
        return

    tag_title = message.content.split(" ")[1]

    with open("bin/tags.json", "r") as f:
        tags = json.load(f)

    if tag_title not in tags:
        await message.channel.send(premsg + "Tag not found in the system, sorry! Run ``?showtags`` to see what tags are available.")
        return

    del tags[tag_title]

    with open("bin/tags.json", "w") as f:
        tags = json.dump(tags, f, indent=4)
    
    await message.channel.send(premsg + "Success! I've sent the tag into the nearest black hole.")
    