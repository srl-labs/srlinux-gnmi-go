package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/karimra/gnmic/target"
	"github.com/karimra/gnmic/types"
	"github.com/karimra/gnmic/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func main() {
	// gnmi target config
	tc := &types.TargetConfig{
		Address:  "clab-srl-srl:57400",
		Username: stringP("admin"),
		Password: stringP("admin"),

		Insecure:   boolP(false), // srlinux uses tls secured gnmi connection
		SkipVerify: boolP(true),  // don't check certification validity
		// Gzip:       boolP(false),
		Timeout: 30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a gnmi target with the specified config
	gnmiTarget := target.NewTarget(tc)

	if err := gnmiTarget.CreateGNMIClient(ctx); err != nil {
		log.Fatalf("failed to create gnmi client %v", err)
	}

	mgmtIPPath, _ := utils.ParsePath("/interface[name=mgmt0]/subinterface[index=0]/ipv4/address/status")

	rsp, err := gnmiTarget.Get(ctx,
		&gnmi.GetRequest{
			Path:     []*gnmi.Path{mgmtIPPath},
			Type:     gnmi.GetRequest_STATE,
			Encoding: gnmi.Encoding_JSON_IETF,
		})
	if err != nil {
		log.Fatalf("failed to execute get request %v", err)
	}

	log.Printf("entire response: %+v\n\n", rsp)

	for _, n := range rsp.GetNotification() {
		for _, u := range n.GetUpdate() {
			b := u.GetVal().GetJsonIetfVal()
			var j map[string][]map[string]string
			json.Unmarshal(b, &j)

			fmt.Println("management address:", j["address"][0]["ip-prefix"])
		}
	}

}

// boolP returns a pointer to a bool.
func boolP(b bool) *bool {
	return &b
}

// stringP returns a pointer to a string.
func stringP(s string) *string {
	return &s
}
