"""cancel_order
"""
import pprint
import korea_investment_stock

with open("../../koreainvestment.key", encoding="utf-8") as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()
ACC_NO = "63398082-01"

broker = korea_investment_stock.KoreaInvestment(
    api_key=key,
    api_secret=secret,
    acc_no=ACC_NO
)

resp = broker.cancel_order(
    org_no="91252",
    order_no="0000120154",
    quantity=2,     # 취소하고자하는 수량
    total=False     # 잔량일부
)
pprint.pprint(resp)
