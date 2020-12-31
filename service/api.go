package adash

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"bitbucket.org/peregrinetraders/mc/apps/adash/service/charter"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/holidayer"
	"bitbucket.org/peregrinetraders/mc/apps/adash/service/redishashtracker"
	"bitbucket.org/peregrinetraders/mc/pkg/aaa"
	"bitbucket.org/peregrinetraders/mc/pkg/random"
	"github.com/axkit/hms"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/hydrogen18/stoppableListener"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"strconv"
)

type API struct {
	log                  zerolog.Logger
	aaaService           *aaa.Service
	router               *chi.Mux
	service              Servicer
	wwwfolder            string
	restartAllowedBefore *time.Time
	restartShellCommand  string
	pnlch                *charter.Charter
	appVersion           AppVersion
	h                    *holidayer.Holidayer
}

type AppVersion struct {
	Version      string `json:"version"`
	BuildAt      string `json:"buildAt"`
	BuildFor     string `json:"buildFor"`
	ReleaseNotes string `json:"releaseNotes"`
}

var TokenAuth *jwtauth.JWTAuth

func init() {
	// TODO: выташить потом секрет ключ?
	TokenAuth = jwtauth.New("HS256", []byte("secret-top-secret"), nil)
}

func NewAPI(service Servicer, a *aaa.Service, looger zerolog.Logger,
	wwwfolder string, rab hms.HHMM, rsc string, av AppVersion,
	h *holidayer.Holidayer) *API {
	res := API{
		service:             service,
		aaaService:          a,
		log:                 looger,
		wwwfolder:           wwwfolder,
		restartShellCommand: strings.TrimSpace(rsc),
		appVersion:          av,
		h:                   h,
	}

	if len(rab) == 5 {
		t := time.Now().UTC().Truncate(time.Hour * 24).Add(rab.MidnightOffset())
		res.restartAllowedBefore = &t
	}

	return &res
}

func (a *API) Start(ctx context.Context, listen string, lp *redishashtracker.TrackingEngine, pnlch *charter.Charter) error {

	a.pnlch = pnlch
	// create http router
	a.router = chi.NewRouter()
	//Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	//f.router.Use(middleware.DefaultLogger)
	a.router.Use(middleware.Compress(5, "text/plain"))
	a.router.Use(cors.Handler)

	/*f.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "files")
	FileServer(root, "/", "frontend", http.Dir(filesDir))
	/**/

	//f.router.Get("/ws", f.WS.Handler())
	//f.router.Get("/ws/", f.WS.Handler())
	a.router.Get("/api/matrix/i", a.getMatrixIHandler)
	a.router.Get("/api/matrix/m", a.getMatrixHandler)
	a.router.Post("/api/upload", a.uploadHandler)
	a.router.Get("/api/config", a.getConfigHandler)
	a.router.Get("/api/version", a.getVersion)
	a.router.Get("/api/holidays", a.getHolidaysHandler)

	a.router.Get("/api/log/{id}/count", a.getLogCountHandler)
	a.router.Get("/api/log/{id}/feed", a.getLogFeedHandler)
	a.router.Get("/api/log/{id}/orders", a.getLogFeedOrdersHandler)
	a.router.Get("/api/log/{id}/orders/{orderid}", a.getLogFeedExecutionsHandler)
	a.router.Get("/api/instance/{id}/info", a.getInfoHandler)
	a.router.Post("/api/instance/{id}/stop", a.postCloseHandler)
	a.router.Post("/api/close-all", a.postCloseAllHandler)
	a.router.Get("/api/platform/status", a.getPlatformStatusHandler)
	a.router.Post("/api/platform/restart", a.postRestartPlatformHandler)

	a.router.Get("/api/log/{id}/download", a.downloadFile)

	a.router.Get("/api/charter/{id}/img", pnlch.MultiChartImage)
	a.router.Get("/api/charter/{id}/array", pnlch.SingleChartJSON)
	a.router.Get("/api/charter/{id}/multiarray", pnlch.MultiChartJSON)
	a.router.Get("/api/charter/{id}/download", pnlch.DownloadFile)

	a.router.Get("/api/pnl/{id}/multiarray", a.getPnLMultiarrayHandler)
	a.router.Get("/api/portfolio/group-pnl", pnlch.GroupPnLHandler)
	a.router.Get("/api/portfolio/debug", pnlch.DebugHandler)
	//f.router.Get("/api/instance/{id}/lpchart", lp.ChartHandler)

	/*f.router.Post("/api/auth/signout", f.signOutHandler)
	f.router.Get("/api/auth/issigned", f.isSignedHandler)
	f.router.Get("/api/auth/user", f.userDataHandler)
	*/

	a.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		var fname = r.RequestURI
		if r.RequestURI == "/" {
			fname = "/index.html"
		}

		fname = a.wwwfolder + fname
		//fmt.Printf("requestURI: %s, fname=%s\n", r.RequestURI, fname)

		buf, err := ioutil.ReadFile(fname)
		if err != nil {
			a.log.Error().Str("errmsg", err.Error()).Msg("index.html reading failed")
			return
		}
		if strings.HasPrefix(r.RequestURI, "/css") {
			w.Header().Set("Content-Type", "text/css")
		}
		if strings.HasPrefix(r.RequestURI, "/js") {
			w.Header().Set("Content-Type", "text/javascript")
		}

		if r.RequestURI == "/" {
			w.Header().Set("Content-Type", "text/html")
		}

		w.Write(buf)

	})

	if err := a.startHTTP(ctx, listen); err != nil {
		return errors.New("listener starting failed. " + err.Error())
	}

	//	f.startWSProxy(ctx, &f.Log, f.WS)

	return nil
}

