package util

// func getLocalIP() string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "127.0.0.1"
// 	}
// 	for _, address := range addrs {
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}
// 	return "127.0.0.1"
// }

// func getGrpcAddress() string {
// 	var (
// 		cfg *warden.ServerConfig
// 		ct  paladin.TOML
// 	)
// 	if err := paladin.Get("grpc.txt").Unmarshal(&ct); err != nil {
// 		return ""
// 	} else if err := ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
// 		return ""
// 	}

// 	tmp := strings.Split(cfg.Addr, ":")
// 	if len(tmp) < 2 {
// 		return ""
// 	}

// 	return fmt.Sprintf("grpc://%s:%s", getLocalIP(), tmp[1])
// }

// 注册discovery
//func RegisterDiscovery(appId string) (context.CancelFunc, error) {
//	resolver.Register(discovery.Builder())
//
//	hn, _ := os.Hostname()
//	dis := discovery.New(nil)
//	ins := &naming.Instance{
//		Zone:     env.Zone,
//		Env:      env.DeployEnv,
//		AppID:    appId,
//		Hostname: hn,
//		Addrs: []string{
//			getGrpcAddress(),
//		},
//	}
//	var (
//		cancel context.CancelFunc
//		err    error
//	)
//	if cancel, err = dis.Register(context.Background(), ins); err != nil {
//		return nil, err
//	}
//
//	return cancel, nil
//}
