# 주식잔고조회
import korea_investment_stock
import pprint

with open("../../koreainvestment.key") as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()
ACC_NO = "63398082-01"

broker = korea_investment_stock.KoreaInvestment(
    api_key = key,
    api_secret = secret,
    acc_no = ACC_NO
)

# fetch_symbols
symbols = broker.fetch_symbols()
print(symbols.head())

cond = symbols['그룹코드'] == 'ST'
print(symbols[cond])
symbols.to_excel("korea_code.xlsx", index=False)

