'''
한국투자증권 python wrapper
'''
import datetime
import json
import os
import pickle
import threading
import time
import zipfile
from collections import deque, defaultdict
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from typing import Literal
from zoneinfo import ZoneInfo  # Requires Python 3.9+

import pandas as pd
import requests

EXCHANGE_CODE = {
    "홍콩": "HKS",
    "뉴욕": "NYS",
    "나스닥": "NAS",
    "아멕스": "AMS",
    "도쿄": "TSE",
    "상해": "SHS",
    "심천": "SZS",
    "상해지수": "SHI",
    "심천지수": "SZI",
    "호치민": "HSX",
    "하노이": "HNX"
}

# 해외주식 주문
# 해외주식 잔고
EXCHANGE_CODE2 = {
    "미국전체": "NASD",
    "나스닥": "NAS",
    "뉴욕": "NYSE",
    "아멕스": "AMEX",
    "홍콩": "SEHK",
    "상해": "SHAA",
    "심천": "SZAA",
    "도쿄": "TKSE",
    "하노이": "HASE",
    "호치민": "VNSE"
}

EXCHANGE_CODE3 = {
    "나스닥": "NASD",
    "뉴욕": "NYSE",
    "아멕스": "AMEX",
    "홍콩": "SEHK",
    "상해": "SHAA",
    "심천": "SZAA",
    "도쿄": "TKSE",
    "하노이": "HASE",
    "호치민": "VNSE"
}

EXCHANGE_CODE4 = {
    "나스닥": "NAS",
    "뉴욕": "NYS",
    "아멕스": "AMS",
    "홍콩": "HKS",
    "상해": "SHS",
    "심천": "SZS",
    "도쿄": "TSE",
    "하노이": "HNX",
    "호치민": "HSX",
    "상해지수": "SHI",
    "심천지수": "SZI"
}

CURRENCY_CODE = {
    "나스닥": "USD",
    "뉴욕": "USD",
    "아멕스": "USD",
    "홍콩": "HKD",
    "상해": "CNY",
    "심천": "CNY",
    "도쿄": "JPY",
    "하노이": "VND",
    "호치민": "VND"
}

MARKET_TYPE_MAP = {
    "KR": ["300"],  # "301", "302"
    "KRX": ["300"],  # "301", "302"
    "NASDAQ": ["512"],
    "NYSE": ["513"],
    "AMEX": ["529"],
    "US": ["512", "513", "529"],
    "TYO": ["515"],
    "JP": ["515"],
    "HKEX": ["501"],
    "HK": ["501", "543", "558"],
    "HNX": ["507"],
    "HSX": ["508"],
    "VN": ["507", "508"],
    "SSE": ["551"],
    "SZSE": ["552"],
    "CN": ["551", "552"]
}

MARKET_TYPE = Literal[
    "KRX",
    "NASDAQ",
    "NYSE",
    "AMEX",
    "TYO",
    "HKEX",
    "HNX",
    "HSX",
    "SSE",
    "SZSE",
]

EXCHANGE_TYPE = Literal[
    "NAS",
    "NYS",
    "AMS"
]

MARKET_CODE_MAP: dict[str, MARKET_TYPE] = {
    "300": "KRX",
    "301": "KRX",
    "302": "KRX",
    "512": "NASDAQ",
    "513": "NYSE",
    "529": "AMEX",
    "515": "TYO",
    "501": "HKEX",
    "543": "HKEX",
    "558": "HKEX",
    "507": "HNX",
    "508": "HSX",
    "551": "SSE",
    "552": "SZSE",
}

EXCHANGE_CODE_MAP: dict[str, EXCHANGE_TYPE] = {
    "NASDAQ": "NAS",
    "NYSE": "NYS",
    "AMEX": "AMS"
}

API_RETURN_CODE = {
    "SUCCESS": "0",  # 조회되었습니다
    "EXPIRED_TOKEN": "1",  # 기간이 만료된 token 입니다
    "NO_DATA": "7",  # 조회할 자료가 없습니다
}


