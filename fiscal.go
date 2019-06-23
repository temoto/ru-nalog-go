package ru_nalog

// kinds of TLV.value
type DataKind uint

const (
	DataKindInvalid DataKind = iota
	DataKindSTLV
	DataKindBool
	DataKindUint
	DataKindVLN
	DataKindFVLN
	DataKindTime
	DataKindString
	DataKindBytes
)

type DocType uint16

const (
	FDRegistration         DocType = 1
	FDRegChange            DocType = 11
	FDCycleOpen            DocType = 2
	FDStateReport          DocType = 21
	FDCheck                DocType = 3
	FDCorrectionCheck      DocType = 31
	FDBSO                  DocType = 4
	FDCorrectionBSO        DocType = 41
	FDCycleClose           DocType = 5
	FDStorageClose         DocType = 6
	FDOperatorConfirmation DocType = 7
)

type Doc struct {
	Number uint32 `fdn:"1040"`
	Type   DocType
	Props  []TLV
}

func (d *Doc) FindByTag(tag Tag) *TLV {
	for _, child := range d.Props {
		if t := child.FindByTag(tag); t != nil {
			return t
		}
	}
	return nil
}
