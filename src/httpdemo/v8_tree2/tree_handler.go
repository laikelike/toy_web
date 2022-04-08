package main

import (
	"errors"
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root *node
}

var supportMethods = [4]string{http.MethodPost, http.MethodGet,
	http.MethodDelete, http.MethodPut}

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
		return
	}
	handler(c)
}

// findRouter 查找完全匹配的URL
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

func (h *HandlerBasedOnTree) Route(method string, pattern string, handleFunc handlerFunc) error {
	//做校验
	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}
	pattern = strings.Trim(pattern, "/") //去掉前后“/”
	paths := strings.Split(pattern, "/") //按照/切开
	cur := h.root
	for index, path := range paths {
		mathChild, ok := h.findMatchChild(cur, path)
		if ok {
			cur = mathChild
		} else {
			h.createSubTree(cur, paths[index:], handleFunc)
			return nil
		}
	}
	// 离开了循环，说明我们加入的是短路径，
	// 比如说我们先加入了 /order/detail
	// 再加入/order，那么会走到这里
	cur.handler = handleFunc
	return nil
}

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")

func (h *HandlerBasedOnTree) validatePattern(pattern string) error {
	// 校验 *，如果存在，必须在最后一个，并且它前面必须是/
	// 即我们只接受 /* 的存在，abc*这种是非法

	pos := strings.Index(pattern, "*")
	// 找到了 *
	if pos > 0 {
		// 必须是最后一个
		if pos != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}
		if pattern[pos-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
}

// findMatchChild 添加通配符匹配
func (h *HandlerBasedOnTree) findMatchChild(parent *node, path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range parent.children {

		if child.path == path && child.path != "*" {
			return child, true
		}
		// 命中了通配符
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
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
	root := &node{}
	return &HandlerBasedOnTree{
		root: root,
	}
}
