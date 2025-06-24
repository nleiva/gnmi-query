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

	/////////////////////////
	// Get one value for PATH
	/////////////////////////
	// pathOne := arista.Root().System().Hostname()
	// pathOne := arista.Root().Interface("Ethernet3").Config()
	pathOne := arista.Root().Interface("Ethernet3").Subinterface(0).Ipv4().Config()

	val, err := ygnmi.Get(ctx, c, pathOne)
	if err != nil {
		log.Fatalf("failed to get one: %v", err)
	}

	p, _, err := ygnmi.ResolvePath(pathOne.PathStruct())
	if err != nil {
		log.Fatalf("failed to resolve path: %v", err)
	}
	fmt.Printf("Path: %v\n", p)
	for _, v := range val.Address {
		fmt.Printf("Address: %v\\%v\n", *v.Ip, *v.PrefixLength)
	}
	fmt.Println()

	///////////////////////////////////
	// Get all value for wildcard PATH
	//////////////////////////////////
	pathAll := arista.Root().InterfaceAny().Subinterface(0).Ipv4().AddressMap().Config()
	vals, err := ygnmi.GetAll(ctx, c, pathAll)
	if err != nil {
		log.Fatalf("failed to get all: %v", err)
	}
	p, _, err = ygnmi.ResolvePath(pathAll.PathStruct())
	if err != nil {
		log.Fatalf("failed to resolve path: %v", err)
	}
	fmt.Printf("Path: %v\n", p)

	// Get the value of each list element.
	// With GetAll it is impossible to know the path associated with a value,
	// so use LookupAll or Batch with with wildcard path instead.
	for idx, val := range vals {
		fmt.Printf("Value %d: ", idx)
		for _, v := range val {
			fmt.Printf("Address: %v\\%v\n", *v.Ip, *v.PrefixLength)
		}
	}

	fmt.Println()

	///////////////////////////////////
	// Reconcile (WIP)
	//////////////////////////////////

	// Define the query root (typed)
	Query := arista.Root().System().Ntp().Config()
	p, _, err = ygnmi.ResolvePath(Query.PathStruct())
	if err != nil {
		log.Fatalf("failed to resolve path: %v", err)
	}
	fmt.Printf("Path: %v\n", p)

	// Create a reconciler for System_Ntp
	r, err := ygnmi.NewReconciler(c, Query)
	if err != nil {
		log.Fatalf("failed to create reconciler: %v", err)
	}

	serverAddress := "100.64.1.1"

	desired := &eos.System_Ntp{
		Server: map[string]*eos.System_Ntp_Server{
			serverAddress: {
				Address: &serverAddress,
			},
		},
	}

	err = r.AddSubReconciler(arista.Root().System().Ntp().Server(serverAddress).Config(), func(cfg, state *ygnmi.Value[*eos.System_Ntp]) error {
		cfgV, _ := cfg.Val()
		if d := cmp.Diff(cfgV, desired); d != "" {
			fmt.Printf(">>>>> unexpected cfg diff detected:\n %s\n", d)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error adding subreconciler: %s", err)
	}

	r.Start(ctx, func(cfg, state *ygnmi.Value[*eos.System_Ntp]) error {
		return nil
	})

	// Check diff
	wantErr := "context deadline exceeded"
	err = r.Await()
	if diff := errdiff.Substring(err, wantErr); !errors.Is(err, io.EOF) && diff != "" {
		fmt.Printf("watch() returned unexpected diff: %s", diff)
		return
	}

	// Reconcile desired state (TODO)
}
