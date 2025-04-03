import os
from unittest import TestCase, skip

from korea_investment_stock import KoreaInvestment
from korea_investment_stock.koreainvestmentstock import API_RETURN_CODE


class TestKoreaInvestment(TestCase):
    @classmethod
    def setUpClass(cls):
        api_key = os.getenv('STOCK_API_KOREA_INVESTMENT_API_KEY')
        api_secret = os.getenv('STOCK_API_KOREA_INVESTMENT_API_SECRET')
        acc_no = os.getenv('STOCK_API_KOREA_INVESTMENT_ACCOUNT_NO')

        cls.kis = KoreaInvestment(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
        )

    def test_stock_info(self):
        test_cases = [
            ("samsung", "005930", "KR"),
            ("apple", "AAPL", "US")
        ]

        for name, ticker, market in test_cases:
            with self.subTest(name=name):
                resp = self.kis.fetch_stock_info(ticker, market)
                self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])
                self.assertEqual(resp['output']['shtn_pdno'], ticker)

    def test_fetch_search_stock_info(self):
        test_cases = [
            ("samsung", "005930"),
            ("etf", "294400")
        ]

        for name, ticker in test_cases:
            with self.subTest(name=name):
                resp = self.kis.fetch_search_stock_info(ticker)
                print(resp)
                self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])
                self.assertIn('output', resp)
                self.assertIn('frbd_mket_lstg_dt', resp['output'])


    def test_fetch_price(self):
        test_cases = [
            ("samsung", "005930"),
            ("etf", "294400"),
            ("apple", "AAPL")
        ]

        for name, ticker in test_cases:
            with self.subTest(name=name):
                resp = self.kis.fetch_price(ticker)
                print(resp)
                self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])


    def test_fetch_oversea_price(self):
        test_cases = [
            ("apple", "AAPL"),
            ("tlt", "TLT"),
            ("jepq", "JEPQ"),
            ("divo", "DIVO")
        ]

        for name, ticker in test_cases:
            with self.subTest(name=name):
                resp = self.kis.fetch_price(ticker)
                print(resp)
                self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])

    @skip("Skipping test_fetch_kospi_symbols")
    def test_fetch_kospi_symbols(self):
        resp = self.kis.fetch_kospi_symbols()
        print(resp)
        self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])

    def test_fetch_price_detail_oversea(self):
        resp = self.kis.fetch_price_detail_oversea("AAPL")
        print(resp)
        self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])
        self.assertNotEqual(resp['output']['rsym'], None)


