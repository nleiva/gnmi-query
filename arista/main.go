package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	arista "github.com/nleiva/yang-data-structures/gnmi/arista/aristapath"
	eos "github.com/nleiva/yang-data-structures/gnmi/arista"
	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func main() {
	targetAddr := "10.0.0.1:6030"
	user := "admin"
	pass := "admin"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := initClient(ctx, targetAddr, user, pass)
	if err != nil {
		log.Fatalf("failed to init client: %v", err)
	}

	//////////////////////////
	// Get one value for PATH
	//////////////////////////
	pathOne := arista.Root().Interface("Ethernet3").Subinterface(0).Ipv4()

	val, err := ygnmi.Get(ctx, c, pathOne.Config())
	if err != nil {
		log.Fatalf("failed to get one: %v", err)
	}

	p, _, err := ygnmi.ResolvePath(pathOne.Config().PathStruct())
	if err != nil {
		log.Fatalf("failed to resolve path: %v", err)
	}

	st, err := ygot.PathToString(p)
	if err != nil {
		log.Fatalf("failed to convert path to string: %v", err)
	}

	fmt.Printf("Path: %v\n", st)
	for _, v := range val.Address {
		fmt.Printf("Address: %v/%v\n", *v.Ip, *v.PrefixLength)
	}
	fmt.Println()

	////////////////////////////////////
	// Get all values for wildcard PATH
	///////////////////////////////////
	pathAll := arista.Root().InterfaceAny().Subinterface(0).Ipv4().AddressMap()

	all, err := ygnmi.LookupAll(context.Background(), c, pathAll.Config())
	if err != nil {
		log.Fatalf("failed to get all paths: %v", err)
	}
	for _, single := range all {
		fmt.Printf("Interface: %v", single.Path.GetElem()[1].Key["name"])
		val, ok := single.Val()
		if !ok {
			continue
		}
		for _, v := range val {
			fmt.Printf("  -> Address: %v/%v\n", *v.Ip, *v.PrefixLength)
		}
	}

	fmt.Println()

	////////////////////////////////////
	// Reconcile (WIP)
	///////////////////////////////////

	serverAddress := "100.64.1.1"

	ntpPath := arista.Root().System().Ntp()
	p, _, err = ygnmi.ResolvePath(ntpPath.Config().PathStruct())
	if err != nil {
		log.Fatalf("failed to resolve path: %v", err)
	}
	fmt.Printf("Path: %v\n", p)

	// Create a reconciler for System_Ntp
	r, err := ygnmi.NewReconciler(c, ntpPath.Config())
	if err != nil {
		log.Fatalf("failed to create reconciler: %v", err)
	}

	// Desired state
	root := new(eos.Eos)
	ntp := root.GetOrCreateSystem().GetOrCreateNtp()
	// ntp.Enabled =  ygot.Bool(true)
	ntp.GetOrCreateServer(serverAddress)


	r.Start(ctx, func(cfg, state *ygnmi.Value[*eos.System_Ntp]) error {
		val, err := ygnmi.Lookup(ctx, c, ntpPath.Config())
		if err != nil {
			return err
		}

		cfgV, ok := cfg.Val()
		if !ok {
			fmt.Printf(">>>>> path %s: %s\n", cfg.Path.String(), ygnmi.ErrNotPresent.Error())
		}
		if d := cmp.Diff(cfgV, ntp); d != "" {
			fmt.Printf(">>>>> unexpected cfg diff detected:\n %s\n", d)

			// Enforce desired state
			res, err := ygnmi.Replace(ctx, c, ntpPath.Config(), ntp)
			if err != nil {
				log.Fatalf("gNMI set replace failed: %v", err)
			}

			fmt.Printf("config enforced at: %v for %v\n\n", res.Timestamp.Format("2006-01-02 15:04:05"), val.Path.String())
		}

		return ygnmi.Continue
	},
	)

	wantErr := "context deadline exceeded"
	err = r.Await()
	if diff := errdiff.Substring(err, wantErr); !errors.Is(err, io.EOF) && diff != "" {
		fmt.Printf("watch() returned unexpected diff: %s", diff)
		return
	}
}
