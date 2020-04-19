import requests

API_ERROR = Exception("API Error Occured")
class TDClient():
    
    def __init__(self, refresh_token, consumer_key):
        self.refresh = refresh_token
        self.conkey = consumer_key
        self.ROOT = "https://api.tdameritrade.com/v1"

    def get_access_token(self):
        payload = {
            "grant_type":"refresh_token",
            "refresh_token":self.refresh,
            "client_id":f"{self.conkey}@AMER.OAUTHAP",
            "redirect_uri":"http://127.0.0.1"
        }
        r = requests.post(f"{self.ROOT}/oauth2/token", data=payload)
        if "access_token" not in r.json():
            raise API_ERROR
            return
        return r.json()["access_token"]

    def fundamentals(self, ticker):
        query = {
            "symbol": ticker,
            "projection": "fundamental"
        }
        headers = {
            "Authorization": f"Bearer {self.get_access_token()}"
        }
        r = requests.get(f"{self.ROOT}/instruments", 
                          params=query,
                          headers=headers)
        if ticker.upper() not in r.json() or "fundamental" not in r.json()[ticker.upper()]:
            raise API_ERROR
        return r.json()


