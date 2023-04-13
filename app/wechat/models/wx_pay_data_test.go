package models

import (
	"testing"
)

func TestWxPayData_ToUrl(t *testing.T) {
	type fields struct {
		m_values map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"test1",
			fields{
				m_values: map[string]interface{}{
					"a": "a",
					"b": "b",
					"c": 1,
				},
			},
			"a=a&b=b&c=1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &WxPayData{
				m_values: tt.fields.m_values,
			}
			got, err := this.ToUrl()
			if (err != nil) != tt.wantErr {
				t.Errorf("WxPayData.ToUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WxPayData.ToUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWxPayData_ToJson(t *testing.T) {
	type fields struct {
		m_values map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"test1",
			fields{
				m_values: map[string]interface{}{
					"a": "a",
					"b": "b",
				},
			},
			`{"a":"a","b":"b"}`,
			false,
		}, {
			"test2",
			fields{
				m_values: map[string]interface{}{
					"a": 1,
					"b": "b",
				},
			},
			`{"a":1,"b":"b"}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &WxPayData{
				m_values: tt.fields.m_values,
			}
			got, err := this.ToJson()
			if (err != nil) != tt.wantErr {
				t.Errorf("WxPayData.ToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WxPayData.ToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
