package apollo

import (
	"context"
	"fmt"
	"github.com/gitcloneese/agollo"
	"os"
	"strings"
	"sync"
	"time"
	"xy3-proto/pkg/conf/paladin"
	"xy3-proto/pkg/log"
)

// APOLLO Env
const (
	_ApolloSwitch   = "APOLLO"          //apollo开关
	_ApolloAppId    = "APOLLO_APP_ID"   //apollo appId
	_ApolloCluster  = "APOLLO_CLUSTER"  // apollo cluster
	_ApolloEndpoint = "APOLLO_ENDPOINT" //apollo endpoint
	_ApolloSecret   = "APOLLO_SECRET"   //apollo secret
)

type _Apollo struct {
	Switch   bool   //开关
	Env      string //环境
	AppId    string //AppId
	Cluster  string //Cluster
	Endpoint string //地址
	Secret   string //密钥
}

func newApolloEnv() *_Apollo {
	on := false
	apolloSwitch := os.Getenv(_ApolloSwitch)
	if strings.Compare(strings.ToLower(apolloSwitch), "on") == 0 {
		on = true
	}
	return &_Apollo{
		Switch:   on,
		AppId:    os.Getenv(_ApolloAppId),
		Cluster:  os.Getenv(_ApolloCluster),
		Endpoint: os.Getenv(_ApolloEndpoint),
		Secret:   os.Getenv(_ApolloSecret),
	}
}

// 基础db配置
//
//nolint:all
var argsNS = []string{
	"db.txt",
	"etcd.txt",
	"grpc.txt",
	"http.txt",
	"imsdk.txt",
	"kafka.txt",
	"memcache.txt",
	"mongodb.txt",
	"redis_lock.txt",
	"redis.txt",
	"ws.txt",
}

const (
	HttpNS     = "http.txt"
	GrpcNS     = "grpc.txt"
	DbNS       = "db.txt"
	EtcdNS     = "etcd.txt"
	IMSdkNS    = "imsdk.txt"
	KafkaNS    = "kafka.txt"
	MemcacheNS = "memcache.txt"
	MongodbNS  = "mongodb.txt"
	RedisNS    = "redis.txt"
	WsNS       = "ws.txt"
)

var (
	defaultClient  *Client
	defaultClient1 *Client
	defaultClient2 *Client
	defaultClient3 *Client
	defaultClient4 *Client
	defaultClient5 *Client
	defaultClient6 *Client
	defaultClient7 *Client
	defaultClient8 *Client
	defaultClient9 *Client
	defaultNS      = []string{
		HttpNS,
		GrpcNS,
	}
	lock sync.Mutex
)

// 多开几个client 增加初始化速度
func assignClient(c *Client) {
	lock.Lock()
	defer lock.Unlock()
	switch {
	case defaultClient == nil:
		defaultClient = c
	case defaultClient1 == nil:
		defaultClient1 = c
	case defaultClient2 == nil:
		defaultClient2 = c
	case defaultClient3 == nil:
		defaultClient3 = c
	case defaultClient4 == nil:
		defaultClient4 = c
	case defaultClient5 == nil:
		defaultClient5 = c
	case defaultClient6 == nil:
		defaultClient6 = c
	case defaultClient7 == nil:
		defaultClient7 = c
	case defaultClient8 == nil:
		defaultClient8 = c
	case defaultClient9 == nil:
		defaultClient9 = c
	}
}

// AddNs
// 初始化 apollo 基础的配置namespace
func AddNs(ns ...string) {
	defaultNS = append(defaultNS, ns...)
}

// Client
// 封装apolloClient 的实现
type Client struct {
	*_Apollo
	configMaps      map[string]paladin.Setter
	priorityConfigs []string
	values          *paladin.Map
	rawVal          map[string]*paladin.Value
	mx              sync.Mutex
	paladin.Client
	aClient agollo.Client
}

func (c *Client) WatchEvent(context.Context, ...string) <-chan paladin.Event {
	return nil
}

func (c *Client) Close() error {
	return nil
}

// Get 获取ns
func (c *Client) Get(key string) *paladin.Value {
	return c.values.Get(key)
}

// Get
// 获取配置
func Get(Key string) *paladin.Value {
	switch {
	case get(Key) != nil:
		return get(Key)
	case get1(Key) != nil:
		return get1(Key)
	case get2(Key) != nil:
		return get2(Key)
	case get3(Key) != nil:
		return get3(Key)
	case get4(Key) != nil:
		return get4(Key)
	case get5(Key) != nil:
		return get5(Key)
	case get6(Key) != nil:
		return get6(Key)
	case get7(Key) != nil:
		return get7(Key)
	case get8(Key) != nil:
		return get8(Key)
	case get9(Key) != nil:
		return get9(Key)
	}
	return nil
}

func get(Key string) *paladin.Value {
	if defaultClient == nil {
		return nil
	}
	return defaultClient.Get(Key)
}
func get1(Key string) *paladin.Value {
	if defaultClient1 == nil {
		return nil
	}
	return defaultClient1.Get(Key)
}
func get2(Key string) *paladin.Value {
	if defaultClient2 == nil {
		return nil
	}
	return defaultClient2.Get(Key)
}
func get3(Key string) *paladin.Value {
	if defaultClient3 == nil {
		return nil
	}
	return defaultClient3.Get(Key)
}
func get4(Key string) *paladin.Value {
	if defaultClient4 == nil {
		return nil
	}
	return defaultClient4.Get(Key)
}
func get5(Key string) *paladin.Value {
	if defaultClient5 == nil {
		return nil
	}
	return defaultClient5.Get(Key)
}
func get6(Key string) *paladin.Value {
	if defaultClient6 == nil {
		return nil
	}
	return defaultClient6.Get(Key)
}
func get7(Key string) *paladin.Value {
	if defaultClient7 == nil {
		return nil
	}
	return defaultClient7.Get(Key)
}
func get8(Key string) *paladin.Value {
	if defaultClient8 == nil {
		return nil
	}
	return defaultClient8.Get(Key)
}
func get9(Key string) *paladin.Value {
	if defaultClient9 == nil {
		return nil
	}
	return defaultClient9.Get(Key)
}

