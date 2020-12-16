// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestGobEncoding(t *testing.T) {
	g := NewGomegaWithT(t)

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []string{
		"P0D",
		"P1D",
		"P1W",
		"P1M",
		"P1Y",
		"PT1H",
		"PT1M",
		"PT1S",
		"P2Y3M4W5D",
		"-P2Y3M4W5D",
		"P2Y3M4W5DT1H7M9S",
		"-P2Y3M4W5DT1H7M9S",
		"P48M",
	}
	for i, c := range cases {
		period := MustParse(c, false)
		var p Period
		err := encoder.Encode(&period)
		g.Expect(err).NotTo(HaveOccurred(), info(i, c))
		if err == nil {
			err = decoder.Decode(&p)
			g.Expect(err).NotTo(HaveOccurred(), info(i, c))
			g.Expect(p).To(Equal(period), info(i, c))
		}
	}
}

func TestPeriodJSONMarshalling(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -3, -11, -59, -59), `"-P1111Y4M3DT11H59M59S"`},
		{New(-1, -10, -31, -5, -4, -20), `"-P1Y10M31DT5H4M20S"`},
		{New(0, 0, 0, 0, 0, 0), `"P0D"`},
		{New(0, 0, 0, 0, 0, 1), `"PT1S"`},
		{New(0, 0, 0, 0, 1, 0), `"PT1M"`},
		{New(0, 0, 0, 1, 0, 0), `"PT1H"`},
		{New(0, 0, 1, 0, 0, 0), `"P1D"`},
		{New(0, 1, 0, 0, 0, 0), `"P1M"`},
		{New(1, 0, 0, 0, 0, 0), `"P1Y"`},
	}
	for i, c := range cases {
		var p Period
		bb, err := json.Marshal(c.value)
		g.Expect(err).NotTo(HaveOccurred(), info(i, c))
		g.Expect(string(bb)).To(Equal(c.want), info(i, c))
		if string(bb) == c.want {
			err = json.Unmarshal(bb, &p)
			g.Expect(err).NotTo(HaveOccurred(), info(i, c))
			g.Expect(p).To(Equal(c.value), info(i, c))
		}
	}
}

func TestPeriodTextMarshalling(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -3, -11, -59, -59), "-P1111Y4M3DT11H59M59S"},
		{New(-1, -9, -31, -5, -4, -20), "-P1Y9M31DT5H4M20S"},
		{New(0, 0, 0, 0, 0, 0), "P0D"},
		{New(0, 0, 0, 0, 0, 1), "PT1S"},
		{New(0, 0, 0, 0, 1, 0), "PT1M"},
		{New(0, 0, 0, 1, 0, 0), "PT1H"},
		{New(0, 0, 1, 0, 0, 0), "P1D"},
		{New(0, 1, 0, 0, 0, 0), "P1M"},
		{New(1, 0, 0, 0, 0, 0), "P1Y"},
	}
	for i, c := range cases {
		var p Period
		bb, err := c.value.MarshalText()
		g.Expect(err).NotTo(HaveOccurred(), info(i, c))
		g.Expect(string(bb)).To(Equal(c.want), info(i, c))
		if string(bb) == c.want {
			err = p.UnmarshalText(bb)
			g.Expect(err).NotTo(HaveOccurred(), info(i, c))
			g.Expect(p).To(Equal(c.value), info(i, c))
		}
	}
}

func TestInvalidPeriodText(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value string
		want  string
	}{
		{``, `cannot parse a blank string as a period`},
		{`not-a-period`, `not-a-period: expected 'P' period mark at the start`},
		{`P000`, `P000: missing designator at the end`},
	}
	for i, c := range cases {
		var p Period
		err := p.UnmarshalText([]byte(c.value))
		g.Expect(err).To(HaveOccurred(), info(i, c))
		g.Expect(err.Error()).To(Equal(c.want), info(i, c))
	}
}
