package model

import "os"

type Post struct {
	ID   int     `db:"id"`
	D0   int     `db:"d0"`
	D1   int     `db:"d1"`
	D2   int     `db:"d2"`
	D3   int     `db:"d3"`
	D4   int     `db:"d4"`
	D5   int     `db:"d5"`
	D6   int     `db:"d6"`
	D7   int     `db:"d7"`
	D8   int     `db:"d8"`
	D9   int     `db:"d9"`
	D10  int     `db:"d10"`
	D11  int     `db:"d11"`
	D12  int     `db:"d12"`
	D13  int     `db:"d13"`
	D14  int     `db:"d14"`
	D15  int     `db:"d15"`
	D16  int     `db:"d16"`
	D17  int     `db:"d17"`
	D18  int     `db:"d18"`
	D19  int     `db:"d19"`
	D20  int     `db:"d20"`
	D21  int     `db:"d21"`
	D22  int     `db:"d22"`
	D23  int     `db:"d23"`
	D24  int     `db:"d24"`
	D608 int     `db:"d608"`
	D609 int     `db:"d609"`
	D610 int     `db:"d610"`
	D611 int     `db:"d611"`
	D612 int     `db:"d612"`
	D613 int     `db:"d613"`
	D614 int     `db:"d614"`
	D618 int     `db:"d618"`
	D619 int     `db:"d619"`
	D620 int     `db:"d620"`
	D621 int     `db:"d621"`
	D622 int     `db:"d622"`
	D623 int     `db:"d623"`
	D624 int     `db:"d624"`
	D625 int     `db:"d625"`
	D626 int     `db:"d626"`
	D627 int     `db:"d627"`
	D628 int     `db:"d628"`
	D629 int     `db:"d629"`
	D630 int     `db:"d630"`
	D631 int     `db:"d631"`
	D632 int     `db:"d632"`
	D633 int     `db:"d633"`
	D634 int     `db:"d634"`
	D635 int     `db:"d635"`
	D800 int     `db:"d800"`
	D802 int     `db:"d802"`
	D804 int     `db:"d804"`
	D806 int     `db:"d806"`
	D808 int     `db:"d808"`
	D810 int     `db:"d810"`
	D812 int     `db:"d812"`
	D814 int     `db:"d814"`
	D816 int     `db:"d816"`
	D818 int     `db:"d818"`
	D820 int     `db:"d820"`
	D106 float32 `db:"d106"`
	D136 float32 `db:"d136"`
	D138 float32 `db:"d138"`
	D140 float32 `db:"d140"`
	D148 float32 `db:"d148"`
	D150 float32 `db:"d150"`
	D166 float32 `db:"d166"`
	D190 float32 `db:"d190"`
	D192 float32 `db:"d192"`
	D364 float32 `db:"d364"`
	D366 float32 `db:"d366"`
	D392 float32 `db:"d392"`
	D534 float32 `db:"d534"`
	D536 float32 `db:"d536"`
	D538 float32 `db:"d538"`
	D540 float32 `db:"d540"`
	D542 float32 `db:"d542"`
	D544 float32 `db:"d544"`
	D546 float32 `db:"d546"`
	D774 float32 `db:"d774"`
	D776 float32 `db:"d776"`
	D778 float32 `db:"d778"`
	D676 float32 `db:"d676"`
	D650 float32 `db:"d650"`
}

type Message struct {
	Address string      `json:"address"`
	Value   interface{} `json:"value"`
}

func (Post) TableName() string {
	tableName := os.Getenv("TABLE_NAME")
	return tableName
}
