// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package handlers

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(in *jlexer.Lexer, out *deviceHandlerListener) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "events":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Events = make(map[string]string)
				} else {
					out.Events = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.Events)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "fires":
			out.Fires = int64(in.Int64())
		case "fire_first":
			if in.IsNull() {
				in.Skip()
				out.FiredFirst = nil
			} else {
				if out.FiredFirst == nil {
					out.FiredFirst = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.FiredFirst).UnmarshalJSON(data))
				}
			}
		case "fire_last":
			if in.IsNull() {
				in.Skip()
				out.FiredLast = nil
			} else {
				if out.FiredLast == nil {
					out.FiredLast = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.FiredLast).UnmarshalJSON(data))
				}
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
func easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(out *jwriter.Writer, in deviceHandlerListener) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"events\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Events == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Events {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"fires\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Fires))
	}
	{
		const prefix string = ",\"fire_first\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.FiredFirst == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.FiredFirst).MarshalJSON())
		}
	}
	{
		const prefix string = ",\"fire_last\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.FiredLast == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.FiredLast).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v deviceHandlerListener) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v deviceHandlerListener) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *deviceHandlerListener) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *deviceHandlerListener) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers(l, v)
}
func easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(in *jlexer.Lexer, out *deviceHandlerDevice) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "register_id":
			out.RegisterId = string(in.String())
		case "id":
			out.Id = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "serial_number":
			out.SerialNumber = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "tasks":
			if in.IsNull() {
				in.Skip()
				out.Tasks = nil
			} else {
				in.Delim('[')
				if out.Tasks == nil {
					if !in.IsDelim(']') {
						out.Tasks = make([]string, 0, 4)
					} else {
						out.Tasks = []string{}
					}
				} else {
					out.Tasks = (out.Tasks)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.Tasks = append(out.Tasks, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "mqtt_topics":
			if in.IsNull() {
				in.Skip()
				out.MQTTTopics = nil
			} else {
				in.Delim('[')
				if out.MQTTTopics == nil {
					if !in.IsDelim(']') {
						out.MQTTTopics = make([]string, 0, 4)
					} else {
						out.MQTTTopics = []string{}
					}
				} else {
					out.MQTTTopics = (out.MQTTTopics)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.MQTTTopics = append(out.MQTTTopics, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "mqtt_subscribers":
			if in.IsNull() {
				in.Skip()
				out.MQTTSubscribers = nil
			} else {
				in.Delim('[')
				if out.MQTTSubscribers == nil {
					if !in.IsDelim(']') {
						out.MQTTSubscribers = make([]string, 0, 4)
					} else {
						out.MQTTSubscribers = []string{}
					}
				} else {
					out.MQTTSubscribers = (out.MQTTSubscribers)[:0]
				}
				for !in.IsDelim(']') {
					var v5 string
					v5 = string(in.String())
					out.MQTTSubscribers = append(out.MQTTSubscribers, v5)
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
					var v6 string
					v6 = string(in.String())
					out.Tags = append(out.Tags, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "config":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Config = make(map[string]interface{})
				} else {
					out.Config = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v7 interface{}
					if m, ok := v7.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v7.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v7 = in.Interface()
					}
					(out.Config)[key] = v7
					in.WantComma()
				}
				in.Delim('}')
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
func easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(out *jwriter.Writer, in deviceHandlerDevice) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"register_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RegisterId))
	}
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"serial_number\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SerialNumber))
	}
	{
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"tasks\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Tasks == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Tasks {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"mqtt_topics\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.MQTTTopics == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v10, v11 := range in.MQTTTopics {
				if v10 > 0 {
					out.RawByte(',')
				}
				out.String(string(v11))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"mqtt_subscribers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.MQTTSubscribers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v12, v13 := range in.MQTTSubscribers {
				if v12 > 0 {
					out.RawByte(',')
				}
				out.String(string(v13))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"tags\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Tags {
				if v14 > 0 {
					out.RawByte(',')
				}
				out.String(string(v15))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"config\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Config == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v16First := true
			for v16Name, v16Value := range in.Config {
				if v16First {
					v16First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v16Name))
				out.RawByte(':')
				if m, ok := v16Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v16Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v16Value))
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v deviceHandlerDevice) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v deviceHandlerDevice) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson65411fd3EncodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *deviceHandlerDevice) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *deviceHandlerDevice) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson65411fd3DecodeGithubComKihamoBoggartComponentsBoggartInternalHandlers1(l, v)
}