func fail(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write([]byte(`{"msg": "` + msg + `"}`))
}

func FileServer(r chi.Router, basePath string, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(basePath+path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (a *API) startHTTP(ctx context.Context, listen string) error {

	originalListener, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	sl, err := stoppableListener.New(originalListener)
	if err != nil {
		return err
	}

	server := http.Server{Handler: a.router}

	go func() {

		a.log.Info().Str("listen", listen).Msg("start serving HTTP requests")

		go server.Serve(sl)
		<-ctx.Done()
		a.log.Info().Msg("Stopping listener")
		sl.Stop()
		a.log.Info().Msg("Waiting on server")

	}()

	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func bind(r io.Reader, data interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(data)
}

func (a *API) getMatrixHandler(w http.ResponseWriter, r *http.Request) {

	//sortBy := "id"
	// v, ok := r.URL.Query()["sort"]
	// if ok {
	// 	sortBy = v[0]
	// }

	var res []InstanceEx
	if err := a.service.Instances(&res); err != nil {
		render.Render(w, r, ErrRender(err))
	}
	render.JSON(w, r, res)

	/*	switch sortBy {
		case "id":
			sort.Slice(res, func(i, j int) bool {
				return res[i].ID < res[j].ID
			})
		case "pnl":
			sort.Slice(res, func(i, j int) bool {
				if res[i].PnlFloat == 0 && res[j].PnlFloat == 0 {
					return res[i].ID < res[j].ID
				}

				if res[i].PnlFloat == 0 {
					return false
				}
				if res[j].PnlFloat == 0 {
					return true
				}
				return res[i].PnlFloat < res[j].PnlFloat
			})
			// deleted in accordance to Floris request
			// mark the latest instance had traded before
			// the instances with zero PnL.
			// for i := range res {
			// 	res[i].IsBeforeZero = false
			// }
			// for i := len(res) - 1; i >= 0; i-- {
			// 	if res[i].PnlFloat > 0 {
			// 		res[i].IsBeforeZero = true
			// 		break
			// 	}
			// }

		case "status":
			sort.Slice(res, func(i, j int) bool {
				if res[i].Status.Status == res[j].Status.Status {
					return res[i].PnlFloat < res[j].PnlFloat
				}
				return res[i].Status.Status < res[j].Status.Status
			})

		}

		render.JSON(w, r, res)*/
}

func (a *API) getMatrixIHandler(w http.ResponseWriter, r *http.Request) {

	res, err := a.service.MatrixI()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.JSON(w, r, res)
}

type Response struct {
	Received bool `json:"received"`
}

func (a *API) uploadHandler(w http.ResponseWriter, r *http.Request) {

	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("file uploading failed.details " + err.Error()))
		return
	}
	name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])

	fname := ""
	ltsname := time.Now().Format("20060102-150405")
	if len(name) > 0 {
		fname = name[0]
		ltsname += "-" + random.GenString(random.Chars, 6) + "-" + name[0]
	}
	if len(name) > 1 {
		fname += "." + name[1]
		ltsname += "." + name[1]
	}

	action := r.FormValue("action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`http form parameter "action" is empty`))
		return
	}

	// Copy the file data to my buffer
	if _, err := io.Copy(&Buf, file); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`internal error. Buffer copy failed: ` + err.Error()))
		return
	}
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example

	if err := a.service.SaveFileToDisk(ltsname, Buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`launcher config file writing filed:` + err.Error()))
		return
	}

	var resp Response

	resp.Received, err = a.service.PublishLaunchConfig(action, fname, Buf.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`sending message to the redis channel failed`))
		return
	}

	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header

	go func() {
		time.Sleep(3 * time.Second)
		if err = a.service.RefreshConfig(); err != nil {
			a.log.Error().Str("errmsg", err.Error()).Msg("refresh CFG:* kets failed")
		}
	}()
	render.JSON(w, r, resp)
	return
}

