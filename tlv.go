package ru_nalog

//go:generate ./script/generate

import (
	"sort"
	"strings"
	"time"
)

type Tag uint16

type TagDesc struct {
	Kind   DataKind
	Tag    Tag
	Length uint16
	Varlen bool
}

var userTags []TagDesc

// Will stable sort argument, return previous value.
// Not thread-safe.
func RegisterTags(ts []TagDesc) (prev []TagDesc) {
	prev, userTags = userTags, ts
	sort.SliceStable(userTags, func(i, j int) bool { return userTags[i].Tag < userTags[j].Tag })
	return prev
}

func searchSortedTagDesc(tag Tag, xs []TagDesc) *TagDesc {
	// please forgive optimisation junkie
	// it's 2x faster than sort.Search() because f() is inlined and early return
	// copy from go stdlib sort.Search
	i, j := uint32(0), uint32(len(xs))
	for i < j {
		h := (i + j) >> 1
		item := &xs[h]
		if item.Tag == tag {
			return item
		} else if item.Tag < tag {
			i = h + 1
		} else {
			j = h
		}
	}
	return nil
}

// Searches in user tags (RegisterTags) first, then in builtin tags.
// Not thread-safe.
func FindTag(tag Tag) *TagDesc {
	// very likely shortcut
	if userTags != nil {
		if u := searchSortedTagDesc(tag, userTags); u != nil {
			return u
		}
	}
	return searchSortedTagDesc(tag, builtinTags[:])
}

type TLV struct {
	TagDesc
	Caption   string
	Printable string
	value     interface{}
}

func NewTLV(tag Tag) *TLV {
	desc := FindTag(tag)
	if desc == nil {
		return nil
	}
	return &TLV{TagDesc: *desc}
}

func (self *TLV) Value() interface{} { return self.value }

func (self *TLV) SetValue(value interface{}) {
	switch self.TagDesc.Kind {
	case DataKindBool:
		self.value = value.(bool)
	case DataKindString:
		self.value = value.(string)
	// case DataKindTime:
	default: // FIXME
		self.value = value
	}
}

func (self *TLV) String() string {
	if self.Kind != DataKindString {
		panic("kind")
	}
	if !self.Varlen {
		return self.FixedString()
	}
	return self.value.(string)
}

func (self *TLV) FixedString() string {
	crude := self.value.(string)
	return strings.TrimRightFunc(crude, isSpace)
}

func (self *TLV) Bool() bool {
	return self.value.(bool)
}

func (self *TLV) Time() time.Time {
	return self.value.(time.Time)
}

func (self *TLV) Uint32() uint32 {
	return self.value.(uint32)
}

// VLN / byte[] bug
func (self *TLV) Uint64() uint64 {
	return self.value.(uint64)
}

type FindByTager interface {
	FindByTag(Tag) *TLV
}

func (self *TLV) FindByTag(tag Tag) *TLV {
	if self.Tag == tag {
		return self
	}
	if self.Kind == DataKindSTLV {
		for _, child := range self.value.([]TLV) {
			if t := child.FindByTag(tag); t != nil {
				return t
			}
		}
	}
	return nil
}

func isSpace(r rune) bool { return r == ' ' }