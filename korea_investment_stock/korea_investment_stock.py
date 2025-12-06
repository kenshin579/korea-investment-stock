'''
한국투자증권 python wrapper
'''
import json
import os
import zipfile
import logging
import re
from pathlib import Path
from typing import Optional
from zoneinfo import ZoneInfo  # Requires Python 3.9+
from datetime import datetime, timedelta

import pandas as pd
import requests

from .token_storage import TokenStorage, FileTokenStorage, RedisTokenStorage
from .constants import (
    MARKET_TYPE_MAP,
    API_RETURN_CODE,
    FID_COND_MRKT_DIV_CODE_STOCK,
)
from .config_resolver import ConfigResolver
from .parsers import parse_kospi_master, parse_kosdaq_master
from .ipo import (
    validate_date_format,
    validate_date_range,
    parse_ipo_date_range,
    format_ipo_date,
    calculate_ipo_d_day,
    get_ipo_status,
    format_number,
)

# 로거 설정
logger = logging.getLogger(__name__)


class KoreaInvestment:
    '''
    한국투자증권 REST API
    '''

    # 기본 캐시 TTL (시간) - 1주일
    DEFAULT_MASTER_TTL_HOURS = 168

    def __init__(
        self,
        api_key: str | None = None,
        api_secret: str | None = None,
        acc_no: str | None = None,
        config: "Config | None" = None,
        config_file: "str | Path | None" = None,
        token_storage: Optional[TokenStorage] = None
    ):
        """한국투자증권 API 클라이언트 초기화

        설정 우선순위 (5단계):
            1. 생성자 파라미터 (최고 우선순위)
            2. config 객체
            3. config_file 파라미터
            4. 환경 변수 (KOREA_INVESTMENT_*)
            5. 기본 config 파일 (~/.config/kis/config.yaml)

        Args:
            api_key (str | None): 발급받은 API key
            api_secret (str | None): 발급받은 API secret
            acc_no (str | None): 계좌번호 체계의 앞 8자리-뒤 2자리 (예: "12345678-01")
            config (Config | None): Config 객체 (Phase 2에서 추가됨)
            config_file (str | Path | None): 설정 파일 경로
            token_storage (Optional[TokenStorage]): 토큰 저장소 인스턴스

        Raises:
            ValueError: api_key, api_secret, 또는 acc_no가 설정되지 않았을 때
            ValueError: acc_no 형식이 올바르지 않을 때

        Examples:
            # 방법 1: 생성자 파라미터 (기존 방식)
            >>> broker = KoreaInvestment(
            ...     api_key="your-api-key",
            ...     api_secret="your-api-secret",
            ...     acc_no="12345678-01"
            ... )

            # 방법 2: 환경 변수 자동 감지
            >>> broker = KoreaInvestment()

            # 방법 3: Config 객체 사용
            >>> config = Config.from_yaml("~/.config/kis/config.yaml")
            >>> broker = KoreaInvestment(config=config)

            # 방법 4: config_file 파라미터
            >>> broker = KoreaInvestment(config_file="./my_config.yaml")

            # 방법 5: 혼합 사용 (일부만 override)
            >>> broker = KoreaInvestment(config=config, api_key="override-key")
        """
        # 5단계 우선순위로 설정 해결 (ConfigResolver 사용)
        resolver = ConfigResolver()
        resolved = resolver.resolve(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
            config=config,
            config_file=config_file,
        )

        self.api_key = resolved["api_key"]
        self.api_secret = resolved["api_secret"]
        acc_no = resolved["acc_no"]

        # 필수값 검증
        missing_fields = []
        if not self.api_key:
            missing_fields.append("api_key (KOREA_INVESTMENT_API_KEY)")
        if not self.api_secret:
            missing_fields.append("api_secret (KOREA_INVESTMENT_API_SECRET)")
        if not acc_no:
            missing_fields.append("acc_no (KOREA_INVESTMENT_ACCOUNT_NO)")

        if missing_fields:
            raise ValueError(
                "API credentials required. Missing: " + ", ".join(missing_fields) + ". "
                "Pass as parameters, use config/config_file, or set KOREA_INVESTMENT_* environment variables."
            )

        # 계좌번호 형식 검증
        if '-' not in acc_no:
            raise ValueError(f"계좌번호 형식이 올바르지 않습니다. '12345678-01' 형식이어야 합니다. 입력값: {acc_no}")

        self.base_url = "https://openapi.koreainvestment.com:9443"

        # account number - 검증 후 split
        parts = acc_no.split('-')
        if len(parts) != 2 or len(parts[0]) != 8 or len(parts[1]) != 2:
            raise ValueError(f"계좌번호 형식이 올바르지 않습니다. 앞 8자리-뒤 2자리여야 합니다. 입력값: {acc_no}")

        self.acc_no = acc_no
        self.acc_no_prefix = parts[0]
        self.acc_no_postfix = parts[1]

        # resolved에서 token_storage 관련 설정 가져오기
        self._resolved_config = resolved

        # 토큰 저장소 초기화
        if token_storage:
            self.token_storage = token_storage
        else:
            self.token_storage = self._create_token_storage()

        # access token
        self.access_token = None
        if self.token_storage.check_token_valid(self.api_key, self.api_secret):
            token_data = self.token_storage.load_token(self.api_key, self.api_secret)
            if token_data:
                self.access_token = f'Bearer {token_data["access_token"]}'
        else:
            self.issue_access_token()

    def _create_token_storage(self) -> TokenStorage:
        """설정 기반 토큰 저장소 생성

        _resolved_config에서 설정을 읽어 토큰 저장소를 생성합니다.
        설정이 없으면 환경 변수에서 읽습니다.

        Returns:
            TokenStorage: 설정된 토큰 저장소 인스턴스

        Raises:
            ValueError: 지원하지 않는 저장소 타입일 때
        """
        # _resolved_config가 있으면 사용, 없으면 환경 변수에서 읽기
        if hasattr(self, "_resolved_config") and self._resolved_config:
            storage_type = self._resolved_config.get("token_storage_type") or "file"
            redis_url = self._resolved_config.get("redis_url") or "redis://localhost:6379/0"
            redis_password = self._resolved_config.get("redis_password")
            token_file = self._resolved_config.get("token_file")
        else:
            # 하위 호환성: 환경 변수에서 읽기
            storage_type = os.getenv("KOREA_INVESTMENT_TOKEN_STORAGE", "file")
            redis_url = os.getenv("KOREA_INVESTMENT_REDIS_URL", "redis://localhost:6379/0")
            redis_password = os.getenv("KOREA_INVESTMENT_REDIS_PASSWORD")
            token_file = os.getenv("KOREA_INVESTMENT_TOKEN_FILE")

        storage_type = storage_type.lower()

        if storage_type == "file":
            file_path = None
            if token_file:
                file_path = Path(token_file).expanduser()
            return FileTokenStorage(file_path)

        elif storage_type == "redis":
            return RedisTokenStorage(redis_url, password=redis_password)

        else:
            raise ValueError(
                f"지원하지 않는 저장소 타입: {storage_type}\n"
                f"'file' 또는 'redis'만 지원됩니다."
            )

    def __enter__(self):
        """컨텍스트 매니저 진입"""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """컨텍스트 매니저 종료 - 리소스 정리"""
        self.shutdown()
        return False  # 예외를 전파

    def shutdown(self):
        """리소스 정리"""
        # 컨텍스트 매니저 종료 시 호출됨
        # 향후 필요한 정리 작업이 있으면 여기에 추가
        pass


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
        dt = datetime.strptime(resp_data['access_token_token_expired'], '%Y-%m-%d %H:%M:%S').replace(
            tzinfo=timezone)
        resp_data['timestamp'] = int(dt.timestamp())
        resp_data['api_key'] = self.api_key
        resp_data['api_secret'] = self.api_secret

        # 토큰 저장소에 저장
        self.token_storage.save_token(resp_data)

    def check_access_token(self) -> bool:
        """check access token

        Returns:
            Bool: True: token is valid, False: token is not valid
        """
        return self.token_storage.check_token_valid(self.api_key, self.api_secret)

    def load_access_token(self):
        """load access token
        """
        token_data = self.token_storage.load_token(self.api_key, self.api_secret)
        if token_data:
            self.access_token = f'Bearer {token_data["access_token"]}'

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

    def fetch_price(self, symbol: str, market: str = "KR") -> dict:
        """국내주식시세/주식현재가 시세
           해외주식현재가/해외주식 현재체결가

        Args:
            symbol (str): 종목코드
            market (str): 시장 코드 ("KR", "KRX", "US" 등)

        Returns:
            dict: API 응답 데이터
        """

        if market == "KR" or market == "KRX":
            stock_info = self.fetch_stock_info(symbol, market)
            symbol_type = self.get_symbol_type(stock_info)
            resp_json = self.fetch_domestic_price(symbol, symbol_type)
        elif market == "US":
            # 기존: resp_json = self.fetch_oversea_price(symbol)  # 메서드 없음
            # 개선: 이미 구현된 fetch_price_detail_oversea() 활용
            resp_json = self.fetch_price_detail_oversea(symbol, market)
            # 참고: 이 API는 현재가 외에도 PER, PBR, EPS, BPS 등 추가 정보 제공
        else:
            raise ValueError("Unsupported market type")

        return resp_json

    def get_symbol_type(self, symbol_info):
        # API 오류 응답 처리
        if symbol_info.get('rt_cd') != '0' or 'output' not in symbol_info:
            return 'Stock'  # 기본값으로 주식 타입 반환

        symbol_type = symbol_info['output']['prdt_clsf_name']
        if symbol_type == '주권' or symbol_type == '상장REITS' or symbol_type == '사회간접자본투융자회사':
            return 'Stock'
        elif symbol_type == 'ETF':
            return 'ETF'

        return "Unknown"

    def fetch_domestic_price(
        self,
        symbol: str,
        symbol_type: str = "Stock"
    ) -> dict:
        """국내 주식/ETF 현재가시세

        Args:
            symbol: 종목코드 (ex: 005930)
            symbol_type: 상품 타입 ("Stock" 또는 "ETF")

        Returns:
            dict: API 응답 데이터
        """
        TR_ID_MAP = {
            "Stock": "FHKST01010100",
            "ETF": "FHPST02400000"
        }

        path = "uapi/domestic-stock/v1/quotations/inquire-price"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": TR_ID_MAP.get(symbol_type, "FHKST01010100")
        }
        params = {
            "fid_cond_mrkt_div_code": FID_COND_MRKT_DIV_CODE_STOCK["KRX"],
            "fid_input_iscd": symbol
        }
        resp = requests.get(url, headers=headers, params=params)
        return resp.json()

    def fetch_kospi_symbols(
        self,
        ttl_hours: int = 168,
        force_download: bool = False
    ) -> pd.DataFrame:
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

        Args:
            ttl_hours (int): 캐시 유효 시간 (기본 1주일 = 168시간)
            force_download (bool): 강제 다운로드 여부

        Returns:
            DataFrame: 코스피 종목 정보
        """
        base_dir = os.getcwd()
        file_name = "kospi_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name

        self.download_master_file(base_dir, file_name, url, ttl_hours, force_download)
        df = parse_kospi_master(base_dir)
        return df

    def fetch_kosdaq_symbols(
        self,
        ttl_hours: int = 168,
        force_download: bool = False
    ) -> pd.DataFrame:
        """코스닥 종목 코드

        Args:
            ttl_hours (int): 캐시 유효 시간 (기본 1주일 = 168시간)
            force_download (bool): 강제 다운로드 여부

        Returns:
            DataFrame: 코스닥 종목 정보
        """
        base_dir = os.getcwd()
        file_name = "kosdaq_code.mst.zip"
        url = "https://new.real.download.dws.co.kr/common/master/" + file_name

        self.download_master_file(base_dir, file_name, url, ttl_hours, force_download)
        df = parse_kosdaq_master(base_dir)
        return df

    def _should_download(
        self,
        file_path: Path,
        ttl_hours: int,
        force: bool
    ) -> bool:
        """다운로드 필요 여부 판단

        Args:
            file_path: ZIP 파일 경로
            ttl_hours: 캐시 유효 시간
            force: 강제 다운로드 여부

        Returns:
            bool: True=다운로드 필요, False=캐시 사용
        """
        if force:
            return True

        if not file_path.exists():
            return True

        mtime = datetime.fromtimestamp(file_path.stat().st_mtime)
        age = datetime.now() - mtime

        if age.total_seconds() > ttl_hours * 3600:
            logger.debug(f"캐시 만료: {file_path} (age={age})")
            return True

        return False

    def download_master_file(
        self,
        base_dir: str,
        file_name: str,
        url: str,
        ttl_hours: int = 168,
        force_download: bool = False
    ) -> bool:
        """master 파일 다운로드 (캐싱 지원)

        Args:
            base_dir (str): 저장 디렉토리
            file_name (str): 파일명 (예: "kospi_code.mst.zip")
            url (str): 다운로드 URL
            ttl_hours (int): 캐시 유효 시간 (기본 1주일 = 168시간)
            force_download (bool): 강제 다운로드 여부

        Returns:
            bool: True=다운로드됨, False=캐시 사용
        """
        zip_path = Path(base_dir) / file_name

        # 다운로드 필요 여부 확인
        if not self._should_download(zip_path, ttl_hours, force_download):
            mtime = datetime.fromtimestamp(zip_path.stat().st_mtime)
            age_hours = (datetime.now() - mtime).total_seconds() / 3600
            logger.info(f"캐시 사용: {zip_path} (age: {age_hours:.1f}h, ttl: {ttl_hours}h)")
            return False

        # 다운로드
        logger.info(f"다운로드 중: {url} -> {zip_path}")
        resp = requests.get(url)
        resp.raise_for_status()

        with open(zip_path, "wb") as f:
            f.write(resp.content)

        # 압축 해제
        with zipfile.ZipFile(zip_path, 'r') as zf:
            zf.extractall(base_dir)

        return True

    def fetch_price_detail_oversea(self, symbol: str, market: str = "KR"):
        """해외주식 현재가상세

        Args:
            symbol (str): symbol
        """
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

        for exchange_code in ["NYS", "NAS", "AMS", "BAY", "BAQ", "BAA"]:
            logger.debug(f"exchange_code: {exchange_code}")
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
        
        # 모든 거래소에서 실패한 경우
        raise ValueError(f"Unable to fetch price for symbol '{symbol}' in any {market} exchange")

    def fetch_stock_info(self, symbol: str, market: str = "KR"):
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
                logger.debug(f"fetch_stock_info 에러: {e}")
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e

    def fetch_search_stock_info(self, symbol: str, market: str = "KR"):
        """
        국내 주식만 제공하는 API이다
        """
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
                logger.debug(f"fetch_search_stock_info 에러: {e}")
                if resp_json['rt_cd'] != API_RETURN_CODE['SUCCESS']:
                    continue
                raise e

    # IPO Schedule API
    def fetch_ipo_schedule(self, from_date: str = None, to_date: str = None, symbol: str = "") -> dict:
        """공모주 청약 일정 조회
        
        예탁원정보(공모주청약일정) API를 통해 공모주 정보를 조회합니다.
        한국투자 HTS(eFriend Plus) > [0667] 공모주청약 화면과 동일한 기능입니다.
        
        Args:
            from_date: 조회 시작일 (YYYYMMDD, 기본값: 오늘)
            to_date: 조회 종료일 (YYYYMMDD, 기본값: 30일 후)
            symbol: 종목코드 (선택, 공백시 전체 조회)
            
        Returns:
            dict: 공모주 청약 일정 정보
                {
                    "rt_cd": "0",  # 성공여부
                    "msg_cd": "응답코드",
                    "msg1": "응답메시지",
                    "output1": [
                        {
                            "record_date": "기준일",
                            "sht_cd": "종목코드",
                            "isin_name": "종목명",
                            "fix_subscr_pri": "공모가",
                            "face_value": "액면가",
                            "subscr_dt": "청약기간",  # "2024.01.15~2024.01.16"
                            "pay_dt": "납입일",
                            "refund_dt": "환불일",
                            "list_dt": "상장/등록일",
                            "lead_mgr": "주간사",
                            "pub_bf_cap": "공모전자본금",
                            "pub_af_cap": "공모후자본금",
                            "assign_stk_qty": "당사배정물량"
                        }
                    ]
                }
                
        Raises:
            ValueError: 날짜 형식 오류시

        Note:
            - 예탁원에서 제공한 자료이므로 정보용으로만 사용하시기 바랍니다.
            - 실제 청약시에는 반드시 공식 공모주 청약 공고문을 확인하세요.
            
        Examples:
            >>> # 전체 공모주 조회 (오늘부터 30일)
            >>> ipos = broker.fetch_ipo_schedule()
            
            >>> # 특정 기간 조회
            >>> ipos = broker.fetch_ipo_schedule(
            ...     from_date="20240101",
            ...     to_date="20240131"
            ... )
            
            >>> # 특정 종목 조회
            >>> ipo = broker.fetch_ipo_schedule(symbol="123456")
        """
        # 날짜 기본값 설정
        if not from_date:
            from_date = datetime.now().strftime("%Y%m%d")
        if not to_date:
            to_date = (datetime.now() + timedelta(days=30)).strftime("%Y%m%d")
        
        # 날짜 유효성 검증
        if not validate_date_format(from_date) or not validate_date_format(to_date):
            raise ValueError("날짜 형식은 YYYYMMDD 이어야 합니다.")
        
        if not validate_date_range(from_date, to_date):
            raise ValueError("시작일은 종료일보다 이전이어야 합니다.")
        
        path = "uapi/domestic-stock/v1/ksdinfo/pub-offer"
        url = f"{self.base_url}/{path}"
        headers = {
            "content-type": "application/json",
            "authorization": self.access_token,
            "appKey": self.api_key,
            "appSecret": self.api_secret,
            "tr_id": "HHKDB669108C0",
            "custtype": "P"  # 개인
        }
        
        params = {
            "SHT_CD": symbol,
            "CTS": "",
            "F_DT": from_date,
            "T_DT": to_date
        }
        
        resp = requests.get(url, headers=headers, params=params)
        resp_json = resp.json()
        
        # 에러 처리
        if resp_json.get('rt_cd') != '0':
            logger.error(f"공모주 조회 실패: {resp_json.get('msg1', 'Unknown error')}")
            return resp_json
        
        return resp_json
