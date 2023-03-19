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

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()

	client, err := b.client(ctx)
	if err != nil {
		widget.InternalError(w, r, err)
	}
	defer client.Close()

	var vars map[string]interface{}

	action := r.URL().Query().Get("action")

	switch action {
	case "logs":
		vars = b.widgetActionLogs(w, r, client)

	case "files":
		vars = b.widgetActionFiles(w, r, client)

	case "accounts":
		vars = b.widgetActionAccounts(w, r, client)

	case "user":
		vars = b.widgetActionUser(w, r, client)
		if vars == nil {
			return
		}

	case "password":
		vars = b.widgetActionPassword(w, r, client)
		if vars == nil {
			return
		}

	case "group":
		vars = b.widgetActionGroup(w, r, client)
		if vars == nil {
			return
		}

	case "user-delete":
		b.widgetActionUserDelete(w, r, client)
		return

	case "group-delete":
		b.widgetActionGroupDelete(w, r, client)
		return

	case "logs-export":
		b.widgetActionLogsExport(w, r, client)
		return

	case "configs-export":
		b.widgetActionConfigsExport(w, r, client)
		return

	case "download":
		b.widgetActionDownload(w, r, client)
		return

	case "preview":
		b.widgetActionPreview(w, r, client)
		return

	case "system":
		if r.IsPost() {
			b.widgetActionSystemPost(w, r, client)
			return
		}

		vars = b.widgetActionSystem(w, r, client)

	default:
		action = "default"

		vars = b.widgetActionDefault(w, r, client)
	}

	vars["action"] = action
	widget.RenderLayout(ctx, action, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) widgetActionDefault(_ *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	cfg := b.config()

	return map[string]interface{}{
		"preview_refresh_interval": cfg.PreviewRefreshInterval.Seconds(),
	}
}