func (a *API) postCloseHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	ok, err := a.service.Publish("STOP:" + id)
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}

	resp := Response{ok}
	render.JSON(w, r, resp)
}

func (a *API) postCloseAllHandler(w http.ResponseWriter, r *http.Request) {

	ok, err := a.service.Publish("STOP:ALL")
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}

	resp := Response{ok}
	render.JSON(w, r, resp)
}

func (a *API) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, a.service.Config())
}

func (a *API) getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	b := a.appVersion
	wrn, ok := r.URL.Query()["withoutReleaseNotes"]
	if ok && wrn[0] == "true" {
		b.ReleaseNotes = ""
	}
	render.JSON(w, r, b)
}

func (a *API) getLogCountHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	i, err := a.service.ListLen(id)
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}

	from, ok := r.URL.Query()["from"]

	type Response struct {
		Count int `json:"count"`
	}

	var fromi int64
	if ok {
		if fromi, err = strconv.ParseInt(from[0], 10, 64); err != nil {
			println(err.Error())
			return
		}
	}

	var resp = Response{Count: i - int(fromi)}

	render.JSON(w, r, &resp)
}

func (a *API) getLogFeedHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	fromstr, ok := r.URL.Query()["from"]
	var from int64

	if ok {
		if from, err = strconv.ParseInt(fromstr[0], 10, 64); err != nil {
			return
		}
	}

	tostr, ok := r.URL.Query()["to"]
	var to int64

	if ok {
		if to, err = strconv.ParseInt(tostr[0], 10, 64); err != nil {
			return
		}
	}

	id := chi.URLParam(r, "id")

	if to > from {
		// to aviod repeatable message in the log
		to--
	}

	res, err := a.service.ListRange(KeyPrefixLog+id, int(from), int(to))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	sep := ""
	buf := bytes.NewBufferString("[")
	for i := range res {
		row := strings.SplitN(res[i], ":", 3)
		buf.WriteString(sep)
		buf.WriteString("{\"at\":")
		buf.WriteString(row[0])
		buf.WriteString(",\"severity\":")
		buf.WriteString(row[1])
		buf.WriteString(",\"msg\":\"")
		row[2] = strings.ReplaceAll(row[2], "\n", "")
		row[2] = strings.ReplaceAll(row[2], "\r", "")
		row[2] = strings.ReplaceAll(row[2], "\t", " ")
		row[2] = strings.ReplaceAll(row[2], "\"", "'")
		row[2] = strings.ReplaceAll(row[2], "\\\"", "'")
		row[2] = strings.TrimSpace(row[2])
		buf.WriteString(row[2])
		buf.WriteString("\"}")
		sep = ",\n"
	}
	buf.WriteString("]\n")

	//render.Data(w, r, buf.Bytes())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(buf.Bytes())
	buf.Reset()
}

func (a *API) getLogFeedOrdersHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	res, err := a.service.ListRange(KeyPrefixLog+id, 0, -1)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	sep := ""
	buf := bytes.NewBufferString("[")
	for i := range res {
		if !strings.Contains(res[i], "[Order:") {
			continue
		}

		row := strings.SplitN(res[i], ":", 3)
		buf.WriteString(sep)
		buf.WriteString("{\"at\":")
		buf.WriteString(row[0])
		buf.WriteString(",\"msg\":\"")
		row[2] = strings.ReplaceAll(row[2], "\n", "")
		row[2] = strings.ReplaceAll(row[2], "\r", "")
		row[2] = strings.ReplaceAll(row[2], "\t", " ")
		row[2] = strings.ReplaceAll(row[2], "\"", "'")
		row[2] = strings.ReplaceAll(row[2], "\\\"", "'")
		row[2] = strings.TrimSpace(row[2])
		buf.WriteString(row[2])
		buf.WriteString("\"}")
		sep = ",\n"
	}
	buf.WriteString("]\n")

	//render.Data(w, r, buf.Bytes())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(buf.Bytes())
	buf.Reset()
}

func (a *API) getLogFeedExecutionsHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	orderid := chi.URLParam(r, "orderid")

	res, err := a.service.ListRange(KeyPrefixLog+id, 0, -1)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	cmp := "[Order:" + orderid + ";Executions:"

	sep := ""
	buf := bytes.NewBufferString("[")
	for i := range res {
		if !strings.Contains(res[i], cmp) {
			continue
		}

		row := strings.SplitN(res[i], ":", 3)
		buf.WriteString(sep)
		buf.WriteString("{\"at\":")
		buf.WriteString(row[0])
		buf.WriteString(",\"msg\":\"")
		row[2] = strings.ReplaceAll(row[2], "\n", "")
		row[2] = strings.ReplaceAll(row[2], "\r", "")
		row[2] = strings.ReplaceAll(row[2], "\t", " ")
		row[2] = strings.ReplaceAll(row[2], "\"", "'")
		row[2] = strings.ReplaceAll(row[2], "\\\"", "'")
		row[2] = strings.TrimSpace(row[2])
		buf.WriteString(row[2])
		buf.WriteString("\"}")
		sep = ",\n"
	}
	buf.WriteString("]\n")

	//render.Data(w, r, buf.Bytes())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(buf.Bytes())
	buf.Reset()
}

