import os
from unittest import TestCase

from parameterized import parameterized

from korea_investment_stock import KoreaInvestment
from korea_investment_stock.koreainvestmentstock import RETURN_CD

class TestKoreaInvestment(TestCase):
    @classmethod
    def setUpClass(cls):
        cls.kis = KoreaInvestment(
            api_key=os.getenv('STOCK_API_KOREA_INVESTMENT_API_KEY'),
            api_secret=os.getenv('STOCK_API_KOREA_INVESTMENT_API_SECRET'),
            acc_no=os.getenv('STOCK_API_KOREA_INVESTMENT_ACCOUNT_NO')
        )

    def test_stock_info(self):
        test_cases = [
            ("samsung", "005930", "KR"),
            ("apple", "AAPL", "US")
        ]

        for name, ticker, market in test_cases:
            with self.subTest(name=name):
                resp = self.kis.fetch_stock_info(ticker, market)
                self.assertEqual(resp['rt_cd'], RETURN_CD["SUCCESS"])
                self.assertEqual(resp['output']['shtn_pdno'], ticker)
