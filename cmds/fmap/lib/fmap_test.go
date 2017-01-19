// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmap

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// Flash map is stored in little-endian.
var fmapName = []byte("Fake flash" + strings.Repeat("\x00", 32-10))
var area0Name = []byte("Area Number 1\x00\x00\x00Hello" + strings.Repeat("\x00", 32-21))
var area1Name = []byte("Area Number 2xxxxxxxxxxxxxxxxxxx")
var fakeFlash = bytes.Join([][]byte{
	// Arbitrary data
	bytes.Repeat([]byte{0x53, 0x11, 0x34, 0x22}, 94387),

	// Signature
	[]byte("__FMAP__"),
	// VerMajor, VerMinor
	{1, 0},
	// Base
	{0xef, 0xbe, 0xad, 0xde, 0xbe, 0xba, 0xfe, 0xca},
	// Size
	{0x11, 0x22, 0x33, 0x44},
	// Name (32 bytes)
	fmapName,
	// NAreas
	{0x02, 0x00},

	// Areas[0].Offset
	{0xef, 0xbe, 0xad, 0xde},
	// Areas[0].Size
	{0x11, 0x11, 0x11, 0x11},
	// Areas[0].Name (32 bytes)
	area0Name,
	// Areas[0].Flags
	{0x13, 0x10},

	// Areas[1].Offset
	{0xbe, 0xba, 0xfe, 0xca},
	// Areas[1].Size
	{0x22, 0x22, 0x22, 0x22},
	// Areas[1].Name (32 bytes)
	area1Name,
	// Areas[1].Flags
	{0x00, 0x00},
}, []byte{})

func TestReadFMap(t *testing.T) {
	r := bytes.NewReader(fakeFlash)
	fmap := ReadFMap(r)
	expected := FMap{
		fMapInternal: fMapInternal{
			VerMajor: 1,
			VerMinor: 0,
			Base:     0xcafebabedeadbeef,
			Size:     0x44332211,
		},
		Areas: []FMapArea{
			{
				Offset: 0xdeadbeef,
				Size:   0x11111111,
				Flags:  0x1013,
			}, {
				Offset: 0xcafebabe,
				Size:   0x22222222,
				Flags:  0x0000,
			},
		},
		Start: 4 * 94387,
	}
	copy(expected.Signature[:], []byte("__FMAP__"))
	copy(expected.Name[:], fmapName)
	copy(expected.Areas[0].Name[:], area0Name)
	copy(expected.Areas[1].Name[:], area1Name)
	if !reflect.DeepEqual(*fmap, expected) {
		t.Errorf("expected: %+v\ngot: %+v", expected, *fmap)
	}
}

func TestFieldNames(t *testing.T) {
	r := bytes.NewReader(fakeFlash)
	fmap := ReadFMap(r)
	for i, expected := range []string{"STATIC|COMPRESSED|0x1010", "0x0"} {
		got := FlagNames(fmap.Areas[i].Flags)
		if got != expected {
			t.Errorf("expected: %s\ngot: %s", expected, got)
		}
	}
}