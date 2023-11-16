package util

import (
	"fmt"
	"testing"
)

func TestGetServiceName(t *testing.T) {
	fmt.Println(GetServiceName("xxx-ss-dd"))
}

func TestGetCpuPercent(t *testing.T) {
	fmt.Println(GetCpuPercent())
}

func TestGetMemPercent(t *testing.T) {
	fmt.Println(GetMemPercent())
}