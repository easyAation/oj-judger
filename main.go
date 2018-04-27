package main

import "github.com/open-fightcoder/oj-judger/judge"

//func main() {
//	cfgFile := flag.String("c", "cfg/cfg.toml.debug", "set config file")
//	flag.Parse()
//
//	common.Init(*cfgFile)
//	defer common.Close()
//
//	router := router.GetRouter()
//
//	graceful.LogListenAndServe(&http.Server{
//		Addr:    fmt.Sprintf(":%d", g.Conf().Run.HTTPPort),
//		Handler: router,
//	})
//
//	common.Close()
//}

func main() {
	judgeCpp := new(judge.JudgeCpp)
	judgeCpp.Compile("./scripts/test/cpp_hello.cpp")
	judgeCpp.Run("./scripts/test/2.in", "user.out")
}
