package r2

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var puts = fmt.Println

type Env struct {
	R    *http.Request
	W    http.ResponseWriter
	Path Path
}

type Handler func(*Env)

const (
	any      = "?"
	regexSep = "!"
)

type trieNode struct {
	// each node may have zero or more children, but at most ONE child can be a parameter.
	children map[string]*trieNode
	// http method to handler
	handlers map[string]Handler
	// the name of the parameter for this node (if any)
	paramName string
	paramRe   *regexp.Regexp
}

type Router struct {
	root    *trieNode
	prefix  string
	regexes map[string]*regexp.Regexp
}

func newNode() *trieNode {
	return &trieNode{children: make(map[string]*trieNode)}
}

func NewRouter(prefix string) *Router {
	return &Router{newNode(), prefix, map[string]*regexp.Regexp{}}
}

func (r *Router) Route(method, path string, handler Handler) {
	r.route(path, handler, method)
}
func (r *Router) Get(path string, handler Handler) {
	r.route(path, handler, "GET")
}
func (r *Router) Post(path string, handler Handler) {
	r.route(path, handler, "POST")
}
func (r *Router) Put(path string, handler Handler) {
	r.route(path, handler, "PUT")
}
func (r *Router) Delete(path string, handler Handler) {
	r.route(path, handler, "DELETE")
}
func (r *Router) Patch(path string, handler Handler) {
	r.route(path, handler, "PATCH")
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// start := time.Now()
	handlerMap, pathVars := r.get(req.URL.Path)

	if handlerMap == nil {
		http.NotFound(w, req)
		return
	}

	handler, found := handlerMap[req.Method]
	if !found {
		handler, found = handlerMap[any]
		if !found {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
	env := &Env{R: req, W: w, Path: pathVars}
	handler(env)
	// puts(time.Now().Sub(start), "\n")
}

func (r *Router) route(urlPath string, handler Handler, method string) {

	failIfEmpty(urlPath, handler)

	node := r.root

	if urlPath == "/" {
		r.add(urlPath, node, handler, method)
		return
	}

	for _, part := range split(urlPath) {

		isParam := strings.HasPrefix(part, ":")
		name, regex := r.separate(part)

		if isParam && !validParam(name, node) {
			log.Fatal(fmt.Sprintf("parameter conflict routing %v with handler %v", urlPath, funcName(handler)))
		}

		key := name
		if isParam {
			key = any
		}
		// if there's already a child node for this part of the path,
		// then use it and descend
		child, found := node.children[key]
		if found {
			node = child
			continue
		}

		child = newNode()
		node.children[key] = child

		if isParam {
			child.paramName = name
			child.paramRe = regex
		}
		node = child
	}

	r.add(urlPath, node, handler, method)
}

func validParam(name string, node *trieNode) bool {

	// if node is a leaf, there's no conflict
	if len(node.children) == 0 {
		return true
	}
	// check to see if a parameter has already been stored
	prevNode, prevParam := node.children[any]
	if !prevParam {
		// no previous param has been stored (only static parts), so no conflict
		return true
	}
	// there was a previous param, so potentially two parameters for this part of path
	// check whether current param name matches previous stored param name.
	// if not the same, there's a conflict
	return name == prevNode.paramName
}

func (r *Router) add(urlPath string, node *trieNode, handler Handler, method string) {

	if node.handlers == nil {
		node.handlers = make(map[string]Handler)
	}

	// check if handler for method already stored
	_, found := node.handlers[method]
	if found {
		log.Fatal(fmt.Sprintf("existing method %v found for path %v", method, urlPath))
	}
	node.handlers[method] = handler
}

func (r *Router) get(path string) (map[string]Handler, Path) {

	if !strings.HasPrefix(path, r.prefix) {
		return nil, nil
	}

	path = path[len(r.prefix):]
	l := len(path)

	if l == 0 || (l == 1 && path[0] == '/') {
		return r.root.handlers, nil
	}
	node := r.root
	var found bool
	var next *trieNode
	var key string
	var vars Path

	// ignore leading and trailing slashes -- could make this an option
	if path[0] == '/' {
		path = path[1:]
		l -= 1
	}
	if path[l-1] == '/' {
		path = path[:l]
	}

	start, end := 0, 0

	for pos, char := range path {

		if char == '/' || pos == l-1 {
			end = pos
			if char != '/' {
				end += 1
			}
			key = path[start:end]

			// check first for a static part
			next, found = node.children[key]
			if !found {
				// not found, so check now if a param is available
				next, found = node.children[any]
				if !found {
					return nil, nil
				}
				// finally, if a regex, check it matches
				if next.paramRe != nil && !next.paramRe.MatchString(key) {
					return nil, nil
				}

				if vars == nil {
					vars = make(Path)
				}
				vars[next.paramName] = key
			}
			start = pos + 1
			node = next
		}
	}
	return node.handlers, vars
}

func (r *Router) separate(s string) (string, *regexp.Regexp) {
	// todo: pass in full path and handler function name to display with error message
	// or return error and show elsewhere

	if !strings.HasPrefix(s, ":") {
		return s, nil
	}

	i := strings.Index(s, regexSep)
	useRegex := i != -1

	// remove initial colon
	name := s[1:]
	if useRegex {
		// if a regex provided, take only up to the regex separator "!"
		name = s[1:i]
	}
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		log.Fatal("parameter must have name: " + s)
	}
	if !useRegex {
		return name, nil
	}
	// add one to skip separator
	pattern := s[i+1:]
	regex, found := r.regexes[pattern]
	if found {
		return name, regex
	}
	regex = compileRe(pattern)
	r.regexes[pattern] = regex
	return name, regex
}

func split(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func compileRe(pattern string) *regexp.Regexp {
	switch pattern {
	case "":
		log.Fatal("no pattern provided")
	case "int":
		pattern = `^-?\d+$`
	case "float":
		pattern = `^-?\d+(?:\.\d+)?$`
	default:
		if !strings.HasPrefix(pattern, "^") {
			pattern = "^" + pattern
		}
		if !strings.HasSuffix(pattern, "$") {
			pattern += "$"
		}
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err.Error())
	}
	return re
}

func failIfEmpty(path string, handler Handler) {
	if handler == nil {
		log.Fatal("nil handler for path: " + path)
	}
	if path == "" {
		log.Fatal("urlPath empty string for handler: " + funcName(handler))
	}
}

func funcName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	dot := strings.LastIndex(name, ".")
	return name[dot+1:]
}

