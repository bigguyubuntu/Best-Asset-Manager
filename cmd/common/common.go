package cmn

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type messageLevel string
type messageLevels struct {
	Debug     messageLevel
	Test      messageLevel
	Info      messageLevel
	Operation messageLevel
	Warning   messageLevel
	Error     messageLevel
	Critical  messageLevel
}

var ErrorLevels = messageLevels{Debug: "DEBUG: ", Info: "INFO: ", Warning: "WARNING: ",
	Error: "ERROR: ", Critical: "CRITICAL: ", Test: "TEST: ", Operation: "Operation"}
var LogLevels = ErrorLevels

var isInitalized = false
var appName = ""

// msg can be a formatted string, and args can be the variables
// workes exacly like fmt.Printf. only that messege leve has to be second
func Log(msg string, level messageLevel, args ...any) {
	logger := log.New(os.Stdout, appName+" ", log.LUTC|log.Ldate|log.Ltime)
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	logger.Printf("%s %s\n", level, msg)
}

// Critical Errors will cause panic
func HandleError(err error, level messageLevel) {
	if err == nil {
		return
	}
	_, f, line, _ := runtime.Caller(1)
	// fmt.Println(f)
	// fn := runtime.FuncForPC(pc).Name()
	Log(fmt.Sprintf("%s:%d %v", f, line, err), level)
	switch level {
	// case ErrorLevels.Error:
	// 	debug.PrintStack()
	case ErrorLevels.Critical:
		panic("Critical Error logged above...")
	}
}
func InitCommon(env string, programName string) {
	// inits common functinality like logging, should be the first to init
	if isInitalized {
		return
	}
	if env != string(Envs.Development) &&
		env != string(Envs.Test) &&
		env != string(Envs.Production) {
		HandleError(errors.New("environment variables are not set correctly"), ErrorLevels.Critical)
	}
	envType := stringToEnvironment(env)
	appName = programName
	Log(fmt.Sprintf("Initalizing program with environment %s", env), LogLevels.Info)
	initEnvironment(envType)
	isInitalized = true
}
func DidProgramInit() bool {
	return isInitalized
}

func ConvertToMil(dollarAmount int) int {
	// 1.25$ is 1250mil todo handle rounding
	return int(dollarAmount * 1000)
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	HandleError(err, ErrorLevels.Error)
	return i
}

// takes an array of int
// and makes it a type any literal. so if a was []int
// we will return []any. This is useful
// for dodging type checks on lists
func ConvertTypeToAny(a []int) []any {
	s := make([]interface{}, len(a))
	for i, v := range a {
		s[i] = v
	}
	return s
}

func RandomStringGenerator(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	l := len(charset)
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		// Getting random character byte
		c := charset[rand.Intn(l)]
		s[i] = c
	}
	return string(s)
}

// converts an int array to a string, each char is seperated by sep param
func IntArrayToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
