// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
)

type AddonBindingStatus struct {
	Name  string              `json:"name"`
	Phase v1alpha1.AddonPhase `json:"phase"`
}

type Metric struct {
	Name   string        `json:"name"`
	Metric string        `json:"metric"`
	Format MetricFormat  `json:"format"`
	Values []MetricValue `json:"values"`
}

type MetricValue struct {
	Time  int64  `json:"time"`
	Value string `json:"value"`
}

type NewCluster struct {
	Name  string                `json:"name"`
	Track v1alpha1.ReleaseTrack `json:"track"`
	Ha    bool                  `json:"ha"`
}

type User struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

type MetricFormat string

const (
	MetricFormatBytes MetricFormat = "Bytes"
	MetricFormatCPU   MetricFormat = "CPU"
	MetricFormatTime  MetricFormat = "Time"
	MetricFormatPlain MetricFormat = "Plain"
)

var AllMetricFormat = []MetricFormat{
	MetricFormatBytes,
	MetricFormatCPU,
	MetricFormatTime,
	MetricFormatPlain,
}

func (e MetricFormat) IsValid() bool {
	switch e {
	case MetricFormatBytes, MetricFormatCPU, MetricFormatTime, MetricFormatPlain:
		return true
	}
	return false
}

func (e MetricFormat) String() string {
	return string(e)
}

func (e *MetricFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MetricFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MetricFormat", str)
	}
	return nil
}

func (e MetricFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
