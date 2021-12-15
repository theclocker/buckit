package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Aliases: []string{
		"s",
	},
	Short: "Serve the app's web interface and REST Api",
	Run: serverCommandHandler,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serverCommandHandler(cmd *cobra.Command, args []string) {
	apiRouteMatch := "{api:[a-zA-Z0-9-/]+}"
	r := mux.NewRouter()
	r.Path(fmt.Sprintf("/a/bulk/%s", apiRouteMatch)).HandlerFunc(DevHandler)
	r.Path(fmt.Sprintf("/a/single/%s", apiRouteMatch)).HandlerFunc(ApiHandler)
	r.Path(fmt.Sprintf("/%s", apiRouteMatch)).HandlerFunc(DevHandler)
	go func(){
		log.Println(http.ListenAndServe(":8080", r))
	}()
	fmt.Println("Serving app on port 8080, See docs for routes and usage")
	select {}
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	apis := map[string]string{
		"ip-api": "http://ip-api.com/",
	}
	url := mux.Vars(r)["api"]
	urlVars := strings.Split(url, "/")
	apiRequested := urlVars[0]

	query := r.URL.Query()
	fmt.Fprintln(w, query.Get("query"))
	fullUrl := apiRequested
	fmt.Fprintln(w, fullUrl)
	fmt.Fprintln(w, r.Method)
	fmt.Fprintf(w, apis[apiRequested])
}

func DevHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Under development")
}
