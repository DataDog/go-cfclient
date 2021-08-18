package cfclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListApps(t *testing.T) {
	Convey("List Apps", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps", []string{listAppsPayload}, "Test-golang", 200, "inline-relations-depth=2", nil},
			{"GET", "/v2/appsPage2", []string{listAppsPayloadPage2}, "Test-golang", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		apps, err := client.ListApps()
		assertAppList(apps, err)
	})
}

func TestListAppsByRoute(t *testing.T) {
	Convey("List Apps by Route", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/routes/a3fb8d86-4620-4725-b643-0b5122a432e0/apps", []string{listAppsPayload}, "Test-golang", 200, "", nil},
			{"GET", "/v2/appsPage2", []string{listAppsPayloadPage2}, "Test-golang", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		apps, err := client.ListAppsByRoute("a3fb8d86-4620-4725-b643-0b5122a432e0")
		assertAppList(apps, err)
	})
}

func assertAppList(apps []App, err error) {
	So(err, ShouldBeNil)

	So(len(apps), ShouldEqual, 2)
	So(apps[0].Guid, ShouldEqual, "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c")
	So(apps[0].CreatedAt, ShouldEqual, "2014-10-10T21:03:13+00:00")
	So(apps[0].UpdatedAt, ShouldEqual, "2014-11-10T14:07:31+00:00")
	So(apps[0].Name, ShouldEqual, "app-test")
	So(apps[0].Memory, ShouldEqual, 256)
	So(apps[0].Instances, ShouldEqual, 1)
	So(apps[0].DiskQuota, ShouldEqual, 1024)
	So(apps[0].SpaceGuid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
	So(apps[0].StackGuid, ShouldEqual, "2c531037-68a2-4e2c-a9e0-71f9d0abf0d4")
	So(apps[0].State, ShouldEqual, "STARTED")
	So(apps[0].Command, ShouldEqual, "")
	So(apps[0].Buildpack, ShouldEqual, "https://github.com/cloudfoundry/buildpack-go.git")
	So(apps[0].DetectedBuildpack, ShouldEqual, "")
	So(apps[0].DetectedBuildpackGuid, ShouldEqual, "0d22f6a1-76c5-417f-ac6c-d9d21463ecbc")
	So(apps[0].HealthCheckHttpEndpoint, ShouldEqual, "")
	So(apps[0].HealthCheckType, ShouldEqual, "port")
	So(apps[0].HealthCheckTimeout, ShouldEqual, 0)
	So(apps[0].Diego, ShouldEqual, true)
	So(apps[0].EnableSSH, ShouldEqual, true)
	So(apps[0].DetectedStartCommand, ShouldEqual, "app-launching-service-broker")
	So(apps[0].DockerImage, ShouldEqual, "")
	So(apps[0].DockerCredentialsJSON["redacted_message"], ShouldEqual, "[PRIVATE DATA HIDDEN]")
	So(apps[0].Environment["FOOBAR"], ShouldEqual, "QUX")
	So(apps[0].StagingFailedReason, ShouldEqual, "")
	So(apps[0].StagingFailedDescription, ShouldEqual, "")
	So(apps[0].PackageState, ShouldEqual, "PENDING")
	So(len(apps[0].Ports), ShouldEqual, 1)
	So(apps[0].Ports[0], ShouldEqual, 8080)

	So(apps[1].Guid, ShouldEqual, "f9ad202b-76dd-44ec-b7c2-fd2417a561e8")
	So(apps[1].Name, ShouldEqual, "app-test2")
}

func TestGetAppByGuidNoInlineCall(t *testing.T) {
	Convey("App By GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", []string{appPayload}, "Test-golang", 200, "", nil},
			{"GET", "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73", []string{spacePayload}, "Test-golang", 200, "", nil},
			{"GET", "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2", []string{orgPayload}, "Test-golang", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}

		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.GetAppByGuidNoInlineCall("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
	})

	Convey("App By GUID and space returns error", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", []string{appPayload}, "Test-golang", 200, "", nil},
			{"GET", "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73", []string{"error"}, "Test-golang", 500, "", nil},
			{"GET", "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2", []string{"error"}, "Test-golang", 500, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}

		client, err := NewClient(c)
		So(err, ShouldBeNil)

		_, err = client.GetAppByGuidNoInlineCall("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldNotBeNil)

		So(err.Error(), ShouldStartWith, "Unable to get the Space for the app test-env")
	})

	Convey("App By GUID with environment variables with different types", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", []string{appPayloadWithEnvironment}, "Test-golang", 200, "", nil},
			{"GET", "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73", []string{spacePayload}, "Test-golang", 200, "", nil},
			{"GET", "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2", []string{orgPayload}, "Test-golang", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
			UserAgent:  "Test-golang",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.GetAppByGuidNoInlineCall("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)
		So(app.Environment["string"], ShouldEqual, "string")
		So(app.Environment["int"], ShouldEqual, 1)
	})
}

func TestAppByGuid(t *testing.T) {
	Convey("App By GUID", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", []string{appPayload}, "", 200, "inline-relations-depth=2", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.GetAppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
	})

	Convey("App By GUID with environment variables with different types", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2", []string{appPayloadWithEnvironment}, "", 200, "inline-relations-depth=2", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app, err := client.GetAppByGuid("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(app.Environment["string"], ShouldEqual, "string")
		So(app.Environment["int"], ShouldEqual, 1)
	})
}

