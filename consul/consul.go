package consul

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var (
	_client *consulapi.Client
	conf    *Config
)

type Config struct {
	Address string
}

// Init 初始化consul连接
func Init(config *Config) error {
	if config == nil {
		config = &Config{
			Address: "127.0.0.1:8500",
		}
	}
	conf = config
	if conf.Address == "" {
		conf.Address = "127.0.0.1:8500"
	}
	return nil
}

func New() (*consulapi.Client, error) {
	if _client != nil {
		return _client, nil
	}

	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = conf.Address
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	_client = client
	return client, nil
}

// Register 注册服务到consul
func Register(name string, addr string, port int, tags ...string) error {
	client, err := New()
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("consul 实例空")
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = fmt.Sprintf("%s-%s:%d", name, addr, port)
	registration.Name = name
	registration.Port = port
	registration.Tags = tags
	registration.Address = addr

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.GRPC = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	return nil
}

// RegisterAPI 注册api服务到consul
func RegisterAPI(name string, addr string, port int, tags ...string) error {
	client, err := New()
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("consul 实例空")
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = fmt.Sprintf("%s-%s:%d", name, addr, port)
	registration.Name = name
	registration.Port = port
	registration.Tags = tags
	registration.Address = addr

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/health/check", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除

	registration.Check = check

	// 注册服务到consul
	if err := client.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	return nil
}

// DeRegister 取消consul注册的服务
func DeRegister(name string, addr string, port int) error {
	client, err := New()
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("consul 实例空")
	}
	client.Agent().ServiceDeregister(fmt.Sprintf("%s-%s:%d", name, addr, port))
	return nil
}

// FindNode 查找节点
func FindNode(servicename, tag string) (*consulapi.AgentService, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, errors.New("consul 实例空")
	}
	services, _, err := client.Health().Service(servicename, tag, true, nil)
	if err != nil {
		return nil, err
	}
	l := len(services)
	if l == 0 {
		return nil, nil
	}
	if l == 1 {
		return services[0].Service, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return services[r.Intn(l)%l].Service, nil
}

// FindService 从consul中发现服务，并返回grpc连接实例
func FindService(servicename, tag string) (*grpc.ClientConn, error) {
	node, err := FindNode(servicename, tag) //无tag视为grpc服务
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("微服务" + servicename + "不可用，稍后再试！")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", node.Address, node.Port), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// CheckHeath 健康检查
func CheckHeath(serviceid string) error {
	client, err := New()
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("consul 实例空")
	}
	// 健康检查
	//	a, b, _ := client.Agent().AgentHealthServiceByID(serviceid)
	return nil
}

// KVPut test
func KVPut(key string, values *[]byte, flags uint64) (*consulapi.WriteMeta, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("consul 实例空")
	}

	return client.KV().Put(&consulapi.KVPair{Key: key, Flags: flags, Value: *values}, nil)
}

// KVGet 获取值
func KVGet(key string, flags uint64) (*[]byte, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("consul 实例空")
	}

	// KV get值
	data, _, _ := client.KV().Get(key, nil)
	if data != nil {
		return &data.Value, nil
	}

	return nil, nil
}