type Path map[string]string

func (p Path) Get(key string) string {
	if p == nil {
		return ""
	}
	val, _ := p[key]
	return val
}

func (p Path) Int(key string) int {
	if p == nil {
		return 0
	}
	str, _ := p[key]
	val, _ := strconv.Atoi(str)
	return val
}

func (p Path) Float(key string) float64 {
	if p == nil {
		return 0
	}
	str, _ := p[key]
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

func (r *Router) Print() {
	printTree("", strings.TrimLeft(r.prefix, "/"), r.root, true)
}

func printTree(prefix, name string, node *trieNode, last bool) {

	s := prefix

	if last {
		s += "└"
	} else {
		s += "├"
	}
	s += "───" + strings.Replace(name, "?", ":", 1) + node.paramName

	for method, h := range node.handlers {
		s += " " + fmt.Sprintf("%v %v", method, funcName(h))
	}
	puts(s)

	for i, part := range sortedParts(node) {

		child := node.children[part]
		lastSib := i == len(node.children)-1
		nextPrefix := prefix

		if last {
			nextPrefix += "    "
		} else {
			nextPrefix += "│    "
		}
		printTree(nextPrefix, part, child, lastSib)
	}
}

func sortedParts(node *trieNode) []string {
	parts := make([]string, len(node.children))
	i := 0
	for part := range node.children {
		parts[i] = part
		i++
	}
	sort.Strings(parts)
	return parts
}