class KoreaInvestment:
    '''
    한국투자증권 REST API
    '''

    def __init__(self, api_key: str, api_secret: str, acc_no: str,
                 mock: bool = False):
        """생성자
        Args:
            api_key (str): 발급받은 API key
            api_secret (str): 발급받은 API secret
            acc_no (str): 계좌번호 체계의 앞 8자리-뒤 2자리
            exchange (str): "서울", "나스닥", "뉴욕", "아멕스", "홍콩", "상해", "심천", # todo: exchange는 제거 예정
                            "도쿄", "하노이", "호치민"
            mock (bool): True (mock trading), False (real trading)
        """
        self.mock = mock
        self.set_base_url(mock)
        self.api_key = api_key
        self.api_secret = api_secret

        # account number
        self.acc_no = acc_no
        self.acc_no_prefix = acc_no.split('-')[0]
        self.acc_no_postfix = acc_no.split('-')[1]
        max_calls = 20

        self.rate_limiter = RateLimiter(max_calls, 1)
        self.executor = ThreadPoolExecutor(max_workers=max_calls)

        # access token
        self.token_file = Path("~/.cache/mojito2/token.dat").expanduser()
        self.access_token = None
        if self.check_access_token():
            self.load_access_token()
        else:
            self.issue_access_token()

    def __execute_concurrent_requests(self, method, stock_list):
        futures = [self.executor.submit(method, symbol_id, market) for symbol_id, market in stock_list]
        results = [future.result() for future in futures]

        self.rate_limiter.print_stats()

        return results

    def shutdown(self):
        self.executor.shutdown(wait=True)

    def set_base_url(self, mock: bool = True):
        """테스트(모의투자) 서버 사용 설정
        Args:
            mock(bool, optional): True: 테스트서버, False: 실서버 Defaults to True.
        """
        if mock:
            self.base_url = "https://openapivts.koreainvestment.com:29443"
        else:
            self.base_url = "https://openapi.koreainvestment.com:9443"

    def issue_access_token(self):
        """OAuth인증/접근토큰발급
        """
        path = "oauth2/tokenP"
        url = f"{self.base_url}/{path}"
        headers = {"content-type": "application/json"}
        data = {
            "grant_type": "client_credentials",
            "appkey": self.api_key,
            "appsecret": self.api_secret
        }

        resp = requests.post(url, headers=headers, json=data)
        resp_data = resp.json()
        self.access_token = f'Bearer {resp_data["access_token"]}'

        # 'expires_in' has no reference time and causes trouble:
        # The server thinks I'm expired but my token.dat looks still valid!
        # Hence, we use 'access_token_token_expired' here.
        # This error is quite big. I've seen 4000 seconds.
        timezone = ZoneInfo('Asia/Seoul')
        dt = datetime.datetime.strptime(resp_data['access_token_token_expired'], '%Y-%m-%d %H:%M:%S').replace(
            tzinfo=timezone)
        resp_data['timestamp'] = int(dt.timestamp())
        resp_data['api_key'] = self.api_key
        resp_data['api_secret'] = self.api_secret

        # dump access token
        self.token_file.parent.mkdir(parents=True, exist_ok=True)
        with self.token_file.open("wb") as f:
            pickle.dump(resp_data, f)

    def check_access_token(self) -> bool:
        """check access token

        Returns:
            Bool: True: token is valid, False: token is not valid
        """

        if not self.token_file.exists():
            return False

        with self.token_file.open("rb") as f:
            data = pickle.load(f)

        expire_epoch = data['timestamp']
        now_epoch = int(datetime.datetime.now().timestamp())
        status = False

        if (data['api_key'] != self.api_key) or (data['api_secret'] != self.api_secret):
            return False

        good_until = data['timestamp']
        ts_now = int(datetime.datetime.now().timestamp())
        return ts_now < good_until

    def load_access_token(self):
        """load access token
        """
        with self.token_file.open("rb") as f:
            data = pickle.load(f)
        self.access_token = f'Bearer {data["access_token"]}'

    def issue_hashkey(self, data: dict):
        """해쉬키 발급
        Args:
            data (dict): POST 요청 데이터
        Returns:
            _type_: _description_
        """
        path = "uapi/hashkey"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "User-Agent": "Mozilla/5.0"
        }
        resp = requests.post(url, headers=headers, data=json.dumps(data))
        haskkey = resp.json()["HASH"]
        return haskkey

    def fetch_search_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests(self.__fetch_search_stock_info, stock_market_list)

    def fetch_price_list(self, stock_list):
        return self.__execute_concurrent_requests(self.__fetch_price, stock_list)

    def __fetch_price(self, symbol: str, market: str = "KR") -> dict:
        """국내주식시세/주식현재가 시세
           해외주식현재가/해외주식 현재체결가

        Args:
            symbol (str): 종목코드

        Returns:
            dict: _description_
        """
        # TODO: 시세 조회 메서드들이 삭제되어 이 메서드도 구현 변경 필요
        raise NotImplementedError("시세 조회 기능이 제거되었습니다.")

    def fetch_symbols(self):
        """fetch symbols from the exchange

        Returns:
            pd.DataFrame: pandas dataframe
        """
        if self.exchange == "서울":  # todo: exchange는 제거 예정
            df = self.fetch_kospi_symbols()
            kospi_df = df[['단축코드', '한글명', '그룹코드']].copy()
            kospi_df['시장'] = '코스피'

            df = self.fetch_kosdaq_symbols()
            kosdaq_df = df[['단축코드', '한글명', '그룹코드']].copy()
            kosdaq_df['시장'] = '코스닥'

            df = pd.concat([kospi_df, kosdaq_df], axis=0)

        return df

    def download_master_file(self, base_dir: str, file_name: str, url: str):
        """download master file

        Args:
            base_dir (str): download directory
            file_name (str: filename
            url (str): url
        """
        os.chdir(base_dir)

        # delete legacy master file
        if os.path.exists(file_name):
            os.remove(file_name)

        # download master file
        resp = requests.get(url)
        with open(file_name, "wb") as f:
            f.write(resp.content)

        # unzip
        kospi_zip = zipfile.ZipFile(file_name)
        kospi_zip.extractall()
        kospi_zip.close()

    def parse_kospi_master(self, base_dir: str):
        """parse kospi master file

        Args:
            base_dir (str): directory where kospi code exists

        Returns:
            _type_: _description_
        """
        file_name = base_dir + "/kospi_code.mst"
        tmp_fil1 = base_dir + "/kospi_code_part1.tmp"
        tmp_fil2 = base_dir + "/kospi_code_part2.tmp"

        wf1 = open(tmp_fil1, mode="w", encoding="cp949")
        wf2 = open(tmp_fil2, mode="w")

        with open(file_name, mode="r", encoding="cp949") as f:
            for row in f:
                rf1 = row[0:len(row) - 228]
                rf1_1 = rf1[0:9].rstrip()
                rf1_2 = rf1[9:21].rstrip()
                rf1_3 = rf1[21:].strip()
                wf1.write(rf1_1 + ',' + rf1_2 + ',' + rf1_3 + '\n')
                rf2 = row[-228:]
                wf2.write(rf2)

        wf1.close()
        wf2.close()

        part1_columns = ['단축코드', '표준코드', '한글명']
        df1 = pd.read_csv(tmp_fil1, header=None, encoding='cp949', names=part1_columns)

        field_specs = [
            2, 1, 4, 4, 4,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 1, 1, 1, 1,
            1, 9, 5, 5, 1,
            1, 1, 2, 1, 1,
            1, 2, 2, 2, 3,
            1, 3, 12, 12, 8,
            15, 21, 2, 7, 1,
            1, 1, 1, 1, 9,
            9, 9, 5, 9, 8,
            9, 3, 1, 1, 1
        ]

        part2_columns = [
            '그룹코드', '시가총액규모', '지수업종대분류', '지수업종중분류', '지수업종소분류',
            '제조업', '저유동성', '지배구조지수종목', 'KOSPI200섹터업종', 'KOSPI100',
            'KOSPI50', 'KRX', 'ETP', 'ELW발행', 'KRX100',
            'KRX자동차', 'KRX반도체', 'KRX바이오', 'KRX은행', 'SPAC',
            'KRX에너지화학', 'KRX철강', '단기과열', 'KRX미디어통신', 'KRX건설',
            'Non1', 'KRX증권', 'KRX선박', 'KRX섹터_보험', 'KRX섹터_운송',
            'SRI', '기준가', '매매수량단위', '시간외수량단위', '거래정지',
            '정리매매', '관리종목', '시장경고', '경고예고', '불성실공시',
            '우회상장', '락구분', '액면변경', '증자구분', '증거금비율',
            '신용가능', '신용기간', '전일거래량', '액면가', '상장일자',
            '상장주수', '자본금', '결산월', '공모가', '우선주',
            '공매도과열', '이상급등', 'KRX300', 'KOSPI', '매출액',
            '영업이익', '경상이익', '당기순이익', 'ROE', '기준년월',
            '시가총액', '그룹사코드', '회사신용한도초과', '담보대출가능', '대주가능'
        ]

        df2 = pd.read_fwf(tmp_fil2, widths=field_specs, names=part2_columns)
        df = pd.merge(df1, df2, how='outer', left_index=True, right_index=True)

        # clean temporary file and dataframe
        del (df1)
        del (df2)
        os.remove(tmp_fil1)
        os.remove(tmp_fil2)
        return df

    def parse_kosdaq_master(self, base_dir: str):
        """parse kosdaq master file

        Args:
            base_dir (str): directory where kosdaq code exists

        Returns:
            _type_: _description_
        """
        file_name = base_dir + "/kosdaq_code.mst"
        tmp_fil1 = base_dir + "/kosdaq_code_part1.tmp"
        tmp_fil2 = base_dir + "/kosdaq_code_part2.tmp"

        wf1 = open(tmp_fil1, mode="w", encoding="cp949")
        wf2 = open(tmp_fil2, mode="w")
        with open(file_name, mode="r", encoding="cp949") as f:
            for row in f:
                rf1 = row[0:len(row) - 222]
                rf1_1 = rf1[0:9].rstrip()
                rf1_2 = rf1[9:21].rstrip()
                rf1_3 = rf1[21:].strip()
                wf1.write(rf1_1 + ',' + rf1_2 + ',' + rf1_3 + '\n')

                rf2 = row[-222:]
                wf2.write(rf2)

        wf1.close()
        wf2.close()

        part1_columns = ['단축코드', '표준코드', '한글명']
        df1 = pd.read_csv(tmp_fil1, header=None, encoding="cp949", names=part1_columns)

        field_specs = [
            2, 1, 4, 4, 4,  # line 20
            1, 1, 1, 1, 1,  # line 27
            1, 1, 1, 1, 1,  # line 32
            1, 1, 1, 1, 1,  # line 38
            1, 1, 1, 1, 1,  # line 43
            1, 9, 5, 5, 1,  # line 48
            1, 1, 2, 1, 1,  # line 54
            1, 2, 2, 2, 3,  # line 64
            1, 3, 12, 12, 8,  # line 69
            15, 21, 2, 7, 1,  # line 75
            1, 1, 1, 9, 9,  # line 80
            9, 5, 9, 8, 9,  # line 85
            3, 1, 1, 1
        ]

        part2_columns = [
            '그룹코드', '시가총액규모', '지수업종대분류', '지수업종중분류', '지수업종소분류',  # line 20
            '벤처기업', '저유동성', 'KRX', 'ETP', 'KRX100',  # line 27
            'KRX자동차', 'KRX반도체', 'KRX바이오', 'KRX은행', 'SPAC',  # line 32
            'KRX에너지화학', 'KRX철강', '단기과열', 'KRX미디어통신', 'KRX건설',  # line 38
            '투자주의', 'KRX증권', 'KRX선박', 'KRX섹터_보험', 'KRX섹터_운송',  # line 43
            'KOSDAQ150', '기준가', '매매수량단위', '시간외수량단위', '거래정지',  # line 48
            '정리매매', '관리종목', '시장경고', '경고예고', '불성실공시',  # line 54
            '우회상장', '락구분', '액면변경', '증자구분', '증거금비율',  # line 64
            '신용가능', '신용기간', '전일거래량', '액면가', '상장일자',  # line 69
            '상장주수', '자본금', '결산월', '공모가', '우선주',  # line 75
            '공매도과열', '이상급등', 'KRX300', '매출액', '영업이익',  # line 80
            '경상이익', '당기순이익', 'ROE', '기준년월', '시가총액',  # line 85
            '그룹사코드', '회사신용한도초과', '담보대출가능', '대주가능'
        ]

        df2 = pd.read_fwf(tmp_fil2, widths=field_specs, names=part2_columns)
        df = pd.merge(df1, df2, how='outer', left_index=True, right_index=True)

        # clean temporary file and dataframe
        del (df1)
        del (df2)
        os.remove(tmp_fil1)
        os.remove(tmp_fil2)
        return df

    def fetch_kospi_symbols(self):
        """코스피 종목 코드

        실제 필요한 종목: ST, RT, EF, IF

        ST	주권
        MF	증권투자회사
        RT	부동산투자회사
        SC	선박투자회사
        IF	사회간접자본투융자회사
        DR	주식예탁증서
        EW	ELW
        EF	ETF
        SW	신주인수권증권
        SR	신주인수권증서
        BC	수익증권
        FE	해외ETF
        FS	외국주권


        Returns:
            DataFrame:
        """
        base_dir = os.getcwd()
        file_name = "kospi_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name
        self.download_master_file(base_dir, file_name, url)
        df = self.parse_kospi_master(base_dir)
        return df

    def fetch_kosdaq_symbols(self):
        """코스닥 종목 코드

        Returns:
            DataFrame:
        """
        base_dir = os.getcwd()
        file_name = "kosdaq_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name
        self.download_master_file(base_dir, file_name, url)
        df = self.parse_kosdaq_master(base_dir)
        return df

    def fetch_price_detail_oversea_list(self, stock_market_list):
        return self.__execute_concurrent_requests(self.__fetch_price_detail_oversea, stock_market_list)

    def __fetch_price_detail_oversea(self, symbol: str, market: str = "KR"):
        """해외주식 현재가상세

        Args:
            symbol (str): symbol
        """
        self.rate_limiter.acquire()

        path = "/uapi/overseas-price/v1/quotations/price-detail"
        url = f"{self.base_url}/{path}"

        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "HHDFS76200200"
        }

        if market == "KR" or market == "KRX":
            # API 호출해서 실제로 확인은 못해봄, overasea 이라서 안될 것으로 판단해서 조건문 추가함
            raise ValueError("Market cannot be either 'KR' or 'KRX'.")

        for market_code in MARKET_TYPE_MAP[market]:
            print("market_code", market_code)
            market_type = MARKET_CODE_MAP[market_code]
            exchange_code = EXCHANGE_CODE_MAP[market_type]
            params = {
                "AUTH": "",
                "EXCD": exchange_code,
                "SYMB": symbol
            }
            resp = requests.get(url, headers=headers, params=params)
            resp_json = resp.json()
            if resp_json['rt_cd'] != API_RETURN_CODE["SUCCESS"] or resp_json['output']['rsym'] == '':
                continue

            return resp_json

    def fetch_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests(self.__fetch_stock_info, stock_market_list)

    def __fetch_stock_info(self, symbol: str, market: str = "KR"):
        self.rate_limiter.acquire()

        path = "uapi/domestic-stock/v1/quotations/search-info"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "CTPF1604R"
        }

        for market_code in MARKET_TYPE_MAP[market]:
            try:
                params = {
                    "PDNO": symbol,
                    "PRDT_TYPE_CD": market_code
                }
                resp = requests.get(url, headers=headers, params=params)
                resp_json = resp.json()

                if resp_json['rt_cd'] == API_RETURN_CODE['NO_DATA']:
                    continue
                return resp_json

            except Exception as e:
                print(e)
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e

    def fetch_search_stock_info_list(self, stock_market_list):
        return self.__execute_concurrent_requests(self.__fetch_search_stock_info, stock_market_list)

    def __fetch_search_stock_info(self, symbol: str, market: str = "KR"):
        """
        국내 주식만 제공하는 API이다
        """

        self.rate_limiter.acquire()

        path = "uapi/domestic-stock/v1/quotations/search-stock-info"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "CTPF1002R"
        }

        if market != "KR" and market != "KRX":
            raise ValueError("Market must be either 'KR' or 'KRX'.")

        for market_ in MARKET_TYPE_MAP[market]:
            try:
                params = {
                    "PDNO": symbol,
                    "PRDT_TYPE_CD": market_
                }
                resp = requests.get(url, headers=headers, params=params)
                resp_json = resp.json()

                if resp_json['rt_cd'] == API_RETURN_CODE['NO_DATA']:
                    continue
                return resp_json

            except Exception as e:
                print(e)
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e


