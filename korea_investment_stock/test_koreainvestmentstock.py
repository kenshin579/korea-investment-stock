import os
from unittest import TestCase

from korea_investment_stock import KoreaInvestment

class TestKoreaInvestment(TestCase):
    @classmethod
    def setUpClass(cls):
        cls.kis = KoreaInvestment(
            api_key=os.getenv('STOCK_API_KOREA_INVESTMENT_API_KEY'),
            api_secret=os.getenv('STOCK_API_KOREA_INVESTMENT_API_SECRET'),
            acc_no=os.getenv('STOCK_API_KOREA_INVESTMENT_ACCOUNT_NO')
        )

    def test_stock_info(self):
        info = self.kis.fetch_stock_info("005930")
        print(info)

