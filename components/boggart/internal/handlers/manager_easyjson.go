// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package handlers

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonEd74d837DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(in *jlexer.Lexer, out *managerHandlerDevice) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "tasks":
			if in.IsNull() {
				in.Skip()
				out.Tasks = nil
			} else {
				in.Delim('[')
				if out.Tasks == nil {
					if !in.IsDelim(']') {
						out.Tasks = make([][]string, 0, 2)
					} else {
						out.Tasks = [][]string{}
					}
				} else {
					out.Tasks = (out.Tasks)[:0]
				}
				for !in.IsDelim(']') {
					var v1 []string
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if v1 == nil {
							if !in.IsDelim(']') {
								v1 = make([]string, 0, 4)
							} else {
								v1 = []string{}
							}
						} else {
							v1 = (v1)[:0]
						}
						for !in.IsDelim(']') {
							var v2 string
							v2 = string(in.String())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Tasks = append(out.Tasks, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.Tags = append(out.Tags, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "id":
			out.ID = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "serial_number":
			out.SerialNumber = string(in.String())
		case "mac":
			out.MAC = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "probe_readiness":
			out.ProbeReadiness = string(in.String())
		case "probe_liveness":
			out.ProbeLiveness = string(in.String())
		case "metrics_descriptions_count":
			out.MetricsDescriptionsCount = uint64(in.Uint64())
		case "metrics_collect_count":
			out.MetricsCollectCount = uint64(in.Uint64())
		case "metrics_empty_count":
			out.MetricsEmptyCount = uint64(in.Uint64())
		case "mqtt_publishes":
			out.MQTTPublishes = int(in.Int())
		case "mqtt_subscribers":
			out.MQTTSubscribers = int(in.Int())
		case "logs_count":
			out.LogsCount = int(in.Int())
		case "has_metrics":
			out.HasMetrics = bool(in.Bool())
		case "has_widget":
			out.HasWidget = bool(in.Bool())
		case "logs_max_level":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.LogsMaxLevel).UnmarshalText(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEd74d837EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(out *jwriter.Writer, in managerHandlerDevice) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"tasks\":"
		out.RawString(prefix[1:])
		if in.Tasks == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v4, v5 := range in.Tasks {
				if v4 > 0 {
					out.RawByte(',')
				}
				if v5 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v6, v7 := range v5 {
						if v6 > 0 {
							out.RawByte(',')
						}
						out.String(string(v7))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Tags {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"serial_number\":"
		out.RawString(prefix)
		out.String(string(in.SerialNumber))
	}
	{
		const prefix string = ",\"mac\":"
		out.RawString(prefix)
		out.String(string(in.MAC))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"probe_readiness\":"
		out.RawString(prefix)
		out.String(string(in.ProbeReadiness))
	}
	{
		const prefix string = ",\"probe_liveness\":"
		out.RawString(prefix)
		out.String(string(in.ProbeLiveness))
	}
	{
		const prefix string = ",\"metrics_descriptions_count\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.MetricsDescriptionsCount))
	}
	{
		const prefix string = ",\"metrics_collect_count\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.MetricsCollectCount))
	}
	{
		const prefix string = ",\"metrics_empty_count\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.MetricsEmptyCount))
	}
	{
		const prefix string = ",\"mqtt_publishes\":"
		out.RawString(prefix)
		out.Int(int(in.MQTTPublishes))
	}
	{
		const prefix string = ",\"mqtt_subscribers\":"
		out.RawString(prefix)
		out.Int(int(in.MQTTSubscribers))
	}
	{
		const prefix string = ",\"logs_count\":"
		out.RawString(prefix)
		out.Int(int(in.LogsCount))
	}
	{
		const prefix string = ",\"has_metrics\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasMetrics))
	}
	{
		const prefix string = ",\"has_widget\":"
		out.RawString(prefix)
		out.Bool(bool(in.HasWidget))
	}
	{
		const prefix string = ",\"logs_max_level\":"
		out.RawString(prefix)
		out.RawText((in.LogsMaxLevel).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v managerHandlerDevice) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEd74d837EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v managerHandlerDevice) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEd74d837EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *managerHandlerDevice) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEd74d837DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *managerHandlerDevice) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEd74d837DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(l, v)
}
