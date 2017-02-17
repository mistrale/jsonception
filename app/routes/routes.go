// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tWebSocket struct {}
var WebSocket tWebSocket


func (_ tWebSocket) ListenExecutionRun(
		room_name string,
		ws interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "room_name", room_name)
	revel.Unbind(args, "ws", ws)
	return revel.MainRouter.Reverse("WebSocket.ListenExecutionRun", args).Url
}


type tGorpController struct {}
var GorpController tGorpController


func (_ tGorpController) Begin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Begin", args).Url
}

func (_ tGorpController) Commit(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Commit", args).Url
}

func (_ tGorpController) Rollback(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Rollback", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Suite(
		suite string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	return revel.MainRouter.Reverse("TestRunner.Suite", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tExecutions struct {}
var Executions tExecutions


func (_ tExecutions) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Executions.Index", args).Url
}

func (_ tExecutions) All(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Executions.All", args).Url
}

func (_ tExecutions) Get(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Executions.Get", args).Url
}

func (_ tExecutions) GetOne(
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Executions.GetOne", args).Url
}

func (_ tExecutions) GetOneTemplate(
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Executions.GetOneTemplate", args).Url
}

func (_ tExecutions) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Executions.Create", args).Url
}

func (_ tExecutions) Delete(
		id_exec int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_exec", id_exec)
	return revel.MainRouter.Reverse("Executions.Delete", args).Url
}

func (_ tExecutions) Update(
		id_exec int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_exec", id_exec)
	return revel.MainRouter.Reverse("Executions.Update", args).Url
}

func (_ tExecutions) Run(
		id_exec int,
		script string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_exec", id_exec)
	revel.Unbind(args, "script", script)
	return revel.MainRouter.Reverse("Executions.Run", args).Url
}


type tLibraries struct {}
var Libraries tLibraries


func (_ tLibraries) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Libraries.Create", args).Url
}

func (_ tLibraries) Delete(
		id_lib int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_lib", id_lib)
	return revel.MainRouter.Reverse("Libraries.Delete", args).Url
}

func (_ tLibraries) Update(
		id_lib int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_lib", id_lib)
	return revel.MainRouter.Reverse("Libraries.Update", args).Url
}

func (_ tLibraries) Run(
		libID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "libID", libID)
	return revel.MainRouter.Reverse("Libraries.Run", args).Url
}

func (_ tLibraries) GetHistory(
		libID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "libID", libID)
	return revel.MainRouter.Reverse("Libraries.GetHistory", args).Url
}

func (_ tLibraries) GetHistoryTemplate(
		libID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "libID", libID)
	return revel.MainRouter.Reverse("Libraries.GetHistoryTemplate", args).Url
}

func (_ tLibraries) GetOne(
		libID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "libID", libID)
	return revel.MainRouter.Reverse("Libraries.GetOne", args).Url
}

func (_ tLibraries) Get(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Libraries.Get", args).Url
}

func (_ tLibraries) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Libraries.Index", args).Url
}

func (_ tLibraries) All(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Libraries.All", args).Url
}


type tTests struct {}
var Tests tTests


func (_ tTests) Create(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Tests.Create", args).Url
}

func (_ tTests) Delete(
		id_test int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_test", id_test)
	return revel.MainRouter.Reverse("Tests.Delete", args).Url
}

func (_ tTests) Update(
		id_test int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id_test", id_test)
	return revel.MainRouter.Reverse("Tests.Update", args).Url
}

func (_ tTests) Run(
		testID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "testID", testID)
	return revel.MainRouter.Reverse("Tests.Run", args).Url
}

func (_ tTests) GetHistory(
		testID string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "testID", testID)
	return revel.MainRouter.Reverse("Tests.GetHistory", args).Url
}

func (_ tTests) GetHistoryTemplate(
		testID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "testID", testID)
	return revel.MainRouter.Reverse("Tests.GetHistoryTemplate", args).Url
}

func (_ tTests) Show(
		testID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "testID", testID)
	return revel.MainRouter.Reverse("Tests.Show", args).Url
}

func (_ tTests) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Tests.Index", args).Url
}

func (_ tTests) All(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Tests.All", args).Url
}

func (_ tTests) Get(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Tests.Get", args).Url
}

func (_ tTests) GetOneTemplate(
		testID int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "testID", testID)
	return revel.MainRouter.Reverse("Tests.GetOneTemplate", args).Url
}


type tTestHistory struct {}
var TestHistory tTestHistory


func (_ tTestHistory) Get(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestHistory.Get", args).Url
}