func (a *API) getInfoHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	var res map[string]string
	err := a.service.GetObject(KeyPrefixInput+id, &res)
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}

	render.JSON(w, r, res)
}

func (a *API) getPlatformStatusHandler(w http.ResponseWriter, r *http.Request) {

	res, err := a.service.PlatformStatus()
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}

	res.PnL, _ = a.pnlch.LastValue("PNL:BUOYANCY")

	pd, pl := 0.0, 0.0
	a.pnlch.Traverse(func(cu *charter.Curve) {
		if cu.Attr() == "pd" {
			v, _ := cu.LastValue()
			pd += v
		}

		if cu.Attr() == "pl" {
			v, _ := cu.LastValue()
			pl += v
		}
	})

	res.PositionDelta = pd
	res.PositionDeltaLimit = pl

	render.JSON(w, r, res)
}

func (a *API) getPortfolioPnlDebugHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *API) getPnLMultiarrayHandler(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		from int64
	)

	if f := r.URL.Query().Get("from"); f != "" {
		if from, err = strconv.ParseInt(f, 10, 64); err != nil {
			fail(w, http.StatusBadRequest, "parsing parameter from failed: "+err.Error())
			return
		}
	}

	id := chi.URLParam(r, "id")

	cu := a.pnlch.Curve(id)
	if cu == nil {
		render.Render(w, r, ErrRender(errors.New("unknown chart id "+id)))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = charter.WritePositiveNegativeJSONTo(cu, int(from), w)
	if err != nil {
		fail(w, http.StatusInternalServerError, "writing JSON failed")
	}
	return
}

func (a *API) downloadFile(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	//println(id)

	m := map[string]string{}

	if err := a.service.GetObject("I:"+id, &m); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	fname, ok := m["stateModelFileName"]
	if !ok {
		render.Render(w, r, ErrRender(errors.New("I:"+id+" has no attribute stateModelFileName")))
		return
	}
	if fname == "" {
		render.Render(w, r, ErrRender(errors.New("I:"+id+" has empty attribute stateModelFileName")))
		return
	}

	//println("stateModelFileName", fname)
	//fname = "/home/gera/temp/kot.txt"

	fi, err := os.Stat(fname)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	fs, err := a.service.FileStream(fname)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	defer fs.Close()

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	//w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "text/plain")
	//	w.Header().Set("Content-Encoding", "gzip")

	//t := time.Now()
	if _, err := io.Copy(w, fs); err != nil {
		if err != io.EOF {
			render.Render(w, r, ErrRender(err))
			return
		}
	}
	//println("completed", time.Since(t).String())
}

func (a *API) postRestartPlatformHandler(w http.ResponseWriter, r *http.Request) {

	a.log.Info().Msg("platform restart requested")

	if a.restartAllowedBefore == nil {
		fail(w, http.StatusForbidden, "platform restart feature not allowed")
		return
	}

	if time.Now().After(*a.restartAllowedBefore) {
		fail(w, http.StatusForbidden, "allowed deadline is over")
		return
	}

	if a.restartShellCommand == "" {
		fail(w, http.StatusForbidden, "platform restart command not specified")
		return
	}

	go func() {
		var (
			prm []string
			err error
		)

		sidx := strings.Index(a.restartShellCommand, " ")
		if sidx > 0 {
			prm = append(prm, a.restartShellCommand[0:sidx])
			prm = append(prm, a.restartShellCommand[sidx+1:])
		} else {
			prm = append(prm, a.restartShellCommand)
		}

		prm[0], err = exec.LookPath(prm[0])
		if err != nil {
			a.log.Error().Msg("script not found in path")
			return
		}

		cmd := exec.Command(prm[0], prm[1:]...) // Command("/usr/bin/bash", "-c", `echo 'hello world'`) //
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			a.log.Error().Str("errmsg", err.Error()).Msg("script execution failed")
		} else {
			a.log.Info().Str("script", a.restartShellCommand).Msg("script executed succesfully")
		}
	}()

	render.JSON(w, r, Response{true})
}

func (a *API) getHolidaysHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if err := a.h.JSON(w); err != nil {
		fail(w, http.StatusInternalServerError, err.Error())
	}
}
