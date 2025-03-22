import os
from unittest import TestCase, skip

from korea_investment_stock import KoreaInvestment
from korea_investment_stock.koreainvestmentstock import RETURN_CD


class TestKoreaInvestment(TestCase):
    @classmethod
    def setUpClass(cls):
        api_key = os.getenv('STOCK_API_KOREA_INVESTMENT_API_KEY')
        api_secret = os.getenv('STOCK_API_KOREA_INVESTMENT_API_SECRET')
        acc_no = os.getenv('STOCK_API_KOREA_INVESTMENT_ACCOUNT_NO')

        cls.kis_korea = KoreaInvestment(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
            exchange='서울' # todo: exchange는 제거 예정
        )

        # todo: exchange 구분없이 동작을 했으면 한다
        cls.kis_nasdaq = KoreaInvestment(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
            exchange='나스닥' # todo: exchange는 제거 예정
        )

    def test_stock_info(self):
        test_cases = [
            ("samsung", "005930", "KR"),
            ("apple", "AAPL", "US")
        ]

        for name, ticker, market in test_cases:
            with self.subTest(name=name):
                resp = self.kis_korea.fetch_stock_info(ticker, market)
                self.assertEqual(resp['rt_cd'], RETURN_CD["SUCCESS"])
                self.assertEqual(resp['output']['shtn_pdno'], ticker)

    def test_fetch_domestic_price(self):
        test_cases = [
            ("samsung", "005930"),
            ("etf", "294400")
        ]

        for name, ticker in test_cases:
            with self.subTest(name=name):
                resp = self.kis_korea.fetch_price(ticker)
                print(resp)
                self.assertEqual(resp['rt_cd'], RETURN_CD["SUCCESS"])


    def test_fetch_oversea_price(self):
        test_cases = [
            ("apple", "AAPL"),
            ("tlt", "TLT"),
            ("jepq", "JEPQ"),
            ("divo", "DIVO") # todo: 이거 조회가 안됨
        ]

        for name, ticker in test_cases:
            with self.subTest(name=name):
                resp = self.kis_nasdaq.fetch_price(ticker)
                print(resp)
                self.assertEqual(resp['rt_cd'], RETURN_CD["SUCCESS"])

    @skip("Skipping test_fetch_kospi_symbols")
    def test_fetch_kospi_symbols(self):
        resp = self.kis_korea.fetch_kospi_symbols()
        print(resp)

    def test_fetch_price_detail_oversea(self):
        resp = self.kis_nasdaq.fetch_price_detail_oversea("AAPL")
        print(resp)
