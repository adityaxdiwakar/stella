#environemnt information
from dotenv import load_dotenv
load_dotenv()

import requests
import sys
import os

argument = sys.argv[1]

if argument == "open":
    text = "🔔🔔🔔 The cash market is now open! Have a great trading day!"

if argument == "open-futures":
    text = "🔔🔔🔔 The futures market is back open for the week! Have a good one!"

if argument == "close":
    text = "🔔🔔🔔 The cash market is now closed! Hope you had a great trading day!"

if argument == "close-futures":
    text = "🔔🔔🔔 The futures market is now closed for the weekend!"

requests.post(os.getenv("WEBHOOK"), data={"content": text})
