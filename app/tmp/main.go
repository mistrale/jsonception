// GENERATED CODE - DO NOT EDIT
package main

import (
	"flag"
	"reflect"
	"github.com/revel/revel"
	_ "github.com/mistrale/jsonception/app"
	controllers "github.com/mistrale/jsonception/app/controllers"
	tests "github.com/mistrale/jsonception/tests"
	controllers0 "github.com/revel/modules/static/app/controllers"
	_ "github.com/revel/modules/testrunner/app"
	controllers1 "github.com/revel/modules/testrunner/app/controllers"
	websocket "golang.org/x/net/websocket"
	"github.com/revel/revel/testing"
)

var (
	runMode    *string = flag.String("runMode", "", "Run mode.")
	port       *int    = flag.Int("port", 0, "By default, read from app.conf")
	importPath *string = flag.String("importPath", "", "Go Import Path for the app.")
	srcPath    *string = flag.String("srcPath", "", "Path to the source root.")

	// So compiler won't complain if the generated code doesn't reference reflect package...
	_ = reflect.Invalid
)

func main() {
	flag.Parse()
	revel.Init(*runMode, *importPath, *srcPath)
	revel.INFO.Println("Running revel server")
	
	revel.RegisterController((*controllers.GorpController)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Begin",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Commit",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Rollback",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.WebSocket)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "ListenExecutionRun",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "room_name", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "ws", Type: reflect.TypeOf((**websocket.Conn)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers0.Static)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Serve",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ServeModule",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "moduleName", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers1.TestRunner)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					72: []string{ 
						"testSuites",
					},
				},
			},
			&revel.MethodType{
				Name: "Suite",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "suite", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Run",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "suite", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "test", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					125: []string{ 
					},
				},
			},
			&revel.MethodType{
				Name: "List",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Executions)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					33: []string{ 
						"exec",
						"ip",
					},
				},
			},
			&revel.MethodType{
				Name: "All",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					42: []string{ 
						"execs",
						"testID",
						"ip",
					},
				},
			},
			&revel.MethodType{
				Name: "Get",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetOne",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetOneTemplate",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					66: []string{ 
						"exec",
						"uuid",
					},
				},
			},
			&revel.MethodType{
				Name: "Create",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Delete",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_exec", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Update",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_exec", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Run",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_exec", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "script", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Libraries)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Create",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Delete",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_lib", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Update",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_lib", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Run",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetHistory",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "libID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteHistory",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_lib", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetHistoryTemplate",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "libID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					170: []string{ 
						"history",
					},
				},
			},
			&revel.MethodType{
				Name: "GetOne",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "libID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Get",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					191: []string{ 
						"library",
					},
				},
			},
			&revel.MethodType{
				Name: "All",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					198: []string{ 
						"libs",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Tests)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Create",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Delete",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_test", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Update",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "id_test", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Run",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "testID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetHistory",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "testID", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetHistoryTemplate",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "testID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					208: []string{ 
						"history",
					},
				},
			},
			&revel.MethodType{
				Name: "Show",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "testID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					214: []string{ 
					},
				},
			},
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					220: []string{ 
						"test",
					},
				},
			},
			&revel.MethodType{
				Name: "All",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					227: []string{ 
						"tests",
					},
				},
			},
			&revel.MethodType{
				Name: "Get",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetOneTemplate",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "testID", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					241: []string{ 
						"test",
					},
				},
			},
			
		})
	
	revel.RegisterController((*controllers.TestHistory)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Get",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetOne",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "uuid", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "testName", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					28: []string{ 
						"hist",
					},
				},
			},
			
		})
	
	revel.DefaultValidationKeys = map[string]map[int]string{ 
	}
	testing.TestSuites = []interface{}{ 
		(*tests.AppTest)(nil),
	}

	revel.Run(*port)
}
