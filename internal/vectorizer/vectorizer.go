package vectorizer

import (
	"time"
)

var MCCRisk = map[string]float32{
	"5411": 0.15, "5812": 0.30, "5912": 0.20, "5944": 0.45,
	"7801": 0.80, "7802": 0.75, "7995": 0.85, "4511": 0.35,
	"5311": 0.25, "5999": 0.50,
}

type Request struct {
	Transaction Transaction  `json:"transaction"`
	Customer    Customer     `json:"customer"`
	Merchant    Merchant     `json:"merchant"`
	Terminal    Terminal     `json:"terminal"`
	LastTx      *LastTx      `json:"last_transaction"`
}

type Transaction struct {
	Amount       float32 `json:"amount"`
	Installments int     `json:"installments"`
	RequestedAt  string  `json:"requested_at"`
}
type Customer struct {
	AvgAmount      float32  `json:"avg_amount"`
	TxCount24h     int      `json:"tx_count_24h"`
	KnownMerchants []string `json:"known_merchants"`
}
type Merchant struct {
	ID        string  `json:"id"`
	MCC       string  `json:"mcc"`
	AvgAmount float32 `json:"avg_amount"`
}
type Terminal struct {
	IsOnline    bool    `json:"is_online"`
	CardPresent bool    `json:"card_present"`
	KmFromHome  float32 `json:"km_from_home"`
}
type LastTx struct {
	Timestamp     string  `json:"timestamp"`
	KmFromCurrent float32 `json:"km_from_current"`
}

func clamp(v float32) float32 {
	if v < 0 { return 0 }
	if v > 1 { return 1 }
	return v
}

func Vectorize(r *Request) ([14]float32, uint8) {
	var v [14]float32

	v[0] = clamp(r.Transaction.Amount / 10000.0)
	v[1] = clamp(float32(r.Transaction.Installments) / 12.0)
	if r.Customer.AvgAmount > 0 {
		v[2] = clamp((r.Transaction.Amount / r.Customer.AvgAmount) / 10.0)
	}

	t, _ := time.Parse("2006-01-02T15:04:05Z", r.Transaction.RequestedAt)
	if t.IsZero() {
		t, _ = time.Parse(time.RFC3339, r.Transaction.RequestedAt)
	}
	v[3] = float32(t.Hour()) / 23.0
	// 0=Monday ... 6=Sunday (Weekday: 0=Sun → ajuste)
	v[4] = float32((int(t.Weekday())+6)%7) / 6.0

	if r.LastTx != nil {
		tPrev, _ := time.Parse("2006-01-02T15:04:05Z", r.LastTx.Timestamp)
		if tPrev.IsZero() {
			tPrev, _ = time.Parse(time.RFC3339, r.LastTx.Timestamp)
		}
		mins := float32(t.Sub(tPrev).Minutes())
		v[5] = clamp(mins / 1440.0)
		v[6] = clamp(r.LastTx.KmFromCurrent / 1000.0)
	} else {
		v[5] = -1
		v[6] = -1
	}

	v[7] = clamp(r.Terminal.KmFromHome / 1000.0)
	v[8] = clamp(float32(r.Customer.TxCount24h) / 20.0)

	if r.Terminal.IsOnline    { v[9]  = 1 }
	if r.Terminal.CardPresent { v[10] = 1 }

	v[11] = 1
	for _, km := range r.Customer.KnownMerchants {
		if km == r.Merchant.ID { v[11] = 0; break }
	}

	if risk, ok := MCCRisk[r.Merchant.MCC]; ok {
		v[12] = risk
	} else {
		v[12] = 0.5
	}
	v[13] = clamp(r.Merchant.AvgAmount / 10000.0)

	var key uint8
	if r.Terminal.IsOnline    { key |= 0x4 }
	if r.Terminal.CardPresent { key |= 0x2 }
	if v[11] == 1             { key |= 0x1 }

	return v, key
}