func (b *Bind) widgetActionSystem(_ *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	widget := b.Widget()

	tm, err := client.OPTime(ctx)
	if err != nil {
		widget.FlashError(r, "Get current time failed with error %v", "", err)
	}

	info, err := client.SystemInfo(ctx)
	if err != nil {
		widget.FlashError(r, "Get system info failed with error %v", "", err)
	}

	storage, err := client.StorageInfo(ctx)
	if err != nil {
		widget.FlashError(r, "Get storage info failed with error %v", "", err)
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
		widget.FlashError(r, "Get work state failed with error %v", "", err)
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
		widget.FlashError(r, "Get channels title failed with error %v", "", err)
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

func (b *Bind) widgetActionSystemPost(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) {
	ctx := r.Context()

	err := r.Original().ParseForm()
	if err == nil {
		for key, value := range r.Original().PostForm {
			if key == "time" {
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
		_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
	} else {
		_ = w.SendJSON(boggart.NewResponseJSON().Success(b.Widget().Translate(ctx, "Config set success", "")))
	}
}

func (b *Bind) widgetActionLogs(_ *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()

	logs, err := client.LogSearch(ctx, time.Now().Add(-time.Hour), time.Now(), 0)
	if err != nil {
		b.Widget().FlashError(r, "Get logs failed with error %v", "", err)
	}

	return map[string]interface{}{
		"logs": logs,
	}
}

func (b *Bind) widgetActionAccounts(_ *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	widget := b.Widget()

	users, err := client.Users(ctx)
	if err != nil {
		widget.FlashError(r, "Get failed failed with error %v", "", err)
	}

	groups, err := client.Groups(ctx)
	if err != nil {
		widget.FlashError(r, "Create groups failed with error %v", "", err)
	}

	viewGroups := make([]struct {
		xmeye.Group
		CanRemove bool
	}, len(groups))

	for i, group := range groups {
		viewGroups[i].Group = group
		viewGroups[i].CanRemove = true

		for _, user := range users {
			if user.Group == group.Name {
				viewGroups[i].CanRemove = false
				break
			}
		}
	}

	return map[string]interface{}{
		"users":   users,
		"groups":  viewGroups,
		"current": b.config().Address.User.Username(),
	}
}

func (b *Bind) widgetActionUser(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	widget := b.Widget()

	var user *xmeye.User

	username := strings.TrimSpace(r.URL().Query().Get("username"))
	if username != "" { // update
		users, err := client.Users(ctx)
		if err == nil {
			for _, u := range users {
				if u.Name == username {
					u := u
					user = &u

					break
				}
			}
		}

		if user == nil {
			widget.NotFound(w, r)
			return nil
		}

		if r.IsPost() {
			user.Name = r.Original().FormValue("name")
			user.Memo = r.Original().FormValue("memo")
			user.Group = r.Original().FormValue("group")
			user.AuthorityList = r.Original().Form["authorities"]

			if err := client.UserUpdate(ctx, username, *user); err != nil {
				widget.FlashError(r, "Update user %s failed with error %v", "", username, err)
			} else {
				widget.FlashError(r, "Update user %s success", "", username)
				widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
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
			user.AuthorityList = r.Original().Form["authorities"]

			if err := client.UserCreate(ctx, *user); err != nil {
				widget.FlashError(r, "Create user failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "Create user %s success", "", user.Name)
				widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		} else if groupName := strings.TrimSpace(r.URL().Query().Get("groupname")); groupName != "" {
			user.Group = groupName
		}
	}

	vars := map[string]interface{}{
		"user":     user,
		"username": username,
	}

	if groups, err := client.Groups(ctx); err != nil {
		widget.FlashError(r, "Get groups failed with error %v", "", err)
	} else {
		vars["groups"] = groups
	}

	if authorities, err := client.FullAuthorityList(ctx); err != nil {
		widget.FlashError(r, "Get authorities list failed with error %v", "", err)
	} else {
		list := make([]struct {
			Name      string
			IsChecked bool
		}, len(authorities))

		for i, name := range authorities {
			list[i].Name = name

			for _, a := range user.AuthorityList {
				if a == name {
					list[i].IsChecked = true
					break
				}
			}
		}

		vars["authorities"] = list
	}

	return vars
}

func (b *Bind) widgetActionPassword(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	widget := b.Widget()

	userName := strings.TrimSpace(r.URL().Query().Get("username"))
	if userName == "" {
		widget.NotFound(w, r)
		return nil
	}

	if r.IsPost() {
		ctx := r.Context()
		oldPassword := r.Original().FormValue("old")
		newPassword := r.Original().FormValue("new")

		if err := client.UserChangePassword(ctx, userName, oldPassword, newPassword); err != nil {
			widget.FlashError(r, "Change user %s password failed with error %v", "", userName, err)
		} else {
			widget.FlashError(r, "Change user %s password success", "", userName)
			widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
			return nil
		}
	}

	return map[string]interface{}{}
}

func (b *Bind) widgetActionGroup(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	ctx := r.Context()
	widget := b.Widget()
	canEditAuthorities := true

	var group *xmeye.Group

	groupName := strings.TrimSpace(r.URL().Query().Get("groupname"))
	if groupName != "" { // update
		groups, err := client.Groups(ctx)
		if err == nil {
			for _, u := range groups {
				if u.Name == groupName {
					u := u
					group = &u

					break
				}
			}
		}

		if group == nil {
			widget.NotFound(w, r)
			return nil
		}

		if users, err := client.Users(ctx); err != nil {
			widget.FlashError(r, "Get users failed with error %v", "", err)
		} else {
			for _, user := range users {
				if user.Group == groupName {
					canEditAuthorities = false
					break
				}
			}
		}

		if r.IsPost() {
			group.Name = r.Original().FormValue("name")
			group.Memo = r.Original().FormValue("memo")

			if canEditAuthorities {
				group.AuthorityList = r.Original().Form["authorities"]
			}

			if err := client.GroupUpdate(ctx, groupName, *group); err != nil {
				widget.FlashError(r, "Update group %s failed with error %v", "", groupName, err)
			} else {
				widget.FlashSuccess(r, "Update group %s success", "", groupName)
				widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	} else { // create
		group = &xmeye.Group{}

		if r.IsPost() {
			group.Name = r.Original().FormValue("name")
			group.Memo = r.Original().FormValue("memo")
			group.AuthorityList = r.Original().Form["authorities"]

			if err := client.GroupCreate(ctx, *group); err != nil {
				widget.FlashError(r, "Create group failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "Create group %s success", "", group.Name)
				widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
				return nil
			}
		}
	}

	vars := map[string]interface{}{
		"group":              group,
		"groupname":          groupName,
		"canEditAuthorities": canEditAuthorities,
	}

	if canEditAuthorities {
		if authorities, err := client.FullAuthorityList(ctx); err != nil {
			widget.FlashError(r, "Get authorities list failed with error %v", "", err)
		} else {
			list := make([]struct {
				Name      string
				IsChecked bool
			}, len(authorities))

			for i, name := range authorities {
				list[i].Name = name

				for _, a := range group.AuthorityList {
					if a == name {
						list[i].IsChecked = true
						break
					}
				}
			}

			vars["authorities"] = list
		}
	}

	return vars
}

func (b *Bind) widgetActionUserDelete(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) {
	widget := b.Widget()

	userName := strings.TrimSpace(r.URL().Query().Get("username"))
	if userName == "" {
		widget.NotFound(w, r)
		return
	}

	if !r.IsPost() {
		widget.MethodNotAllowed(w, r)
		return
	}

	ctx := r.Context()

	if err := client.UserDelete(ctx, userName); err != nil {
		widget.FlashError(r, "Remove user %s failed with error %v", "", userName, err)
	} else {
		widget.FlashSuccess(r, "Remove user %s success", "", userName)
	}

	widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
}

func (b *Bind) widgetActionGroupDelete(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) {
	widget := b.Widget()

	groupName := strings.TrimSpace(r.URL().Query().Get("groupname"))
	if groupName == "" {
		widget.NotFound(w, r)
		return
	}

	if !r.IsPost() {
		widget.MethodNotAllowed(w, r)
		return
	}

	ctx := r.Context()

	if err := client.GroupDelete(ctx, groupName); err != nil {
		widget.FlashError(r, "Remove group %s failed with error %v", "", groupName, err)
	} else {
		widget.FlashSuccess(r, "Remove group %s success", "", groupName)
	}

	widget.Redirect(r.URL().Path+"?action=accounts", http.StatusFound, w, r)
}

func (b *Bind) widgetActionFiles(_ *dashboard.Response, r *dashboard.Request, client *xmeye.Client) map[string]interface{} {
	query := r.URL().Query()
	ctx := r.Context()
	widget := b.Widget()
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
			widget.FlashError(r, "Unknown event type %s", "", et)
		}
	}

	var (
		channelID    uint32
		channelTitle string
	)

	if raw := query.Get("channel"); raw != "" {
		if val, err := strconv.ParseUint(raw, 10, 64); err == nil {
			channelID = uint32(val)
		} else {
			widget.FlashError(r, "Parse channel ID failed with error %v", "", err)
		}
	}

	if raw := query.Get("from"); raw != "" {
		if val, err := time.Parse(time.RFC3339, raw); err == nil {
			start = val
		} else {
			widget.FlashError(r, "Parse date from failed with error %v", "", err)
		}
	}

	if raw := query.Get("to"); raw != "" {
		if val, err := time.Parse(time.RFC3339, raw); err == nil {
			end = val
		} else {
			widget.FlashError(r, "Parse date to failed with error %v", "", err)
		}
	}

	channels, err := client.ConfigChannelTitleGet(ctx)
	if err != nil {
		widget.FlashError(r, "Get channels title failed with error %v", "", err)
	} else {
		for id, title := range channels {
			if uint32(id) == channelID {
				channelTitle = title
				break
			}
		}
	}

	files := make([]xmeye.FileSearch, 0)

	filesH264, err := client.FileSearch(ctx, start, end, channelID, eventType, xmeye.FileSearchH264)
	if err != nil {
		widget.FlashError(r, "Get files H264 failed with error %v", "", err)
	} else {
		files = append(files, filesH264...)
	}

	filesJPEG, err := client.FileSearch(ctx, start, end, channelID, eventType, xmeye.FileSearchJPEG)
	if err != nil {
		widget.FlashError(r, "Get files JPEG failed with error %v", "", err)
	} else {
		files = append(files, filesJPEG...)
	}

	for i := range files {
		files[i].FileLength *= 1024
	}

	return map[string]interface{}{
		"event_type":    eventType,
		"channels":      channels,
		"channel_id":    channelID,
		"channel_title": channelTitle,
		"date_from":     start,
		"date_to":       end,
		"files":         files,
	}
}

func (b *Bind) widgetActionLogsExport(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) {
	ctx := r.Context()

	reader, err := client.LogExport(ctx)
	if err != nil {
		b.Widget().FlashError(r, "Export logs failed with error %v", "", err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")

	filename := b.Meta().ID() + time.Now().Format("_logs_20060102150405.zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	_, _ = io.Copy(w, reader)
}

func (b *Bind) widgetActionConfigsExport(w *dashboard.Response, r *dashboard.Request, client *xmeye.Client) {
	ctx := r.Context()

	reader, err := client.ConfigExport(ctx)
	if err != nil {
		b.Widget().FlashError(r, "Export config failed with error %v", "", err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")

	filename := b.Meta().ID() + time.Now().Format("_config_20060102150405.zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	_, _ = io.Copy(w, reader)
}

func (b *Bind) widgetActionDownload(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) {
	widget := b.Widget()

	name := strings.TrimSpace(r.URL().Query().Get("name"))
	if name == "" {
		widget.NotFound(w, r)
		return
	}

	filename := strings.TrimSpace(r.URL().Query().Get("filename"))
	if filename == "" {
		filename = name
	}

	begin := time.Now().Add(-time.Hour * 24 * 30)
	end := time.Now()

	reader, err := client.PlayStream(r.Context(), begin, end, name)
	if err != nil {
		widget.NotFound(w, r)
		return
	}

	switch strings.ToLower(filepath.Ext(name)) {
	case ".h264":
		w.Header().Set("Content-Type", "video/H264")
	case ".jpeg", ".jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	_, _ = io.Copy(w, reader)
}

func (b *Bind) widgetActionPreview(w http.ResponseWriter, r *dashboard.Request, client *xmeye.Client) {
	var (
		ch  uint64
		err error
	)

	query := r.URL().Query()
	ctx := r.Context()
	widget := b.Widget()
	cfg := b.config()

	if channel := query.Get("channel"); channel == "" {
		ch = cfg.WidgetChannel
	} else {
		ch, err = strconv.ParseUint(channel, 10, 64)
		if err != nil {
			widget.NotFound(w, r)
			return
		}
	}

	var reader io.Reader
	if cfg.PreviewUseRTSP {
		reader, err = client.SnapshotRTSP(ctx, ch)
	} else {
		reader, err = client.Snapshot(ctx, ch)
	}

	if err != nil {
		widget.InternalError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	if download := query.Get("download"); download != "" {
		filename := b.Meta().ID() + time.Now().Format("_20060102150405.jpg")

		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	}

	_, _ = io.Copy(w, reader)
}
