package scaffold

// Entry defines an entry in the ISO 4217 standard.
type Entry struct {
	Code        string `xml:"Ccy,omitempty" json:"AlphanumericCode,omitempty"`
	MinorUnits  string `xml:"CcyMnrUnts,omitempty" json:"MinorUnits,omitempty"`
	Country     string `xml:"CtryNm,omitempty" json:"CountryName,omitempty"`
	Description string `xml:"CcyNm,omitempty" json:"CurrencyName,omitempty"`
}

// Table defines the table containing the ISO entries.
type Table struct {
	Entries []*Entry `xml:"CcyNtry,omitempty" json:"Entries,omitempty"`
}

// ISO4217 defines the top level xml response for the ISO 4217 standard.
type ISO4217 struct {
	AttrPublished string `xml:"Pblshd,attr" json:",omitempty"` // maxLength=10
	Table         *Table `xml:"CcyTbl,omitempty" json:"Table,omitempty"`
}
