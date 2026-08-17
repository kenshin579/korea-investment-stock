package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInquireExpClosingPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/exp-closing-price`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "exp_closing_price_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireExpClosingPrice(context.Background(), domestic.InquireExpClosingPriceParams{
		RankSortClsCode: "0",
		Symbol:          "0000",
		BlngClsCode:     "0",
	})
	require.NoError(t, err)
	require.Len(t, res.Output1, 2)

	// wire param keys 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11173", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_RANK_SORT_CLS_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "0", capturedQuery.Get("FID_BLNG_CLS_CODE"))

	item := res.Output1[0]
	assert.Equal(t, "005930", item.StckShrnIscd)
	assert.Equal(t, "삼성전자", item.HtsKorIsnm)
	assert.Equal(t, "82500", item.StckPrpr.String())
	assert.Equal(t, "500", item.PrdyVrss.String())
	assert.Equal(t, "2", item.PrdyVrssSign)
	assert.InDelta(t, 0.61, item.PrdyCtrt, 0.001)
	assert.Equal(t, "1000", item.SdprVrssPrpr.String())
	assert.InDelta(t, 1.23, item.SdprVrssPrprRate, 0.001)
	assert.Equal(t, int64(125000), item.CntgVol)
}

func TestInquireChkHoliday(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/chk-holiday`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "chk_holiday_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireChkHoliday(context.Background(), domestic.InquireChkHolidayParams{
		BassDt:    "20260817",
		CtxAreaNk: "",
		CtxAreaFk: "",
	})
	require.NoError(t, err)

	// non-FID UPPERCASE wire param keys 검증
	assert.Equal(t, "20260817", capturedQuery.Get("BASS_DT"))
	assert.Empty(t, capturedQuery.Get("FID_INPUT_DATE_1"), "BASS_DT 파라미터만 사용해야 함 (FID_ 아님)")

	// output 은 배열이다. BASS_DT 하루가 아니라 그 날짜부터의 달력을 페이지로 돌려준다.
	// 예전 픽스처는 output 을 단일 객체로 적어놨고 타입도 그에 맞춰져 있어서,
	// 테스트는 통과하는데 실제 호출은 "cannot unmarshal array into Go struct field" 로
	// 깨졌다. 실측 응답 형태로 고정한다.
	require.Len(t, res.Output, 4)

	// 첫 항목이 BASS_DT 인 것은 관측된 동작일 뿐이라 순서에 의존하지 않고 날짜로 찾는다.
	find := func(bassDt string) *domestic.ChkHolidayItem {
		for i := range res.Output {
			if res.Output[i].Bassdt == bassDt {
				return &res.Output[i]
			}
		}
		return nil
	}

	// 2026-08-17(월)은 공휴일 — 영업일도 개장일도 아니다.
	closed := find("20260817")
	require.NotNil(t, closed)
	assert.Equal(t, "02", closed.WdayDvsnCd)
	assert.Equal(t, "N", closed.BzdyYn)
	assert.Equal(t, "N", closed.OpndYn)
	assert.Equal(t, "N", closed.SttlDayYn)
	assert.Equal(t, "Y", closed.TrDayYn, "거래일여부는 개장일여부와 다르다 — 증권 업무 가능일 기준")

	// 2026-08-18(화)은 정상 개장일.
	open := find("20260818")
	require.NotNil(t, open)
	assert.Equal(t, "03", open.WdayDvsnCd)
	assert.Equal(t, "Y", open.BzdyYn)
	assert.Equal(t, "Y", open.OpndYn)

	// 주말(토)도 개장일이 아니다.
	sat := find("20260822")
	require.NotNil(t, sat)
	assert.Equal(t, "N", sat.OpndYn)

	// 페이지가 더 남았다는 신호. 우리가 필요한 날짜가 1페이지에 있으면 연속조회는 불필요하다.
	assert.Contains(t, res.Msg1, "조회가 계속됩니다")
}

func TestInquireViStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-vi-status`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "vi_status_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireViStatus(context.Background(), domestic.InquireViStatusParams{
		DivClsCode:      "0",
		MrktClsCode:     "0",
		Symbol:          "",
		RankSortClsCode: "0",
		InputDate1:      "20260507",
		TrgtClsCode:     "",
		TrgtExlsCode:    "",
	})
	require.NoError(t, err)
	require.NotNil(t, res.Output)

	// wire param keys 검증
	assert.Equal(t, "20139", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_MRKT_CLS_CODE"))
	assert.Equal(t, "20260507", capturedQuery.Get("FID_INPUT_DATE_1"))

	out := res.Output
	assert.Equal(t, "삼성전자", out.HtsKorIsnm)
	assert.Equal(t, "005930", out.MkscShrnIscd)
	assert.Equal(t, "Y", out.ViClsCode)
	assert.Equal(t, "20260507", out.BsopDate)
	assert.Equal(t, "100530", out.CntgViHour)
	assert.Equal(t, "100830", out.ViCnclHour)
	assert.Equal(t, "1", out.ViKindCode)
	assert.Equal(t, "82500", out.ViPrc.String())
	assert.Equal(t, "81000", out.ViStndPrc.String())
	assert.InDelta(t, 1.85, out.ViDprt, 0.001)
	assert.Equal(t, "80500", out.ViDmcStndPrc.String())
	assert.InDelta(t, 2.48, out.ViDmcDprt, 0.001)
	assert.Equal(t, int64(3), out.ViCount)
}

func TestInquireCaptureUplowprice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/capture-uplowprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "capture_uplowprice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireCaptureUplowprice(context.Background(), domestic.InquireCaptureUplowpriceParams{
		PrcClsCode: "0",
		DivClsCode: "0",
		Symbol:     "0000",
	})
	require.NoError(t, err)
	require.Len(t, res.Output, 2)

	// wire param keys 검증
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "11300", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_PRC_CLS_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_DIV_CLS_CODE"))
	assert.Equal(t, "0000", capturedQuery.Get("FID_INPUT_ISCD"))

	item := res.Output[0]
	assert.Equal(t, "005930", item.MkscShrnIscd)
	assert.Equal(t, "삼성전자", item.HtsKorIsnm)
	assert.Equal(t, "82500", item.StckPrpr.String())
	assert.Equal(t, "2", item.PrdyVrssSign)
	assert.Equal(t, "500", item.PrdyVrss.String())
	assert.InDelta(t, 0.61, item.PrdyCtrt, 0.001)
	assert.Equal(t, int64(8750000), item.AcmlVol)
	assert.Equal(t, int64(250000), item.TotalAskpRsqn)
	assert.Equal(t, int64(310000), item.TotalBidpRsqn)
	assert.Equal(t, int64(45000), item.AskpRsqn1)
	assert.Equal(t, int64(62000), item.BidpRsqn1)
	assert.Equal(t, int64(9200000), item.PrdyVol)
	assert.Equal(t, int64(4200000), item.SelnCnqn)
	assert.Equal(t, int64(4550000), item.ShnuCnqn)
	assert.Equal(t, "57750", item.StckLlam.String())
	assert.Equal(t, "107250", item.StckMxpr.String())
	assert.InDelta(t, -4.89, item.PrdyVrssVolRate, 0.001)
}
