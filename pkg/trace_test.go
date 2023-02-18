// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package gont_test

import (
	"flag"
	"testing"

	g "github.com/stv0g/gont/pkg"
	gopt "github.com/stv0g/gont/pkg/options"
	copt "github.com/stv0g/gont/pkg/options/capture"
	topt "github.com/stv0g/gont/pkg/options/trace"
	"github.com/stv0g/gont/pkg/trace"
)

var captureSocketAddr = flag.String("trace-socket", "tcp:[::1]:42125", "Listen address for capture socket")

func TestTracer(t *testing.T) {
	var (
		err error
		n   *g.Network
		h1  *g.Host
	)

	c1 := gopt.Capture(
		copt.Listener(*captureSocketAddr),
	)

	t1 := g.NewTracer(
		topt.ToFilename("trace.log"),
		topt.ToCapture(c1),
	)

	if err := t1.StartLocal(); err != nil {
		t.Fatalf("Failed to start: %s", err)
	}

	if n, err = g.NewNetwork(*nname,
		gopt.Customize(globalNetworkOptions,
			t1, c1,
		)...,
	); err != nil {
		t.Fatalf("Failed to create network: %s", err)
	}

	if h1, err = n.AddHost("h1",
		gopt.LogToDebug(true),
	); err != nil {
		t.Fatalf("Failed to add host: %s", err)
	}

	for i := 0; i < 5; i++ {
		_, _, err = h1.Run("../test/tracee/tracee", i)
		if err != nil {
			t.Fatalf("Failed to run tracee: %s", err)
		}

		trace.PrintfWithData(i, "Hello from test process")
	}
}
