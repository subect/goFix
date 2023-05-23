package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime/pprof"
)

func Pinyin(c *gin.Context) {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	Mytttt()
	c.JSON(200, "nihoa")
}

func Mytttt() {
	for i := 0; i < 100; i++ {
		for i := 0; i < 100; i++ {
			for i := 0; i < 100; i++ {
				fmt.Println("nibihagcg")
			}
		}
	}

}
