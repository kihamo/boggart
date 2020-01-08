package xmeye

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
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
	client, err := b.Bind().(*Bind).client()
	if err != nil {
		t.InternalError(w, r, err)
	}
	defer client.Close()

	var vars map[string]interface{}

	action := r.URL().Query().Get("action")

	switch action {
	case "logs":
		vars = t.widgetActionLogs(w, r, client)

	case "files":
		vars = t.widgetActionFiles(w, r, client)

	case "accounts":
		vars = t.widgetActionAccounts(w, r, client, b)

	case "user":
		vars = t.widgetActionUser(w, r, client)
		if len(vars) == 0 {
			return
		}

	case "group":
		vars = t.widgetActionGroup(w, r, client)
		if len(vars) == 0 {
			return
		}

	case "logs-export":
		t.widgetActionLogsExport(w, r, client, b)
		return

	case "configs-export":
		t.widgetActionConfigsExport(w, r, client, b)
		return

	case "download":
		t.widgetActionDownload(w, r, client)
		return

	default:
		if r.IsPost() {
			t.widgetActionDefaultPost(w, r, client)
			return
		}

		vars = t.widgetActionDefault(w, r, client)
		action = "default"
	}

	vars["action"] = action

	t.RenderLayout(r.Context(), action, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (t Type) widgetActionDefault(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()

	tm, err := client.OPTime(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get current time failed with error %v", "", err))
	}

	info, err := client.SystemInfo(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get system info failed with error %v", "", err))
	}

	storage, err := client.StorageInfo(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get storage info failed with error %v", "", err))
	} else {
		for s, st := range storage {
			for p := range st.Partition {
				storage[s].Partition[p].RemainSpace *= 1024
				storage[s].Partition[p].TotalSpace *= 1024
			}
		}
	}

	type chanInfo struct {
		Number  int
		Title   string
		Bitrate uint64
		Record  bool
	}
	channels := make([]chanInfo, client.ChannelsCount())

	state, err := client.WorkState(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get work state failed with error %v", "", err))
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

	titles, err := client.ConfigChannelTitleGet(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get channels title failed with error %v", "", err))
	} else {
		for i, title := range titles {
			channels[i].Number = i + 1
			channels[i].Title = title

			if i == len(channels)-1 {
				break
			}
		}
	}

	return map[string]interface{}{
		"time_system":  time.Now(),
		"time_current": tm,
		"system_info":  info,
		"storage_info": storage,
		"channels":     channels,
	}
}

func (t Type) widgetActionDefaultPost(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) {
	ctx := r.Context()

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
					err = client.OPTimeSetting(ctx, t)
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
}

func (t Type) widgetActionLogs(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()

	logs, err := client.LogSearch(ctx, time.Now().Add(-time.Hour), time.Now(), 0)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get logs failed with error %v", "", err))
	}

	return map[string]interface{}{
		"logs": logs,
	}
}

func (t Type) widgetActionAccounts(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client, b boggart.BindItem) map[string]interface{} {
	ctx := r.Context()

	users, err := client.Users(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get failed failed with error %v", "", err))
	}

	groups, err := client.Groups(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Create groups failed with error %v", "", err))
	}

	return map[string]interface{}{
		"users":   users,
		"groups":  groups,
		"current": b.Bind().(*Bind).config.Address.User.Username(),
	}
}

func (t Type) widgetActionUser(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	var user *xmeye.User

	username := strings.TrimSpace(r.URL().Query().Get("username"))
	if username != "" { // update
		users, err := client.Users(ctx)
		if err == nil {
			for _, u := range users {
				if u.Name == username {
					user = &u
					break
				}
			}
		}

		if user == nil {
			t.NotFound(w, r)
			return nil
		}

		if r.IsPost() {
			user.Name = r.Original().FormValue("name")
			user.Memo = r.Original().FormValue("memo")
			user.Group = r.Original().FormValue("group")

			if err := client.UserUpdate(ctx, username, *user); err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Update user %s failed with error %v", "", username, err))
			} else {
				r.Session().FlashBag().Success(t.Translate(ctx, "Update user %s success", "", username))
				t.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	} else { // create
		user = &xmeye.User{}

		if r.IsPost() {
			user.Name = r.Original().FormValue("name")
			user.Memo = r.Original().FormValue("memo")
			user.Password = r.Original().FormValue("password")
			user.Group = r.Original().FormValue("group")

			if err := client.UserCreate(ctx, *user); err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Create user failed with error %v", "", err))
			} else {
				r.Session().FlashBag().Success(t.Translate(ctx, "Create user %s success", "", user.Name))
				t.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	}

	vars := map[string]interface{}{
		"user": user,
	}

	if groups, err := client.Groups(ctx); err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Create groups failed with error %v", "", err))
	} else {
		vars["groups"] = groups
	}

	return vars
}

func (t Type) widgetActionGroup(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	var group *xmeye.Group

	groupName := strings.TrimSpace(r.URL().Query().Get("groupname"))
	if groupName != "" { // update
		groups, err := client.Groups(ctx)
		if err == nil {
			for _, u := range groups {
				if u.Name == groupName {
					group = &u
					break
				}
			}
		}

		if group == nil {
			t.NotFound(w, r)
			return nil
		}

		if r.IsPost() {
			group.Name = r.Original().FormValue("name")
			group.Memo = r.Original().FormValue("memo")

			if err := client.GroupUpdate(ctx, groupName, *group); err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Update group %s failed with error %v", "", groupName, err))
			} else {
				r.Session().FlashBag().Success(t.Translate(ctx, "Update group %s success", "", groupName))
				t.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	} else { // create
		group = &xmeye.Group{}

		if r.IsPost() {
			group.Name = r.Original().FormValue("name")
			group.Memo = r.Original().FormValue("memo")

			if err := client.GroupCreate(ctx, *group); err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Create group failed with error %v", "", err))
			} else {
				r.Session().FlashBag().Success(t.Translate(ctx, "Create group %s success", "", group.Name))
				t.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	}

	return map[string]interface{}{
		"group": group,
	}
}

func (t Type) widgetActionFiles(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	query := r.URL().Query()
	ctx := r.Context()

	var channel uint32
	eventType := xmeye.FileSearchEventAll
	end := time.Now()
	start := end.Add(-time.Hour * 24 * 30)

	if et := query.Get("event-type"); et != "" {
		switch strings.ToUpper(et) {
		case "*":
			eventType = xmeye.FileSearchEventAll
		case "A":
			eventType = xmeye.FileSearchEventAlarm
		case "M":
			eventType = xmeye.FileSearchEventMotionDetect
		case "R":
			eventType = xmeye.FileSearchEventGeneral
		case "H":
			eventType = xmeye.FileSearchEventManual
		default:
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Unknown event type %s", "", et))
		}
	}

	if channelID := query.Get("channel"); channelID != "" {
		if cID, err := strconv.ParseUint(channelID, 10, 64); err == nil {
			channel = uint32(cID)
		} else {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse channel ID failed with error %v", "", err))
		}
	}

	if channelID := query.Get("channel"); channelID != "" {
		if cID, err := strconv.ParseUint(channelID, 10, 64); err == nil {
			channel = uint32(cID)
		} else {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse channel ID failed with error %v", "", err))
		}
	}

	if queryTime := query.Get("from"); queryTime != "" {
		if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
			start = tm
		} else {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse date from failed with error %v", "", err))
		}
	}

	if queryTime := query.Get("to"); queryTime != "" {
		if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
			end = tm
		} else {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse date to failed with error %v", "", err))
		}
	}

	channels, err := client.ConfigChannelTitleGet(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get channels title failed with error %v", "", err))
	}

	files := make([]xmeye.FileSearch, 0)

	filesH264, err := client.FileSearch(ctx, start, end, channel, eventType, xmeye.FileSearchH264)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get files H264 failed with error %v", "", err))
	} else {
		files = append(files, filesH264...)
	}

	filesJPEG, err := client.FileSearch(ctx, start, end, channel, eventType, xmeye.FileSearchJPEG)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Get files JPEG failed with error %v", "", err))
	} else {
		files = append(files, filesJPEG...)
	}

	for i := range files {
		files[i].FileLength = files[i].FileLength * 1024
	}

	return map[string]interface{}{
		"event_type": eventType,
		"channels":   channels,
		"channel":    channel,
		"date_from":  start,
		"date_to":    end,
		"files":      files,
	}
}

func (t Type) widgetActionLogsExport(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client, b boggart.BindItem) {
	ctx := r.Context()

	reader, err := client.LogExport(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Export logs failed with error %v", "", err))
		return
	}

	w.Header().Set("Content-Type", "application/zip")

	filename := b.ID() + time.Now().Format("_logs_20060102150405.zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	_, _ = io.Copy(w, reader)
}

func (t Type) widgetActionConfigsExport(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client, b boggart.BindItem) {
	ctx := r.Context()

	reader, err := client.ConfigExport(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx, "Export config failed with error %v", "", err))
		return
	}

	w.Header().Set("Content-Type", "application/zip")

	filename := b.ID() + time.Now().Format("_config_20060102150405.zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	_, _ = io.Copy(w, reader)
}

func (t Type) widgetActionDownload(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) {
	name := strings.TrimSpace(r.URL().Query().Get("name"))
	if name == "" {
		t.NotFound(w, r)
		return
	}

	begin := time.Now().Add(-time.Hour * 24 * 30)
	end := time.Now()

	reader, err := client.PlayStream(r.Context(), begin, end, name)
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
}
