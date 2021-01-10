package templates

import (
	"github.com/sirupsen/logrus"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	log *logrus.Logger
)

// SetLogger set output log
func SetLogger(l *logrus.Logger) {
	log = l
}

// This function writes the separator s when its called except in the last iteration where we return an empty string
// TODO: Test in more range loop scenarios (1 loop, 3 loops ...)
func separator(s string, iter int) func(maxIter int) string {
	i := 0

	return func(maxIter int) string {
		i++
		c := i + iter*i

		if c == maxIter {
			return ""
		}
		return s
	}
}

// multiply returns the product of a and b.
func multiply(b, a int) int {
	return b * a
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func regexFind(regex string, s string) string {
	r := regexp.MustCompile(regex)
	return r.FindString(s)
}

// Initialize a templates.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom templates functions and the functions themselves.
var functions = template.FuncMap{
	"separator":  separator,
	"humanData":  humanDate,
	"multiply":   multiply,
	"regexFind":  regexFind,
	"trimSuffix": func(a, b string) string { return strings.TrimSuffix(b, a) },
	"trimPrefix": func(a, b string) string { return strings.TrimPrefix(b, a) },
}

// Parse the files when starting the service and store the parsed templates in an in-memory cache
func NewTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.tmpl'. This essentially gives us a slice of all the
	// templates for the application.
	templates, err := filepath.Glob(filepath.Join(dir, "*.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the templates one-by-one
	for _, tpl := range templates {
		// Extract the file name (like '*.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(tpl)

		// The templates.FuncMap must be registered with the templates set before you
		// call the ParseFiles() method. This means we have to use templates.New() to
		// create an empty templates set, use the Funcs() method to register the
		// templates.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(tpl)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any templates to the
		// templates set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	log.Debugf("Loaded Template files on memory cache %v", cache)


	return cache, nil
}
