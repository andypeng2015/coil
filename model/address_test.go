package model

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/cybozu-go/coil"
)

const (
	containerID = "5451faf2-da4f-4690-b024-b77518982f61"
)

func testGetAddressInfo(t *testing.T) {
	t.Parallel()
	m := NewTestEtcdModel(t)

	_, subnet, err := net.ParseCIDR("10.11.0.0/28")
	if err != nil {
		t.Fatal(err)
	}

	err = m.AddPool(context.Background(), "default", subnet, 2)
	if err != nil {
		t.Fatal(err)
	}

	block, err := m.AcquireBlock(context.Background(), "node1", "default")
	if err != nil {
		t.Fatal(err)
	}

	assignment := coil.IPAssignment{
		ContainerID: containerID,
		Namespace:   "default",
		Pod:         "test",
		CreatedAt:   time.Now().UTC(),
	}
	ip, err := m.AllocateIP(context.Background(), block, assignment)
	if err != nil {
		t.Fatal(err)
	}

	info, _, err := m.GetAddressInfo(context.Background(), ip)
	if err != nil {
		t.Fatal(err)
	}

	if *info != assignment {
		t.Errorf("expected info: %v, actual: %v", assignment, *info)
	}

	_, _, err = m.GetAddressInfo(context.Background(), net.ParseIP("10.11.1.0"))
	if err != ErrNotFound {
		t.Errorf("expected error: ErrNotFound, actual: %v", err)
	}
}

func TestAddress(t *testing.T) {
	t.Run("GetAddressInfo", testGetAddressInfo)
}
