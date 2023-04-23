package db

import (
	"db-proxy/utils"
	"log"
)

var sem *utils.Sem

func InitQ(qLen int) {
	sem = utils.NewSem(qLen)
}

func enterQ(prompt string) {
	sem.Acquire()
	log.Printf(" ==> %s\n", prompt)
}

func leaveQ(prompt string) {
	log.Printf(" <== %s\n", prompt)
	sem.Release()
}
