package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"gopkg.in/macaron.v1"
)

var config Config

func init() {
	fileChoices := []string{"config.toml", "/etc/laminar-webform.toml"}
	file := ""
	for _, f := range fileChoices {
		fpath, _ := filepath.Abs(f)
		exists := fileExists(fpath)
		if exists {
			file = fpath
			break
		}
	}

	log.Println("Parsing config:", file)

	if _, err := toml.DecodeFile(file, &config); err != nil {
		fmt.Printf("Could not parse config: %s\n", err.Error())
		os.Exit(1)
		return
	}
}

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())

	m.Get("/", func(ctx *macaron.Context, log *log.Logger, req *http.Request) {
		log.Println(req.Header.Get("Authorization"))

		ctx.Data["config"] = config.General
		ctx.Data["forms"] = config.Forms
		ctx.HTML(200, "main")
	})

	m.Get("/form/:form", func(ctx *macaron.Context, x csrf.CSRF) {
		ctx.Data["config"] = config.General
		ctx.Data["form"] = config.Forms[ctx.Params("form")]
		ctx.Data["submit"] = "/form/" + ctx.Params("form")
		ctx.Data["csrf_token"] = x.GetToken()
		// ctx.Data["validation_error"] = "something went wrong"
		ctx.HTML(200, "form")
	})

	m.Post("/form/:form", csrf.Validate, func(ctx *macaron.Context, x csrf.CSRF, log *log.Logger) {

		var renderFormWithError = func(errortext string) {
			ctx.Data["form"] = config.Forms[ctx.Params("form")]
			ctx.Data["submit"] = "/form/" + ctx.Params("form")
			ctx.Data["csrf_token"] = x.GetToken()
			ctx.Data["validation_error"] = errortext
			ctx.HTML(200, "form")
		}

		args := []string{"queue", config.Forms[ctx.Params("form")].Job}

		form := config.Forms[ctx.Params("form")]
		for _, f := range form.Fields {
			val := ctx.Req.FormValue(f.Name)
			if val == "" {
				renderFormWithError("required value not submitted")
				return
			}

			if f.Type == "select" && !sliceContains(val, f.Options) {
				renderFormWithError("selected value invalid")
				return
			}

			if f.Type == "text" {
				var validText = regexp.MustCompile(f.Filter)
				if !validText.MatchString(val) {
					renderFormWithError("text value contains one or more disallowed characters")
					return
				}
			}

			args = append(args, f.Name+"="+val)
		}

		// Find laminarc
		path, err := exec.LookPath("laminarc")
		if err != nil {
			renderFormWithError("laminarc is not installed, so this job cannot be queued. please contact an administrator.")
			return
		}

		if config.General.Debug {
			output := make([]interface{}, len(args)+1)
			output[0] = "laminarc"
			for i, arg := range args {
				output[i+1] = arg
			}
			log.Println("running laminarc command:")
			log.Println(output...)
		}

		cmd := exec.Command(path, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()

		if config.General.Debug {
			log.Println("laminarc output: ", out.String())
		}

		if err != nil {
			renderFormWithError("job was not queued successfully. error was: " + err.Error())
			return
		}

		ctx.Redirect(fmt.Sprintf("%s/jobs/%s", config.General.LaminarURL, config.Forms[ctx.Params("form")].Job))
	})

	m.Run()
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		log.Println(err.Error())
		return false
	}
	return true
}

func sliceContains(needle string, haystack []string) bool {
	for _, b := range haystack {
		if b == needle {
			return true
		}
	}
	return false
}

// @todo get this to return a username if there's anything obvious available.
// also consider using jwt.sub from caddy's http.login cookie.
func getReason() string {
	return "Manually executed via laminar-webform"
}
