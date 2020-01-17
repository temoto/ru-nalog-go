package ru_nalog

//go:generate ./script/generate

import (
	"fmt"
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
	case DataKindBytes:
		if x, ok := value.([]byte); ok {
			self.value = x
		}
		if x, ok := value.(string); ok {
			self.value = []byte(x)
		}
	case DataKindString:
		self.value = value.(string)
	case DataKindVLN:
		self.value = toVLN(value, self.TagDesc.Length)
	// case DataKindTime:
	default: // FIXME
		self.value = value
	}
}

func (self *TLV) Bytes() []byte {
	return self.value.([]byte)
}

func (self *TLV) String() string {
	if self.Kind != DataKindString {
		panic("kind")
	}
	if !self.Varlen {
		return self.FixedString()
	}
	x := toString(self.value)
	if s, ok := x.(string); ok {
		return s
	}
	panic("value")
}

func (self *TLV) FixedString() string {
	crude := toString(self.value)
	if s, ok := crude.(string); ok {
		return strings.TrimRightFunc(s, isSpace)
	}
	panic("value")
}

func (self *TLV) Bool() bool {
	return self.value.(bool)
}

func (self *TLV) Float64() float64 {
	return self.value.(float64)
}

func (self *TLV) Time() time.Time {
	return self.value.(time.Time)
}

func (self *TLV) Uint32() uint32 {
	return self.value.(uint32)
}

// VLN / byte[] bug
func (self *TLV) Uint64() uint64 {
	if u64, ok := self.value.(uint64); ok {
		return u64
	}
	if u32, ok := self.value.(uint32); ok {
		return uint64(u32)
	}
	if ui, ok := self.value.(uint); ok {
		return uint64(ui)
	}
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

func toString(v interface{}) interface{} {
	if s, ok := v.(string); ok {
		return s
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Errorf("toString v=%q", v)
}

func toVLN(v interface{}, length uint16) interface{} {
	var u uint64
	switch n := v.(type) {
	case int32:
		u = uint64(n)
	case uint32:
		u = uint64(n)
	case int:
		u = uint64(n)
	case uint:
		u = uint64(n)
	case int64:
		u = uint64(n)
	case uint64:
		u = uint64(n)
	// TODO case string: strconv.ParseInt
	default:
		return fmt.Errorf("toVLN v=%q", v)
	}
	if length <= 6 {
		return uint32(u)
	}
	return u
}
