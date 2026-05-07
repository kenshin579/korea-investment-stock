// Package domestic 은 한국투자증권 OpenAPI 의 국내주식 카테고리 메서드.
//
// Phase 1.2 메서드 (7):
//
//   - InquirePrice                 — 주식현재가 시세 (FHKST01010100)
//   - SearchInfo                   — 상품기본조회 (CTPF1604R)
//   - SearchStockInfo              — 주식기본조회 (CTPF1002R)
//   - InquireDailyItemChartPrice   — 국내주식기간별시세 일/주/월/년 (FHKST03010100)
//   - InquireTimeItemChartPrice    — 주식당일분봉조회 (FHKST03010200)
//   - FetchKospiSymbols            — KRX KOSPI 마스터 (한투 API 가 아닌 KRX 공개 다운로드)
//   - FetchKosdaqSymbols           — KRX KOSDAQ 마스터
//
// Phase 1.3 메서드 (9):
//
//   - InquireVolumeRank            — 거래량순위 (FHPST01710000)
//   - InquireFluctuation           — 등락률 순위 (FHPST01700000)
//   - InquireMarketCap             — 시가총액 상위 (FHPST01740000)
//   - InquireDividendRate          — 배당률 상위 (HHKDB13470100)
//   - InquireFinancialRatio        — 재무비율 (FHKST66430300)
//   - InquireIncomeStatement       — 손익계산서 (FHKST66430200)
//   - InquireBalanceSheet          — 대차대조표 (FHKST66430100)
//   - InquireProfitRatio           — 수익성비율 (FHKST66430400)
//   - InquireGrowthRatio           — 성장성비율 (FHKST66430800)
//
// Phase 1.4 메서드 (6):
//
//   - InquireInvestorTradeByStockDaily — 종목별 투자자매매동향 일별 (FHPTJ04160001)
//   - InquireInvestorDailyByMarket    — 시장별 투자자매매동향 일별 (FHPTJ04040000)
//   - InquireInvestorTimeByMarket     — 시장별 투자자매매동향 시세 (FHPTJ04030000)
//   - InquireIndexPrice               — 국내업종 현재지수 (FHPUP02100000)
//   - InquireIndexCategoryPrice       — 국내업종 구분별 전체시세 (FHPUP02140000)
//   - InquirePubOffer                 — 예탁원정보 공모주청약일정 (HHKDB669108C0)
//
// Phase 2.1 메서드 (3):
//
//   - InquireAskingPriceExpCcn  — 주식현재가 호가/예상체결 (FHKST01010200)
//   - InquireCcnl               — 주식현재가 체결 (FHKST01010300)
//   - InquireDailyPrice         — 주식현재가 일자별 (FHKST01010400)
//
// Phase 2.2 메서드 (5):
//
//   - InquireNearNewHighlow      — 국내주식 신고/신저근접종목 상위 (FHPST01870000)
//   - InquireOvertimePrice       — 국내주식 시간외현재가 (FHPST02300000)
//   - InquireOvertimeAskingPrice — 국내주식 시간외호가 (FHPST02300400)
//   - InquireOvertimeVolume      — 국내주식 시간외거래량순위 (FHPST02350000)
//   - InquireOvertimeFluctuation — 국내주식 시간외등락율순위 (FHPST02340000)
//
// Phase 2.4 메서드 (11):
//
//   - InquireKsdDividend    — 예탁원정보 배당일정 (HHKDB669102C0)
//   - InquireKsdBonusIssue  — 예탁원정보 무상증자 (HHKDB669101C0)
//   - InquireKsdPaidinCapin — 예탁원정보 유상증자 (HHKDB669100C0) [output key: output]
//   - InquireKsdSharehldMeet — 예탁원정보 주주총회 (HHKDB669111C0)
//   - InquireKsdMergerSplit  — 예탁원정보 합병/분할 (HHKDB669104C0) [no isin_name]
//   - InquireKsdRevSplit     — 예탁원정보 액면변경 (HHKDB669105C0) [+MARKET_GB]
//   - InquireKsdForfeit      — 예탁원정보 실권주청약 (HHKDB669109C0)
//   - InquireKsdMandDeposit  — 예탁원정보 의무보호예수 (HHKDB669110C0) [depo_date]
//   - InquireKsdCapDcrs      — 예탁원정보 감자 (HHKDB669106C0)
//   - InquireKsdPurreq       — 예탁원정보 주식매수청구 (HHKDB669103C0)
//   - InquireKsdListInfo     — 예탁원정보 주식상장정보 (HHKDB669107C0) [list_dt]
//
// Phase 2.5 — 투자자/매매 동향 (v1.8.0)
//
//	InquireInvestorTrendEstimate       HHPTJ04160200  투자자 매매 추정 가집계
//	InquireForeignInstitutionTotal     FHPTJ04400000  외인기관 매매종목가 집계
//	InquireProgramTradeByStockDaily    FHPPG04650201  종목별 프로그램매매 추이(일별)
//	InquireProgramTradeByStock         FHPPG04650101  종목별 프로그램매매 추이(체결)
//	InquireCompProgramTradeToday       FHPPG04600101  프로그램매매 종합현황(시간)
//	InquireCompProgramTradeDaily       FHPPG04600001  프로그램매매 종합현황(일별)
//	InquireInvestorProgramTradeToday   HHPPG046600C1  당일 투자자별 프로그램매매 동향
//
// Phase 2.7 — 업종/지수 (v1.10.0)  [Phase 2.5+ 마지막 sub-phase]
//
//	EP3  InquireIndexDailyPrice      — 국내업종 일자별지수       FHPUP02120000
//	EP4  InquireIndexTimeprice       — 국내업종 시간별지수 분    FHPUP02110200
//	EP5  InquireIndexTickprice       — 국내업종 시간별지수 초    FHPUP02110100
//	EP6  InquireDailyIndexchartprice — 국내주식업종기간별시세    FHKUP03500100
//	EP7  InquireTimeIndexchartprice  — 업종 분봉조회             FHKUP03500200
//	EP8  ExpTotalIndex               — 예상체결 전체지수         FHKUP11750000
//	EP9  ExpIndexTrend               — 예상체결지수 추이         FHPST01840000
//
// Anomalies:
//
//	EP1+EP2 already in Phase 1.4 → Phase 2.7 = 7 NEW (not 9)
//	EP8 lowercase fid_* query params (KIS 유일 예외)
//	EP8/EP9 prdy_ctrt short form (NOT bstp_nmix_prdy_ctrt)
//	EP9 KIS docs Korean labels scrambled — field names are correct
//
// Phase 4.1 — 종목정보/분석 (v1.12.0)
//
// Opinion 3:
//
//	EP1  InquireInvestOpinion    — 종목투자의견        FHKST663300C0
//	EP2  InquireInvestOpbysec    — 증권사별투자의견    FHKST663400C0
//	EP3  InquireEstimatePerform  — 종목추정실적        HHKST668300C0
//
// Extended 7:
//
//	EP4  InquireVolumePower      — 체결강도상위        FHPST01680000
//	EP5  InquireBulkTransNum     — 대량체결건수상위    FHKST190900C0
//	EP6  InquireTradprtByamt     — 체결금액별매매비중  FHKST111900C0
//	EP7  InquireHtsTopView       — HTS조회상위20종목   HHMCM000100C0
//	EP8  InquirePbarTraRatio     — 매물대거래비중      FHPST01130000
//	EP9  InquireExpPriceTrend    — 예상체결가추이      FHPST01810000
//	EP10 InquireExpTransUpdown   — 예상체결상승/하락상위 FHPST01820000
//
// Anomalies (Phase 4.1):
//
//	EP3 quad-output (output1/2/3/4), SHT_CD non-FID param name, KIS docs 오표기
//	EP4/EP5/EP9/EP10 lowercase fid_* query params
//	EP5 mksc_shrn_iscd (시장구분 포함) ≠ stck_shrn_iscd
//	EP6 whol_shun_vol_rate KIS wire typo (shun ≠ shnu) — 그대로 보존
//	EP7 zero params endpoint
//
// Phase 4.2 — 시장운영/특수상태 (v1.13.0)
//
//	EP4  InquireExpClosingPrice    — 장마감 예상체결가   FHKST117300C0
//	EP5  InquireChkHoliday         — 휴장일 조회         CTCA0903R
//	EP6  InquireViStatus           — 변동성완화장치 현황 FHPST01390000
//	EP7  InquireCaptureUplowprice  — 상하한가 포착       FHKST130000C0
//
// Anomalies (Phase 4.2):
//
//	EP4 output1 (not output) array, FID_INPUT_ISCD=시장구분코드 (종목코드 아님)
//	EP4 FID_COND_SCR_DIV_CODE="11173" hardcoded
//	EP5 non-FID UPPERCASE params (BASS_DT/CTX_AREA_NK/CTX_AREA_FK), CTCA prefix TR_ID
//	EP5 단시간 다수 호출 자제 (1일 1회 권장)
//	EP6 FID_COND_SCR_DIV_CODE="20139" hardcoded, output {} 단일 Object (runtime 배열 가능)
//	EP7 FID_COND_SCR_DIV_CODE="11300" hardcoded
//	WebSocket 제외: H0STMKO0/H0NXMKO0/H0UNMKO0 (장운영정보 KRX/NXT/통합) → Phase 5 이연
//
// Phase 4.3 — ranking/흐름 (v1.14.0)  [Phase 4 완료]
//
//	EP1   InquireShortSale              — 공매도 상위          FHPST04820000
//	EP2   InquireDailyShortSale         — 공매도 일별추이       FHPST04830000
//	EP3   InquireCreditBalance          — 신용잔고 상위         FHKST17010000
//	EP4   InquireDailyCreditBalance     — 신용잔고 일별추이     FHPST04760000
//	EP5   InquireLendableByCompany      — 당사 대주가능         CTSC2702R
//	EP6   InquireQuoteBalance           — 호가잔량 순위         FHPST01720000
//	EP7   InquireAfterHourBalance       — 시간외잔량 순위       FHPST01760000
//	EP8   InquireOvertimeExpTransFluct  — 시간외 예상체결 등락률 FHKST11860000
//	EP9   InquireMarketValue            — 시장가치 순위         FHPST01790000
//	EP10  InquireDisparity              — 이격도 순위           FHPST01780000
//	EP11  InquirePreferDisparateRatio   — 우선주 괴리율         FHPST01770000
//	EP12  InquireProfitAssetIndex       — 수익자산지표 순위     FHPST01730000
//	EP13  InquireMktfunds               — 증시자금 종합         FHKST649100C0
//
// Anomalies (Phase 4.3):
//
//	EP2/EP3 이중 output (output1 단일+output2 배열)
//	EP5 비표준 UPPERCASE 파라미터 (FID_ 접두어 없음: EXCG_DVSN_CD/PDNO 등), CTSC prefix TR_ID
//	EP6/EP7/EP9/EP10/EP11/EP12 lowercase fid_* query params
//	EP8 output 단일 객체 (배열 아님), UPPERCASE FID_ params
//	EP13 KOFIA 데이터: 고객예탁금/신용융자/미수금/펀드/담보대출 (단위: 억원/백만원)
//
// 사용자는 root kis.Client 의 Domestic 필드로 접근.
package domestic
