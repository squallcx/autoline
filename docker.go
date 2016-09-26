package main

import (
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var client *docker.Client

func init() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ = docker.NewClient(endpoint)
	version, _ := client.Version()
	fmt.Println("Version: ", version.Get("Version"))
	fmt.Println("API Version: ", version.Get("ApiVersion"))
}

func check(msg string, err error) bool {
	if err != nil {
		log.Printf("============find %s : %s===========/n", msg, err)
		return false
	}
	return true
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", setHandler)
	r.HandleFunc("/view", viewHandler)
	r.HandleFunc("/check", checkHandler)
	r.HandleFunc("/kill", killHandler)
	http.Handle("/", r)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Panicln(err)
	}
}

type userData struct {
	Username  string
	Password  string
	DealySec  string
	Context   string
	Container string
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	rco := docker.RemoveContainerOptions{}
	rco.ID = "line"
	// rco.RemoveVolumes = true
	rco.Force = true
	err = client.RemoveContainer(rco)
	if err != nil {
		log.Println("not find container or if stop the way is error")
		// log.Panicln(err)
	}
}
func checkHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://127.0.0.1:8787/vnc_viewonly.html?password=1234&autoconnect=true"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	rsp, _ := client.Get(url)
	// fmt.Println(rsp.Status)
	if rsp.StatusCode == 200 {
		fmt.Fprintln(w, "ok")
		return
	}

}
func setHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = r.ParseForm()
	html := `
<html  style="text-align:center">

<body>
username: <input type="text" id="username"><br>
password: <input type="text" id="password"><br>
dealy: <input type="text" id="dealysec"><br>
context: <textarea rows="4" id="context" cols="50"></textarea><br>
 <button onclick="submit()">submit</button><br><br>
 <button onclick="reloadIFrame()" value="reloadview">reloadview</button><br>
 <button onclick="killIFrame()" value="reloadview">stop</button><br>


<script src="https://code.jquery.com/jquery-1.9.1.min.js"></script>
<script>

function killIFrame() {
$.get( "/kill");
}
function submit() {
$.get( "/view", {
      username: $("#username").val(),
      password: $("#password").val(),
      dealysec: $("#dealysec").val(),
      context: $("#context").val(),
 } );

  window.setTimeout("reloadIFrame();", 5000);
}


function reloadIFrame() {

var jqxhr = $.get( "/check", function() {
})
  .done(function() {
  document.getElementById("line").src="http://127.0.0.1:8787/vnc_viewonly.html?password=1234&autoconnect=true";
  })
  .fail(function() {
  window.setTimeout("reloadIFrame();", 1000);
  });

}
$( document ).ready(function() {

  window.setTimeout("reloadIFrame();", 1000);
});
</script>
<iframe id="line" name="line"  width="1024" height="768"  scrolling="no" style="width: 1024px; height: 768px;" src=""></iframe>
</body>
</html>
`
	fmt.Fprintln(w, html)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = r.ParseForm()

	u := userData{
		Username: r.FormValue("username"),

		DealySec:  r.FormValue("dealysec"),
		Context:   r.FormValue("context"),
		Container: fmt.Sprintf("line%s", ""),
	}
	getdata(u)
	html := `
<html  style="text-align:center">
 <button onclick="reloadIFrame()" value="reloadview"></button><br>
 <button onclick="killIFrame()" value="stop"></button><br>

<script src="https://code.jquery.com/jquery-1.9.1.min.js"></script>
<script>
function reloadIFrame() {
  document.getElementById("line").src="http://127.0.0.1:8787/vnc_viewonly.html?password=1234&autoconnect=true";
}
$( document ).ready(function() {
  window.setTimeout("reloadIFrame();", 5000);
});
</script>
<iframe id="line" name="line"  width="1024" height="768"  scrolling="no" style="width: 1024px; height: 768px;" src=""></iframe>
</html>
`
	fmt.Fprintln(w, html)

}
func getdata(user userData) string {
	var err error
	var buf bytes.Buffer
	log.Println(user)
	create(user)
	err = client.Logs(
		docker.LogsOptions{
			Container:    user.Container,
			OutputStream: &buf,
			ErrorStream:  &buf,
			Stderr:       true,
			Stdout:       true,
		},
	)
	if err != nil {
		log.Panicln(err)
	}
	return ""
}
func remove(user userData) {
	var err error
	rco := docker.RemoveContainerOptions{}
	rco.ID = user.Container
	// rco.RemoveVolumes = true
	rco.Force = true
	err = client.RemoveContainer(rco)
	if err != nil {
		log.Println("not find container or if stop the way is error")
		// log.Panicln(err)
	}
}
func create(user userData) {
	var err error
	remove(user)

	dockerContainer, err := client.CreateContainer(
		docker.CreateContainerOptions{
			Config: &docker.Config{
				Image: "squallcx/buildline",
				ExposedPorts: map[docker.Port]struct{}{
					"8787/tcp": {}},
				Env: []string{
					fmt.Sprintf("username=%s", user.Username),
					fmt.Sprintf("password=%s", user.Password),
					fmt.Sprintf("dealysec=%s", user.DealySec),
					fmt.Sprintf("context=%s", user.Context),
				},
			},
			Name: user.Container,
			HostConfig: &docker.HostConfig{
				PortBindings: map[docker.Port][]docker.PortBinding{
					"8787/tcp": []docker.PortBinding{
						docker.PortBinding{
							HostPort: "8787",
						},
					}},
				PublishAllPorts: true,
				Privileged:      false,
			},
		},
	)
	if err != nil {
		log.Panicln(err)
	}

	err = client.StartContainer(dockerContainer.ID, nil)
	if err != nil {
		log.Panicln(err)
	}

}
