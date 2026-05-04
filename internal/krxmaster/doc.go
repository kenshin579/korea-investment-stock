// Package krxmaster handles KRX (Korea Exchange) master file parsing.
//
// Master files from KRX are encoded in cp949 (Korean EUC-KR).
// Uses golang.org/x/text/encoding/korean for decoding.
package krxmaster

import (
	_ "golang.org/x/text/encoding/korean"
)
