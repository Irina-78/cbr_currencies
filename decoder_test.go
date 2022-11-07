package main

import (
	"testing"
)

func TestCbrDecoder(t *testing.T) {
	s := ``
	decoder := newCbrDecoder(s)
	result := CbrResult{}
	if err := decoder.Decode(&result); err == nil {
		t.Fatalf("expected an error got nil")
	}

	s = `
	<?xml version="1.0" encoding="windows-1251"?>
	<ValCurs Date="20.01.2007" name="Foreign Currency Market">
		<Valute ID="R01035">
			<NumCode>826</NumCode>
			<CharCode>GBP</CharCode>
			<Nominal>1</Nominal>
			<Name>British Pound Sterling</Name>
			<Value>52,3656</Value>
		</Valute>`
	decoder = newCbrDecoder(s)
	result = CbrResult{}
	if err := decoder.Decode(&result); err == nil {
		t.Fatalf("expected an error got nil")
	}

	s = `
	<?xml version="1.0" encoding="windows-1251"?>
	<ValCurs Date="20.01.2007" name="Foreign Currency Market">
		<Valute ID="R01035">
			<NumCode>826</NumCode>
			<CharCode>GBP</CharCode>
			<Nominal>1</Nominal>
			<Name>British Pound Sterling</Name>
			<Value>52,3656</Value>
		</Valute>
		<Valute ID="R01090">
			<NumCode>974</NumCode>
			<CharCode>BYR</CharCode>
			<Nominal>1000</Nominal>
			<Name>Belarussian Ruble</Name>
			<Value>12,3701</Value>
		</Valute>
		<Valute ID="R01235">
			<NumCode>840</NumCode>
			<CharCode>USD</CharCode>
			<Nominal>1</Nominal>
			<Name>US Dollar</Name>
			<Value>26,5075</Value>
		</Valute>
		<Valute ID="R01239">
			<NumCode>978</NumCode>
			<CharCode>EUR</CharCode>
			<Nominal>1</Nominal>
			<Name>Euro</Name>
			<Value>34,4173</Value>
		</Valute>
		<Valute ID="R01820">
			<NumCode>392</NumCode>
			<CharCode>JPY</CharCode>
			<Nominal>100</Nominal>
			<Name>Japanese Yen</Name>
			<Value>21,8528</Value>
		</Valute>
	</ValCurs>`
	decoder = newCbrDecoder(s)
	result = CbrResult{}
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("got an error: %v", err)
	}
}
