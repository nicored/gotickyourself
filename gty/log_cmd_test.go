package main

import (
	"testing"

	"time"

	"fmt"

	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

func TestGetHours(t *testing.T) {
	testCases := []struct {
		Desc   string
		Arg    string
		Hours  float32
		IsTime bool
	}{
		{Desc: "Test Hours 2.0", Arg: "2.0", Hours: 2.0, IsTime: true},
		{Desc: "Test Hours 2", Arg: "2", Hours: 2.0, IsTime: true},
		{Desc: "Test Hours 0.0", Arg: "0.0", Hours: 0.0, IsTime: true},
		{Desc: "Test Hours 12.5", Arg: "12.5", Hours: 12.5, IsTime: true},
		{Desc: "Test Hours 12.567", Arg: "12.567", Hours: 12.567, IsTime: true},
		{Desc: "Test Hours 12.", Arg: "12.", Hours: 12.0, IsTime: true},
		{Desc: "Test string with time 2.0", Arg: "some time 2.0", Hours: 0.0, IsTime: false},
	}

	for _, tc := range testCases {
		hours, isTime := getHours(tc.Arg)
		assert.Equal(t, tc.IsTime, isTime, tc.Desc)
		assert.Equal(t, tc.Hours, hours, tc.Desc)
	}
}

func TestGetDate(t *testing.T) {
	fmt.Println(parseDate("today"))
}

func TestGetTimePeriodStart(t *testing.T) {
	testCases := []struct {
		Desc   string
		Arg    string
		Time   string
		IsTime bool
	}{
		{Desc: "Test today", Arg: "today", Time: now.BeginningOfDay().String(), IsTime: true},
		{Desc: "Test 2 days", Arg: "2 days", Time: now.BeginningOfDay().Add(-48 * time.Hour).String(), IsTime: true},
		{Desc: "Test 2 days with spaces", Arg: "2         days            ", Time: now.BeginningOfDay().Add(-48 * time.Hour).String(), IsTime: true},
		{Desc: "Test 3 day in singular", Arg: "3 day", Time: now.BeginningOfDay().Add(-24 * 3 * time.Hour).String(), IsTime: true},
		{Desc: "Test 1 week", Arg: "1 week", Time: now.Monday().Add(-24 * 7 * time.Hour).String(), IsTime: true},
		{Desc: "Test 1", Arg: "1", Time: "", IsTime: false},
		{Desc: "Test yesterday", Arg: "yesterday", Time: now.BeginningOfDay().Add(-24 * time.Hour).String(), IsTime: true},
	}

	for _, tc := range testCases {
		actualTime, isTime := getTimePeriodStart(tc.Arg)
		assert.Equal(t, tc.IsTime, isTime, tc.Desc)

		if tc.Time == "" {
			assert.Nil(t, actualTime, tc.Desc)
		} else {
			assert.NotNil(t, actualTime, tc.Desc)

			if actualTime != nil {
				assert.Equal(t, tc.Time, actualTime.String(), tc.Desc)
			}
		}
	}
}
