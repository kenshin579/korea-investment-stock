# 주식잔고조회
import kis
import pprint

with open("../../koreainvestment.key") as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()

broker = kis.KoreaInvestment(
    api_key = key,
    api_secret = secret,
    acc_no = "63398082-01"
)

# fetch_tickers
symbols = broker.fetch_kosdaq_symbols()
symbols.to_excel("kosdaq_code.xlsx", index=False)