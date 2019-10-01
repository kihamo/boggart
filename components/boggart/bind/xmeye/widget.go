package xmeye

import (
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xmeye"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	query := r.URL().Query()
	action := query.Get("action")
	ctx := r.Context()

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	case "logs":
		logs, err := bind.client.LogSearch(ctx, time.Now().Add(-time.Hour), time.Now(), 0)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get logs failed with error %s", "", err.Error()))
		}

		vars["logs"] = logs

	case "logs-export":
		reader, err := bind.client.LogExport(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Export logs failed with error %s", "", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/zip")

		filename := b.ID() + time.Now().Format("_logs_20060102150405.zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		_, _ = io.Copy(w, reader)

		return

	case "configs-export":
		reader, err := bind.client.ConfigExport(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Export config failed with error %s", "", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/zip")

		filename := b.ID() + time.Now().Format("_config_20060102150405.zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		_, _ = io.Copy(w, reader)

		return

	case "download":
		name := strings.TrimSpace(query.Get("name"))
		if name == "" {
			t.NotFound(w, r)
			return
		}

		begin := time.Now().Add(-time.Hour * 24 * 30)
		end := time.Now()

		reader, err := bind.client.PlayStream(ctx, begin, end, name)
		if err != nil {
			t.NotFound(w, r)
			return
		}

		switch strings.ToLower(filepath.Ext(name)) {
		case ".h264":
			w.Header().Set("Content-Type", "video/H264")
		case ".jpeg", ".jpg":
			w.Header().Set("Content-Type", "image/jpeg")
		}

		w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
		_, _ = io.Copy(w, reader)

		return

	case "files":
		var channel uint32 = 1
		begin := time.Now().Add(-time.Hour * 24 * 30)
		end := time.Now()
		eventType := xmeye.FileSearchEventAll

		files := make([]xmeye.FileSearch, 0)

		filesH264, err := bind.client.FileSearch(ctx, begin, end, channel, eventType, xmeye.FileSearchH264)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get files H264 failed with error %s", "", err.Error()))
		} else {
			files = append(files, filesH264...)
		}

		filesJPEG, err := bind.client.FileSearch(ctx, begin, end, channel, eventType, xmeye.FileSearchJPEG)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get files JPEG failed with error %s", "", err.Error()))
		} else {
			files = append(files, filesJPEG...)
		}

		for i := range files {
			files[i].FileLength = files[i].FileLength * 1024
		}

		vars["files"] = files

	default:
		if r.IsPost() {
			err := r.Original().ParseForm()
			if err == nil {
				for key, value := range r.Original().PostForm {
					switch key {
					case "time":
						var t time.Time

						if strings.Compare(value[0], "now") == 0 {
							t = time.Now()
						} else {
							t, err = time.Parse("2006.01.02 15:04:05", value[0])
						}

						if err == nil {
							err = bind.client.OPTimeSetting(ctx, t)
						}
					}
				}
			}

			if err != nil {
				_ = w.SendJSON(response{
					Result:  "failed",
					Message: err.Error(),
				})

			} else {
				_ = w.SendJSON(response{
					Result:  "success",
					Message: t.Translate(ctx, "Config set success", ""),
				})
			}

			return
		}

		tm, err := bind.client.OPTime(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get current time failed with error %s", "", err.Error()))
		} else {
			vars["time_current"] = tm
		}
		vars["time_system"] = time.Now()

		info, err := bind.client.SystemInfo(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get system info failed with error %s", "", err.Error()))
		}
		vars["system_info"] = info

		storage, err := bind.client.StorageInfo(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get storage info failed with error %s", "", err.Error()))
		} else {
			for s, st := range storage {
				for p := range st.Partition {
					storage[s].Partition[p].RemainSpace *= 1024
					storage[s].Partition[p].TotalSpace *= 1024
				}
			}
		}
		vars["storage_info"] = storage

		type chanInfo struct {
			Number  int
			Title   string
			Bitrate uint64
			Record  bool
		}
		channels := make([]chanInfo, bind.client.ChannelsCount())

		state, err := bind.client.WorkState(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get work state failed with error %s", "", err.Error()))
		} else {
			for i, state := range state.ChannelState {
				channels[i].Number = i + 1
				channels[i].Bitrate = state.Bitrate
				channels[i].Record = state.Record

				if i == len(channels)-1 {
					break
				}
			}
		}

		titles, err := bind.client.ConfigChannelTitleGet(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get channels title failed with error %s", "", err.Error()))
		} else {
			for i, title := range titles {
				channels[i].Number = i + 1
				channels[i].Title = title

				if i == len(channels)-1 {
					break
				}
			}
		}

		vars["channels"] = channels
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
