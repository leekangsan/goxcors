package server

import (
	"appengine"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"appengine/urlfetch"
	"io/ioutil"
	"strings"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/post", getCrossDomainRequest)
	http.Handle("/", r)
}

/* curl -s -A "Mozilla/5.0 (X11; Linux x86_64; rv:5.0) Gecko/20100101 Firefox/5.0" "http://translate.google.com/translate_a/t?client=p&sl=&tl=ko&text=compensation" */
func getCrossDomainRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	reqURL := strings.Split(r.URL.RawQuery,"/post?")[0]
	c.Infof("Reqeusted URL: %q", reqURL)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	w.Header().Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:5.0) Gecko/20100101 Firefox/5.0)")

	req, err := http.NewRequest("POST", reqURL, nil)
	req.Header.Add("User-Agent", "myClient")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Errorf("ioutil error Get %s", err)
		fmt.Fprintf(w, "{'err':'%s'", err)
		return
	}
	fmt.Fprintf(w, "%s", body)
}