class RateLimiter:
    def __init__(self, max_calls, per_seconds):
        self.max_calls = max_calls
        self.per_seconds = per_seconds
        self.lock = threading.Lock()
        self.call_timestamps = deque()

        # 호출 통계 추적을 위한 딕셔너리
        self.calls_per_second = defaultdict(int)

    def acquire(self):
        with self.lock:
            now = time.time()

            # 만료된 타임스탬프 제거
            while self.call_timestamps and self.call_timestamps[0] <= now - self.per_seconds:
                self.call_timestamps.popleft()

            # 호출 제한에 도달했다면 대기
            if len(self.call_timestamps) >= self.max_calls:
                wait_time = self.per_seconds - (now - self.call_timestamps[0])
                if wait_time > 0:
                    time.sleep(wait_time)
                    now = time.time()  # 대기 후 시간 업데이트

            # 호출 기록
            current_second = int(now)
            self.calls_per_second[current_second] += 1
            self.call_timestamps.append(now)

    def get_stats(self):
        """호출 통계 분석 결과 반환"""
        stats = {
            "calls_per_second": dict(self.calls_per_second),
            "max_calls_in_one_second": max(self.calls_per_second.values()) if self.calls_per_second else 0,
            "total_calls": sum(self.calls_per_second.values()),
            "seconds_tracked": len(self.calls_per_second)
        }
        return stats

    def print_stats(self):
        """호출 통계 출력"""
        if not self.calls_per_second:
            print("호출 데이터가 없습니다.")
            return

        print("\n===== 초당 API 호출 횟수 분석 =====")
        max_calls = max(self.calls_per_second.values())

        for second, count in sorted(self.calls_per_second.items()):
            timestamp = datetime.datetime.fromtimestamp(second).strftime('%H:%M:%S')
            print(f"시간: {timestamp}, 호출 수: {count}")

        print(f"\n최대 초당 호출 횟수: {max_calls}")
        print(f"설정된 max_calls: {self.max_calls}")
        print(f"제한 준수 여부: {'준수' if max_calls <= self.max_calls else '초과'}")
        print(f"총 호출 횟수: {sum(self.calls_per_second.values())}")
        print("================================\n")


