package main

import (
	"net/http"
	"strings"
)

type Routeable interface {
	Route(method string, pattern string, handleFunc handlerFunc)
}

type Handler interface {
	ServeHTTP(c *Context)
	Routeable //负责路由
}

type HandlerBasedOnTree struct {
	root *node
}

type node struct {
	path     string
	children []*node

	// 如果这是叶子节点
	// 那么匹配上之后就可以调用该方法
	handler handlerFunc
}

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	handler, found := h.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not found!"))
	}
	handler(c)
}

func (h *HandlerBasedOnTree) findRouter(path string) (handlerFunc, bool) {
	// 去除头尾可能有的/，然后按照/切割成段
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, p := range paths {
		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMatchChild(cur, p)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	// 到这里，应该是找完了
	if cur.handler == nil {
		// 到达这里是因为这种场景
		// 比如说你注册了 /user/profile
		// 然后你访问 /user
		return nil, false
	}
	return cur.handler, true
}

func (h *HandlerBasedOnTree) Route(method string, pattern string, handleFunc handlerFunc) {
	pattern = strings.Trim(pattern, "/") //去掉前后“/”
	paths := strings.Split(pattern, "/") //按照/切开
	cur := h.root
	for index, path := range paths {
		mathChild, ok := h.findMatchChild(cur, path)
		if ok {
			cur = mathChild
		} else {
			h.createSubTree(cur, paths[index:], handleFunc)
			return
		}
	}
	// 离开了循环，说明我们加入的是短路径，
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handleFunc
}

func (h *HandlerBasedOnTree) findMatchChild(parent *node, path string) (*node, bool) {
	for _, child := range parent.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handlerFn handlerFunc) {
	//paths可能是friends/xiaoming/address
	cur := root
	for _, path := range paths {
		nn := newNode(path)
		//user.children = [profile, home, friends]
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFn
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 4),
	}
}
func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: &node{},
	}
}
