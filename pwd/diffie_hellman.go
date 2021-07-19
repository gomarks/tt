package main

// 解决秘钥在双方不直接传递秘钥的情况下完成秘钥交换
// 甲：素数 p，底数 g，随机数 a. 计算 A = g ^ a mod p，把 p g A 的值发给乙
// 乙：选择随机数 b，计算 B = g ^ b mode p，同时再计算 s = A ^ b mode p，把 B 的值发给甲
// 甲：计算 s = B ^ a mode p
// 乙、甲计算出的 s 值相等
// s 是共同协商出的秘钥
// a 是甲的私钥，A 是甲的公钥
// b 是乙的私钥，B 是乙的公钥
// 甲知道：甲的公私钥，乙的公钥，协商出的秘钥，不知道乙的私钥
// 乙知道：乙的公司有，甲的公钥，协商出的秘钥，不知道甲的私钥
// ——————————————————
// meta 存 a
// chunk 存 p B

import (
	"context"
	dh "github.com/bensema/go-diffie-hellman"
	"golang.org/x/sync/semaphore"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"tt/common"
)

func main() {
	common.ShowTime(4, 100, dhm())
}

func dhm() func(s *semaphore.Weighted, waitGroup *sync.WaitGroup, idx int) {
	return func(s *semaphore.Weighted, waitGroup *sync.WaitGroup, idx int) {
		s.Acquire(context.Background(), 1)

		// 素数
		p, _ := dh.Prime(512)
		// 底数
		g := big.NewInt(1024)

		alice := dh.Dh{
			P:          p,
			G:          g,
			PrivateKey: big.NewInt(int64(rand.Intn(512))),
		}
		A := alice.Public()
		// log.Printf("A ： %s", A)

		bob := dh.Dh{
			P:          p,
			G:          g,
			PrivateKey: big.NewInt(int64(rand.Intn(512))),
		}
		B := bob.Public()
		// log.Printf("B ： %s", B)

		alice.AnswerKey = B
		bob.AnswerKey = A
		if alice.Sha256() != bob.Sha256() {
			log.Printf(alice.Sha256())
			log.Printf(bob.Sha256())
			panic("同样的 p 和 g 交换到的 s 不一样")
		}

		log.Print(idx)
		waitGroup.Done()
		s.Release(1)
	}
}
