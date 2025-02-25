import korea_investment_stock
import pprint

with open("../../koreainvestment.key") as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()
ACC_NO = "63398082-01"

broker = korea_investment_stock.KoreaInvestment(
    api_key=key,
    api_secret=secret,
    acc_no=ACC_NO
)

resp = broker.fetch_price("005930")
pprint.pprint(resp)