func TestGetAppInstances(t *testing.T) {
	Convey("App completely running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances", []string{appInstancePayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appInstances, err := client.GetAppInstances("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appInstances["0"].State, ShouldEqual, "RUNNING")
		So(appInstances["1"].State, ShouldEqual, "RUNNING")

		var d0 float64 = 1455210430.5104606
		var d1 float64 = 1455210430.3912115
		date0 := time.Unix(int64(d0), 0)
		date1 := time.Unix(int64(d1), 0)

		So(appInstances["0"].Since.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appInstances["1"].Since.Format(time.UnixDate), ShouldEqual, date1.Format(time.UnixDate))
		So(appInstances["0"].Since.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appInstances["1"].Since.ToTime(), ShouldHaveSameTypeAs, date1)

	})

	Convey("App partially running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances", []string{appInstanceUnhealthyPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appInstances, err := client.GetAppInstances("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appInstances["0"].State, ShouldEqual, "RUNNING")
		So(appInstances["1"].State, ShouldEqual, "STARTING")

		var d0 float64 = 1455210430.5104606
		var d1 float64 = 1455210430.3912115
		date0 := time.Unix(int64(d0), 0)
		date1 := time.Unix(int64(d1), 0)

		So(appInstances["0"].Since.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appInstances["1"].Since.Format(time.UnixDate), ShouldEqual, date1.Format(time.UnixDate))
		So(appInstances["0"].Since.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appInstances["1"].Since.ToTime(), ShouldHaveSameTypeAs, date1)

	})
}

func TestGetAppStats(t *testing.T) {
	Convey("App stats completely running", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/stats", []string{appStatsPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appStats, err := client.GetAppStats("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(appStats["0"].State, ShouldEqual, "RUNNING")
		So(appStats["1"].State, ShouldEqual, "RUNNING")
		So(appStats["2"].State, ShouldEqual, "RUNNING")
		So(appStats["3"].State, ShouldEqual, "RUNNING")

		date0, _ := time.Parse("2006-01-02 15:04:05 -0700", "2016-09-17 15:46:17 +0000")
		date1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2016-09-17 15:46:17 +0000")
		date2, _ := time.Parse(time.RFC3339Nano, "2017-04-06T20:32:19.273294439Z")
		date3, _ := time.Parse("2006-01-02 15:04:05 MST", "2017-04-12 15:27:44 UTC")

		So(appStats["0"].Stats.Usage.Time.Format(time.UnixDate), ShouldEqual, date0.Format(time.UnixDate))
		So(appStats["1"].Stats.Usage.Time.Format(time.RFC3339), ShouldEqual, date1.Format(time.RFC3339))
		So(appStats["2"].Stats.Usage.Time.Format(time.RFC3339Nano), ShouldEqual, date2.Format(time.RFC3339Nano))
		So(appStats["3"].Stats.Usage.Time.Format("2006-01-02 15:04:05 MST"), ShouldEqual, date3.Format("2006-01-02 15:04:05 MST"))
		So(appStats["0"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date0)
		So(appStats["1"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date1)
		So(appStats["2"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date2)
		So(appStats["3"].Stats.Usage.Time.ToTime(), ShouldHaveSameTypeAs, date3)
		So(appStats["0"].Stats.Usage.CPU, ShouldEqual, 0.36580239597146486)
		So(appStats["1"].Stats.Usage.CPU, ShouldEqual, 0.33857742931636664)
		So(appStats["2"].Stats.Usage.CPU, ShouldEqual, 0.33857742931636664)
		So(appStats["3"].Stats.Usage.CPU, ShouldEqual, 0.33857742931636664)
		So(appStats["0"].Stats.Usage.Mem, ShouldEqual, 518123520)
		So(appStats["1"].Stats.Usage.Mem, ShouldEqual, 530731008)
		So(appStats["2"].Stats.Usage.Mem, ShouldEqual, 530731008)
		So(appStats["3"].Stats.Usage.Mem, ShouldEqual, 530731008)
		So(appStats["0"].Stats.Usage.Disk, ShouldEqual, 151150592)
		So(appStats["1"].Stats.Usage.Disk, ShouldEqual, 151150592)
		So(appStats["2"].Stats.Usage.Disk, ShouldEqual, 151150592)
		So(appStats["3"].Stats.Usage.Disk, ShouldEqual, 151150592)

	})
}

func TestGetAppRoutes(t *testing.T) {
	Convey("List app routes", t, func() {
		setup(MockRoute{"GET", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/routes", []string{appRoutesPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		routes, err := client.GetAppRoutes("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)

		So(len(routes), ShouldEqual, 1)
		So(routes[0].Guid, ShouldEqual, "311d34d1-c045-4853-845f-05132377ad7d")
		So(routes[0].Host, ShouldEqual, "host-36")
		So(routes[0].Path, ShouldEqual, "/foo")
		So(routes[0].DomainGuid, ShouldEqual, "40a499f7-198a-4289-9aa2-605ba43f92ee")
		So(routes[0].SpaceGuid, ShouldEqual, "c7c0dd06-b078-43d7-adcb-3974cd785fdd")
		So(routes[0].ServiceInstanceGuid, ShouldEqual, "")
		So(routes[0].Port, ShouldEqual, 0)

	})
}

func TestUploadAppBits(t *testing.T) {
	Convey("Upload app bits", t, func() {
		expectedPayload := "this should really be zipped binary data"
		mr := MockRoute{
			Method:   "PUT-FILE",
			Endpoint: "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/bits",
			Status:   201,
			PostForm: &expectedPayload,
			Output:   []string{""},
		}
		setup(mr, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		bits := bytes.NewBufferString(expectedPayload)
		err = client.UploadAppBits(bits, "9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)
	})
}

func TestGetAppBits(t *testing.T) {
	Convey("Get app bits", t, func() {

		next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/gzip")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("apptarbinarydata"))
		})
		s := httptest.NewServer(next)
		defer s.Close()

		mr := MockRouteWithRedirect{
			MockRoute: MockRoute{
				Method:   "GET",
				Endpoint: "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/download",
				Status:   302,
				Output:   []string{""},
			},
			RedirectLocation: s.URL,
		}
		setupWithRedirect(mr, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		bits, err := client.GetAppBits("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)
		So(bits, ShouldNotBeNil)

		var download string
		if b, err := ioutil.ReadAll(bits); err == nil {
			download = string(b)
		}
		So(download, ShouldEqual, "apptarbinarydata")
	})
}

func TestUploadDropletBits(t *testing.T) {
	Convey("Upload droplet bits", t, func() {
		expectedPayload := "this should really be tar'd and gzipped binary data"
		mr := MockRoute{
			Method:   "PUT-FILE",
			Endpoint: "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/droplet/upload",
			Status:   201,
			PostForm: &expectedPayload,
			Output:   []string{`{"metadata":{"guid": "123", "url":"abc"}}`},
		}
		setup(mr, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		bits := bytes.NewBufferString(expectedPayload)
		_, err = client.UploadDropletBits(bits, "9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)
	})
}

func TestGetDropletBits(t *testing.T) {
	Convey("Get droplet bits", t, func() {

		next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/gzip")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("apptarbinarydata"))
		})
		s := httptest.NewServer(next)
		defer s.Close()

		mr := MockRouteWithRedirect{
			MockRoute: MockRoute{
				Method:   "GET",
				Endpoint: "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/droplet/download",
				Status:   302,
				Output:   []string{""},
			},
			RedirectLocation: s.URL,
		}
		setupWithRedirect(mr, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		bits, err := client.GetDropletBits("9902530c-c634-4864-a189-71d763cb12e2")
		So(err, ShouldBeNil)
		So(bits, ShouldNotBeNil)

		var download string
		if b, err := ioutil.ReadAll(bits); err == nil {
			download = string(b)
		}
		So(download, ShouldEqual, "apptarbinarydata")
	})
}

func TestKillAppInstance(t *testing.T) {
	Convey("Kills an app instance", t, func() {
		setup(MockRoute{"DELETE", "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/instances/0", []string{""}, "", 204, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		So(client.KillAppInstance("9902530c-c634-4864-a189-71d763cb12e2", "0"), ShouldBeNil)
	})
}

func TestAppEnv(t *testing.T) {
	Convey("Find app space", t, func() {
		setup(MockRoute{"GET", "/v2/apps/a7c47787-a982-467c-95d7-9ab17cbcc918/env", []string{appEnvPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		appEnv, err := client.GetAppEnv("a7c47787-a982-467c-95d7-9ab17cbcc918")
		So(err, ShouldBeNil)
		So(appEnv.StagingEnv, ShouldResemble, map[string]interface{}{"STAGING_ENV": "staging_value"})
		So(appEnv.RunningEnv, ShouldResemble, map[string]interface{}{"RUNNING_ENV": "running_value"})
		So(appEnv.Environment, ShouldResemble, map[string]interface{}{"env_var": "env_val"})
		So(appEnv.SystemEnv["VCAP_SERVICES"].(map[string]interface{})["abc"], ShouldEqual, 123)
		So(appEnv.ApplicationEnv["VCAP_APPLICATION"].(map[string]interface{})["application_name"], ShouldEqual, "name-2245")
	})
}

func TestAppSummary(t *testing.T) {
	Convey("Get app summary", t, func() {
		setup(MockRoute{"GET", "/v2/apps/b5f0d1bd-a3a9-40a4-af1a-312ad26e5379/summary", []string{appSummaryPayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app := &App{
			Guid: "b5f0d1bd-a3a9-40a4-af1a-312ad26e5379",
			c:    client,
		}

		summary, err := app.Summary()
		So(err, ShouldBeNil)

		So(summary.Guid, ShouldEqual, "b5f0d1bd-a3a9-40a4-af1a-312ad26e5379")
		So(summary.Name, ShouldEqual, "test-app")
		So(summary.ServiceCount, ShouldEqual, 1)
		So(summary.RunningInstances, ShouldEqual, 1)
		So(summary.SpaceGuid, ShouldEqual, "494d8b64-8181-4183-a6d3-6279db8fec6e")
		So(summary.StackGuid, ShouldEqual, "67e019a3-322a-407a-96e0-178e95bd0e55")
		So(summary.Buildpack, ShouldEqual, "ruby_buildpack")
		So(summary.DetectedBuildpack, ShouldEqual, "")
		So(summary.Memory, ShouldEqual, 256)
		So(summary.Instances, ShouldEqual, 1)
		So(summary.DiskQuota, ShouldEqual, 512)
		So(summary.State, ShouldEqual, "STARTED")
		So(summary.Command, ShouldEqual, "")
		So(summary.PackageState, ShouldEqual, "STAGED")
		So(summary.HealthCheckType, ShouldEqual, "port")
		So(summary.HealthCheckTimeout, ShouldEqual, 0)
		So(summary.StagingFailedReason, ShouldEqual, "")
		So(summary.StagingFailedDescription, ShouldEqual, "")
		So(summary.Diego, ShouldEqual, true)
		So(summary.DockerImage, ShouldEqual, "")
		So(summary.DetectedStartCommand, ShouldEqual, "rackup -p $PORT")
		So(summary.EnableSSH, ShouldEqual, true)
		So(summary.DockerCredentials["redacted_message"], ShouldEqual, "[PRIVATE DATA HIDDEN]")
	})
}

func TestAppSpace(t *testing.T) {
	Convey("Find app space", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/foobar", []string{spacePayload}, "", 200, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		app := &App{
			Guid:     "123",
			Name:     "test app",
			SpaceURL: "/v2/spaces/foobar",
			c:        client,
		}
		space, err := app.Space()
		So(err, ShouldBeNil)

		So(space.Name, ShouldEqual, "test-space")
		So(space.Guid, ShouldEqual, "a72fa1e8-c694-47b3-85f2-55f61fd00d73")
	})
}

func TestDeleteApps(t *testing.T) {
	Convey("Delete app", t, func() {
		setup(MockRoute{"DELETE", "/v2/apps/a537761f-9d93-4b30-af17-3d73dbca181b", []string{""}, "", 204, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteApp("a537761f-9d93-4b30-af17-3d73dbca181b")
		So(err, ShouldBeNil)
	})
}

func TestCreateApp(t *testing.T) {
	Convey("Delete app", t, func() {
		setup(MockRoute{"POST", "/v2/apps", []string{appPayload}, "", 201, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		req := AppCreateRequest{
			Name:      "test-env",
			SpaceGuid: "a72fa1e8-c694-47b3-85f2-55f61fd00d73",
		}
		app, err := client.CreateApp(req)
		So(err, ShouldBeNil)
		So(app.Guid, ShouldEqual, "9902530c-c634-4864-a189-71d763cb12e2")
		So(app.Name, ShouldEqual, "test-env")
		So(app.SpaceGuid, ShouldEqual, "a72fa1e8-c694-47b3-85f2-55f61fd00d73")
	})
}

func TestStartApp(t *testing.T) {
	Convey("Start app", t, func() {
		expectedBody := `{ "state": "STARTED" }`
		setup(MockRoute{"PUT", "/v2/apps/a537761f-9d93-4b30-af17-3d73dbca181b", []string{appPayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		So(client.StartApp("a537761f-9d93-4b30-af17-3d73dbca181b"), ShouldBeNil)
	})
}

func TestStopApp(t *testing.T) {
	Convey("Stop app", t, func() {
		expectedBody := `{ "state": "STOPPED" }`
		setup(MockRoute{"PUT", "/v2/apps/a537761f-9d93-4b30-af17-3d73dbca181b", []string{appPayload}, "", http.StatusCreated, "", &expectedBody}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		So(client.StopApp("a537761f-9d93-4b30-af17-3d73dbca181b"), ShouldBeNil)
	})
}
