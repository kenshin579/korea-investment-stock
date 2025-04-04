"""
나스닥 지정가 매수
"""
import pprint
import korea_investment_stock

with open("../../../koreainvestment.key", encoding='utf-8') as f:
    lines = f.readlines()

key = lines[0].strip()
secret = lines[1].strip()
acc_no = lines[2].strip()

broker = korea_investment_stock.KoreaInvestment(
    api_key=key,
    api_secret=secret,
    acc_no=acc_no,
    # exchange='나스닥' # todo: exchange는 제거 예정
)

resp = broker.create_limit_sell_order(
    symbol="TQQQ",
    price=30,
    quantity=5
)
pprint.pprint(resp)
