package semgrep

/*
#cgo LDFLAGS: -L. semgrep_bridge_core.so

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

typedef enum BridgeErrorCode {
  BEC_OK = 0, BEC_OUT_OF_MEMORY = 1, BEC_UNKNOWN_ERROR = 2, BEC_ERROR_MESSAGE = 3,
} BridgeErrorCode;

typedef BridgeErrorCode (*ReadFileFunc)(char const *, char **, size_t *, void *);

void bridge_ml_startup();
void bridge_ml_shutdown();

char *bridge_ml_semgrep_analyze(char const * const [], ReadFileFunc, void *);

int semgrep_analyze_proxy(char const * const argv[]) {
  int fd[2];
  int ret = pipe(fd);
  pid_t p = fork();
  if (p == 0) {
    close(fd[0]);
    dup2(fd[1], STDOUT_FILENO);
    bridge_ml_startup();
    bridge_ml_semgrep_analyze(argv, NULL, NULL);
    bridge_ml_shutdown();
    exit(0);
  } else if (p > 0) {
    close(fd[1]);
    wait(NULL);
    return fd[0];
  }
}

*/
import "C"
import (
	"bufio"
	"os"
	"strings"
	"unsafe"
)

func exec(argv []string) string {
	var c_argv [](*C.char) = make([](*C.char), 0, len(argv))
	c_argv = append(c_argv, C.CString("semgrep"))
	for _, v := range argv {
		c_argv = append(c_argv, C.CString(v))
	}
	c_argv = append(c_argv, nil)

	fd := C.semgrep_analyze_proxy((*(*C.char))(unsafe.Pointer(&c_argv[0])))

	r := os.NewFile(uintptr(fd), "")

	var sb strings.Builder
	fscanner := bufio.NewScanner(r)
	for fscanner.Scan() {
		s := fscanner.Text()
		if s != "." {
			sb.WriteString(fscanner.Text())
		}
	}
	r.Close()

	return sb.String()
}

func Run(lang string, rules string, fn string) string {
	return exec([]string{"-lang", lang, "-rules", rules, "-json", fn})
}
