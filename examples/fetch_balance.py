# 주식잔고조회
import korea_investment_stock
import pprint

with open("../../koreainvestment.key") as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()
acc_no = lines[2].strip()

broker = korea_investment_stock.KoreaInvestment(
    api_key=key,
    api_secret=secret,
    acc_no=acc_no
)
resp = broker.fetch_balance()
pprint.pprint(resp)
