#environemnt information
from dotenv import load_dotenv
load_dotenv()

import requests
import sys
import os

argument = sys.argv[1]

text = {
	"open": "🔔🔔🔔 The cash market is now open! Have a great trading day!",
	"open-futures": "🔔🔔🔔 The futures market is back open for the week! Have a good one!",
	"close": "🔔🔔🔔 The cash market is now closed! Hope you had a great trading day!",
	"close-futures": "🔔🔔🔔 The futures market is now closed for the weekend!",
	"euro-open": "🔔🔔🔔 The LSE (London Stock Exchange) is now open",
	"euro-close": "🔔🔔🔔 The LSE (London Stock Exchange) is now closed",
	"tokyo-open": "🔔🔔🔔 The TSE (Tokyo Stock Exchange) is now open",
	"tokyo-close": "🔔🔔🔔 The TSE (Tokyo Stock Exchange) is now closed",
	"hk-open": "🔔🔔🔔 The HKE (Hong Kong Stock Exchange) is now open",
	"hk-close": "🔔🔔🔔 The HKE (Hong Long Stock Exchange) is now closed",
	"sse-open": "🔔🔔🔔 The SSE (Shanghai Stock Exchange) is now open",
	"sse-close": "🔔🔔🔔 The SSE (Shanghai Stock Exchange) is now close"
}[argument]

for webhook in [os.getenv("WEBHOOK_1"), os.getenv("WEBHOOK_2")]:
	requests.post(webhook, json={"content": text})