// GetAll return all config key->value map.
func (c *Client) GetAll() *paladin.Map {
	return nil
}

func (c *Client) loadValuesFromPaths(ns []string) (map[string]*paladin.Value, error) {
	var err error
	values := make(map[string]*paladin.Value, len(ns))
	for _, v := range ns {
		if values[v], err = c.loadValue(v); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (c *Client) loadValue(ns string) (*paladin.Value, error) {
	data := c.aClient.GetContent(agollo.WithNamespace(ns))
	if data == "" {
		return nil, fmt.Errorf("no config found for %s", ns)
	}
	return paladin.NewValue(data, []byte(data)), nil
}

// Watch watch watch on a key. The configuration implements the setter interface, which is invoked when the configuration changes.
func (c *Client) Watch() {
	keys := c.priorityConfigs
	mm := c.configMaps
	var err error
	ll := len(keys)
	// 优先初始化的文件
	for idx := 0; idx < ll; idx++ {
		k := keys[idx]
		err = c.SetConf(k, mm[k])
		if err != nil {
			panic(fmt.Sprintf("file:%s, err:%v", k, err))
		}
	}

	// 其他文件初始化
	for k := range mm {
		if paladin.IsExistStringArray(k, keys) {
			continue
		}
		err = c.SetConf(k, mm[k])
		if err != nil {
			panic(fmt.Sprintf("file:%s, err:%v", k, err))
		}
	}
	//go func(mm map[string]paladin.Setter) {
	//	for event := range c.Client.WatchEvent(context.TODO(), c.Config.Namespaces...) {
	//		m, ok := mm[event.Key]
	//		if !ok || m == nil {
	//			continue
	//		}
	//		log.Debug("Apollo Watch file event :%s, size:%d", event.Key, len(event.Value))
	//		err = m.Set(event.Value)
	//		if err != nil {
	//			log.Error("Apollo Watch Set key:%s err:%v", event.Key, err)
	//		}
	//	}
	//}(mm)
}

func (c *Client) SetConf(key string, setter paladin.Setter) (err error) {
	v := c.aClient.GetContent(agollo.WithNamespace(key))
	if v == "" {
		return paladin.ErrNotExist
	}
	raw := []byte(v)
	err = setter.Set(raw)
	if err != nil {
		return
	}
	log.Debug("Apollo SetConf :%s, size:%d", key, len(raw))
	return
}

func Init(priority []string, configMaps map[string]paladin.Setter) (paladin.Client, error) {
	c := &Client{
		_Apollo:         newApolloEnv(),
		priorityConfigs: priority,
		configMaps:      configMaps,
	}

	ns := make([]string, 0, len(c.configMaps))
	for key := range c.configMaps {
		ns = append(ns, key)
	}
	ns = append(ns, defaultNS...)
	aClient, err := agollo.Start(&agollo.Conf{
		AppID:           c._Apollo.AppId,
		Cluster:         c._Apollo.Cluster,
		NameSpaceNames:  ns,
		MetaAddr:        c._Apollo.Endpoint,
		AccesskeySecret: c._Apollo.Secret,
		SyncTimeout:     9000000,
		Retry:           3,
		PollTimeout:     9000000,
	}, agollo.SkipLocalCache())
	if err != nil {
		log.Error("agollo start err:%v", err)
		panic(err)
	}
	c.aClient = aClient

	rawVal, err := c.loadValuesFromPaths(ns)
	if err != nil {
		return nil, err
	}
	c.rawVal = rawVal
	valMap := &paladin.Map{}
	valMap.Store(rawVal)
	c.values = valMap

	aClient.OnUpdate(c.updateFun)

	c.Watch()

	assignClient(c)
	return c, nil
}

func (c *Client) updateFun(event *agollo.ChangeEvent) {
	log.Info("apollo config changed: %s", event.Namespace)
	log.Info("apollo config changed: %s", event.Changes)
	c.reloadValue(event)
	// 基础db配置跳过update
	if c.configMaps[event.Namespace] == nil {
		return
	}
	err := c.SetConf(event.Namespace, c.configMaps[event.Namespace])
	if err != nil {
		log.Error("apollo config changed: %s, err:%v", event.Namespace, err)
	}
}

func (c *Client) reloadValue(event *agollo.ChangeEvent) {
	time.Sleep(200 * time.Millisecond)
	c.mx.Lock()
	defer c.mx.Unlock()
	name := event.Namespace
	val, err := c.loadValue(name)
	if err != nil {
		log.Error("load file %s error: %s, skipped", name, err)
		return
	}

	c.rawVal[name] = val
	c.values.Store(c.rawVal)
}

// Switch
// 开关
func Switch() bool {
	on := false
	getEnv := os.Getenv(_ApolloSwitch)
	if strings.Compare(strings.ToLower(getEnv), "on") == 0 {
		on = true
	}
	return on
}