if __name__ == "__main__":
    with open("../koreainvestment.key", encoding='utf-8') as key_file:
        lines = key_file.readlines()

    key = lines[0].strip()
    secret = lines[1].strip()
    acc_no = lines[2].strip()

    broker = KoreaInvestment(
        api_key=key,
        api_secret=secret,
        acc_no=acc_no,
        # exchange="나스닥" # todo: exchange는 제거 예정
    )

    balance = broker.fetch_present_balance()
    print(balance)

    # result = broker.fetch_oversea_day_night()
    # pprint.pprint(result)

    # minute1_ohlcv = broker.fetch_today_1m_ohlcv("005930")
    # pprint.pprint(minute1_ohlcv)

    # broker = KoreaInvestment(key, secret, exchange="나스닥")
    # import pprint
    # resp = broker.fetch_price("005930")
    # pprint.pprint(resp)
    #
    # b = broker.fetch_balance("63398082")
    # pprint.pprint(b)
    #
    # resp = broker.create_market_buy_order("63398082", "005930", 10)
    # pprint.pprint(resp)
    #
    # resp = broker.cancel_order("63398082", "91252", "0000117057", "00", 60000, 5, "Y")
    # print(resp)
    #
    # resp = broker.create_limit_buy_order("63398082", "TQQQ", 35, 1)
    # print(resp)



    # import pprint
    # broker = KoreaInvestment(key, secret, exchange="나스닥")
    # resp_ohlcv = broker.fetch_ohlcv("TSLA", '1d', to="")
    # print(len(resp_ohlcv['output2']))
    # pprint.pprint(resp_ohlcv['output2'][0])
    # pprint.pprint(resp_ohlcv['output2'][-1])
